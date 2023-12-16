package tmux

import (
	"bytes"
	"context"
	"fmt"
)

var windowOutputFormat = outputFormat(
	"window_id", "window_name", "window_path", "window_index",
	"window_width", "window_height",
)

// Window represents a tmux window.
type Window struct {
	tmux   Runner
	sess   *Session
	id     string
	name   string
	path   string
	cmds   []string
	index  string
	width  string
	height string
	env    map[string]string
	panes  []*Pane
	active bool
	state  state
}

// NewWindow creates a new Window instance configured with the provided options
// and belonging to the provided session.
//
// The session must be applied before being passed to this function.
//
// NOTE: The window is not created until [Window.Apply] is called.
func NewWindow(runner Runner, session *Session, opts ...WindowOption) (*Window, error) {
	if runner == nil {
		return nil, ErrNilRunner
	}

	if session == nil {
		return nil, ErrNilSession
	}

	if err := session.checkState(); err != nil {
		return nil, fmt.Errorf("checking session state: %w", err)
	}

	w := &Window{tmux: runner, sess: session, state: stateNew}

	for _, opt := range opts {
		if err := opt(w); err != nil {
			return nil, fmt.Errorf("applying window option: %w", err)
		}
	}

	return w, nil
}

// Apply creates the tmux window by invoking the new-window command using its
// internal [Runner] instance.
//
// If the window is already applied, this method is a no-op.
//
// If the window is the first window in the session, it is created with the -k
// flag to override the default initial window created by the new-session
// command.
//
// https://man.archlinux.org/man/tmux.1#new-window
func (w *Window) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if w.IsClosed() {
		return ErrSessionClosed
	}

	if w.IsApplied() {
		return nil
	}

	args := []string{"new-window", "-P", "-F", windowOutputFormat}

	if w.sess.NumWindows() == 0 {
		args = append(args, "-k", "-t", fmt.Sprintf("%s:^", w.sess.Name()))
	} else {
		args = append(args, "-t", fmt.Sprintf("%s:", w.sess.Name()))
	}

	args = append(args, envArgs(w.sess.env, w.env)...)

	if w.name != "" {
		args = append(args, "-n", w.name)
	}

	if w.path != "" {
		args = append(args, "-c", w.path)
	}

	output, err := w.tmux.Run(ctx, args...)
	if err != nil {
		return fmt.Errorf("running new-window command: %w", err)
	}

	if w.tmux.IsDryRun() {
		output = []byte(w.dryRunRecord())
	}

	records, err := parseOutput(output)
	if err != nil {
		return fmt.Errorf("parsing new-window command output: %w", err)
	}

	if err := w.update(records[0]); err != nil {
		return fmt.Errorf("updating window data: %w", err)
	}

	w.sess.addWindow(w)
	w.log("window created")

	cmds := append(w.sess.onWindowCommands(), w.cmds...)

	return w.RunCommands(ctx, cmds...)
}

// Select selects the window by invoking the select-window command using its
// internal [Runner] instance.
//
// If the window is not applied, the method returns [ErrWindowNotApplied].
//
// If the window has a pane configured to be the active pane, the pane is also
// selected.
//
// https://man.archlinux.org/man/tmux.1#select-window
func (w *Window) Select(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := w.checkState(); err != nil {
		return fmt.Errorf("checking window state: %w", err)
	}

	if _, err := w.tmux.Run(ctx, "select-window", "-t", w.Name()); err != nil {
		return fmt.Errorf("running select-window command: %w", err)
	}

	w.log("window selected")

	return w.selectPane(ctx)
}

// selectPane selects the pane configured as the active pane.
//
// If no pane is configured as the active pane, the first pane is selected.
func (w *Window) selectPane(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := w.checkState(); err != nil {
		return fmt.Errorf("checking window state: %w", err)
	}

	if w.NumPanes() == 0 {
		return nil
	}

	var activePane *Pane

	for _, p := range w.panes {
		if p.IsActive() {
			activePane = p
		}
	}

	// If no pane is set as active, select the first pane.
	if activePane == nil {
		baseIndex, err := w.tmux.Run(ctx, "show-option", "-gqv", "pane-base-index")
		if err != nil {
			return fmt.Errorf("running show-option command: %w", err)
		}

		pTarget := fmt.Sprintf("%s.%s", w.Name(), bytes.TrimSpace(baseIndex))
		if _, err := w.tmux.Run(ctx, "select-pane", "-t", pTarget); err != nil {
			return fmt.Errorf("running select-pane command: %w", err)
		}

		return nil
	}

	if err := activePane.Select(ctx); err != nil {
		return fmt.Errorf("selecting active window pane: %w", err)
	}

	return nil
}

// RunCommands runs the provided commands inside the window by invoking the
// send-keys tmux command using its internal [Runner] instance.
//
// The commands are automatically followed by a carriage return.
//
// If the window is not applied, the method returns [ErrWindowNotApplied].
//
// If no commands are provided, the method is a no-op.
//
// https://man.archlinux.org/man/tmux.1#send-keys
func (w *Window) RunCommands(ctx context.Context, cmds ...string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if len(cmds) == 0 {
		return nil
	}

	if err := w.checkState(); err != nil {
		return fmt.Errorf("checking window state: %w", err)
	}

	for _, cmd := range cmds {
		args := []string{"send-keys", "-t", w.Name(), cmd, "C-m"}
		if _, err := w.tmux.Run(ctx, args...); err != nil {
			return fmt.Errorf("running send-keys command: %w", err)
		}

		w.log("window send-keys", "cmd", cmd+"<cr>")
	}

	return nil
}

// Name returns the window's fully qualified name.
//
// The name is composed of the session name and the window name separated by a
// colon.
func (w *Window) Name() string {
	return fmt.Sprintf("%s:%s", w.sess.Name(), w.name)
}

// IsApplied returns true if the window has been applied with [Window.Apply].
func (w *Window) IsApplied() bool {
	return w.state == stateApplied
}

// IsClosed returns true if the session has been closed.
func (w *Window) IsClosed() bool {
	return w.state == stateClosed || w.sess.IsClosed()
}

// IsActive returns true if the window is configured as the active window of its
// session.
func (w *Window) IsActive() bool {
	return w.active
}

// NumPanes returns the number of panes in the window.
func (w *Window) NumPanes() int {
	return len(w.panes)
}

// String returns a string representation of the window.
func (w *Window) String() string {
	return fmt.Sprintf("window %s", w.Name())
}

// update updates the window's internal state from the provided output record.
func (w *Window) update(record outputRecord) error {
	fieldsMap := map[string]*string{
		"window_id":     &w.id,
		"window_name":   &w.name,
		"window_path":   &w.path,
		"session_path":  &w.path,
		"window_index":  &w.index,
		"window_width":  &w.width,
		"window_height": &w.height,
	}

	for k, v := range record {
		if v == "" {
			continue
		}

		if field, ok := fieldsMap[k]; ok {
			*field = v
		}
	}

	w.state = stateApplied

	return nil
}

func (w *Window) addPane(p *Pane) {
	w.panes = append(w.panes, p)
}

func (w *Window) checkState() error {
	if w.IsClosed() {
		return ErrSessionClosed
	}

	if !w.IsApplied() {
		return ErrWindowNotApplied
	}

	return nil
}

func (w *Window) log(msg string, args ...any) {
	w.tmux.Log(msg, append(args, "session", w.sess.Name(), "window", w.Name())...)
}

// dryRunRecord returns an output record string to use when running in dry-run
// mode.
func (w *Window) dryRunRecord() string {
	wID := w.sess.NumWindows() + 1

	return outputRecord{
		"window_id":     fmt.Sprintf("@%d", wID),
		"window_name":   w.name,
		"window_path":   w.path,
		"window_index":  fmt.Sprintf("%d", wID),
		"window_width":  "80",
		"window_height": "24",
	}.String()
}

// WindowOption configures a [Window].
type WindowOption func(*Window) error

// WindowWithName configures the [Window] with a name.
func WindowWithName(name string) WindowOption {
	return func(w *Window) error {
		if name == "" {
			return fmt.Errorf("window name cannot be empty")
		}

		w.name = name

		return nil
	}
}

// WindowWithPath configures the [Window] with a working directory.
//
// If a window is not configured with a working directory, the session's working
// directory is used instead.
func WindowWithPath(p string) WindowOption {
	return func(w *Window) error {
		w.path = p
		return nil
	}
}

// WindowWithCommands configures the [Window] with an initial shell commands.
//
// The commands will run in the window after it has been created.
//
// NOTE: Commands are appended to the list of commands, so applying this option
// multiple times will add to the list of commands.
func WindowWithCommands(cmds ...string) WindowOption {
	return func(w *Window) error {
		w.cmds = append(w.cmds, cmds...)
		return nil
	}
}

// WindowAsActive configures the [Window] to be the active window of its
// session.
func WindowAsActive() WindowOption {
	return func(w *Window) error {
		w.active = true
		return nil
	}
}

// WindowWithEnv configures the [Window] with environment variables.
//
// Environment variables are inherited from session to window to pane. If a
// an environment variable is is named the same as an inherited variable, it
// will take precedence.
func WindowWithEnv(env map[string]string) WindowOption {
	return func(w *Window) error {
		w.env = env
		return nil
	}
}
