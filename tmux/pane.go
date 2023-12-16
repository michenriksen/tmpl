package tmux

import (
	"context"
	"fmt"
)

var paneOutputFormat = outputFormat("pane_id", "pane_path", "pane_index", "pane_width", "pane_height")

// Pane represents a tmux window pane.
type Pane struct {
	tmux       Runner
	sess       *Session
	win        *Window
	pane       *Pane
	env        map[string]string
	id         string
	path       string
	cmds       []string
	size       string
	width      string
	height     string
	index      string
	panes      []*Pane
	horizontal bool
	active     bool
	state      state
}

// NewPane creates a new Pane instance configured with the provided options and
// belonging to the provided window and optional parent pane.
//
// The window must be applied before being passed to this function.
// If a parent pane is provided, it must be applied before being passed to this
// function.
//
// NOTE: The pane is not created until [Pane.Apply] is called.
func NewPane(runner Runner, window *Window, parentPane *Pane, opts ...PaneOption) (*Pane, error) {
	if runner == nil {
		return nil, ErrNilRunner
	}

	if window == nil {
		return nil, ErrNilWindow
	}

	if window.sess == nil {
		return nil, ErrNilSession
	}

	if err := window.sess.checkState(); err != nil {
		return nil, fmt.Errorf("checking session state: %w", err)
	}

	if parentPane != nil {
		if err := parentPane.checkState(); err != nil {
			return nil, fmt.Errorf("checking parent pane state: %w", err)
		}
	}

	p := &Pane{
		tmux:  runner,
		sess:  window.sess,
		win:   window,
		pane:  parentPane,
		state: stateNew,
	}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, fmt.Errorf("applying pane option: %w", err)
		}
	}

	return p, nil
}

// Apply creates the tmux pane by invoking the split-window command using its
// internal [Runner] instance.
//
// If the pane is already applied, this method is a no-op.
//
// https://man.archlinux.org/man/tmux.1#split-window
func (p *Pane) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if p.IsClosed() {
		return ErrSessionClosed
	}

	if p.IsApplied() {
		return nil
	}

	target := p.win.Name()
	if p.pane != nil {
		target = p.pane.Name()
	}

	args := []string{"split-window", "-d", "-P", "-F", paneOutputFormat, "-t", target}

	args = append(args, p.envArgs()...)

	if p.path != "" {
		args = append(args, "-c", p.path)
	}

	if p.size != "" {
		args = append(args, "-l", p.size)
	}

	if p.horizontal {
		args = append(args, "-h")
	}

	output, err := p.tmux.Run(ctx, args...)
	if err != nil {
		return fmt.Errorf("running split-window command: %w", err)
	}

	if p.tmux.IsDryRun() {
		output = []byte(p.dryRunRecord())
	}

	records, err := parseOutput(output)
	if err != nil {
		return fmt.Errorf("parsing split-window command output: %w", err)
	}

	if err := p.update(records[0]); err != nil {
		return fmt.Errorf("updating pane data: %w", err)
	}

	if p.pane != nil {
		p.addPane(p)
	} else {
		p.win.addPane(p)
	}

	p.log("pane created")

	cmds := append(p.sess.onPaneCommands(), p.cmds...)

	return p.RunCommands(ctx, cmds...)
}

// RunCommands runs the provided commands inside the pane by invoking the
// send-keys tmux command using its internal [Runner] instance.
//
// The commands are automatically followed by a carriage return.
//
// If the pane is not applied, the method returns [ErrPaneNotApplied].
//
// If no commands are provided, the method is a no-op.
//
// https://man.archlinux.org/man/tmux.1#send-keys
func (p *Pane) RunCommands(ctx context.Context, cmds ...string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if len(cmds) == 0 {
		return nil
	}

	if err := p.checkState(); err != nil {
		return fmt.Errorf("checking pane state: %w", err)
	}

	for _, cmd := range cmds {
		args := []string{"send-keys", "-t", p.Name(), cmd, "C-m"}
		if _, err := p.tmux.Run(ctx, args...); err != nil {
			return fmt.Errorf("running send-keys command: %w", err)
		}

		p.log("pane send-keys", "cmd", cmd+"<cr>")
	}

	return nil
}

// Select selects the pane by invoking the select-pane command using its
// internal [Runner] instance.
//
// If the pane is not applied, the method returns [ErrPaneNotApplied].
//
// https://man.archlinux.org/man/tmux.1#select-pane
func (p *Pane) Select(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := p.checkState(); err != nil {
		return fmt.Errorf("checking pane state: %w", err)
	}

	if _, err := p.tmux.Run(ctx, "select-pane", "-t", p.Name()); err != nil {
		return fmt.Errorf("running select-pane command: %w", err)
	}

	p.log("pane selected")

	return nil
}

// Name returns the pane's fully qualified name.
//
// The name is composed of the session name, the window name separated by a
// colon, followed by a dot and the pane index.
func (p *Pane) Name() string {
	return fmt.Sprintf("%s.%s", p.win.Name(), p.index)
}

// IsApplied returns true if the pane has been applied with [Pane.Apply].
func (p *Pane) IsApplied() bool {
	return p.state == stateApplied
}

// IsClosed return true if the session has been closed.
func (p *Pane) IsClosed() bool {
	return p.state == stateClosed || p.win.IsClosed()
}

// IsActive returns true if the pane is configured as the active pane of its
// window.
func (p *Pane) IsActive() bool {
	return p.active
}

// NumPanes returns the number of panes in the pane.
func (p *Pane) NumPanes() int {
	return len(p.panes)
}

// String returns a string representation of the pane.
func (p *Pane) String() string {
	return fmt.Sprintf("pane %s", p.Name())
}

// update updates the pane's internal state from the provided output record.
func (p *Pane) update(record outputRecord) error {
	fieldsMap := map[string]*string{
		"pane_id":     &p.id,
		"pane_path":   &p.path,
		"pane_index":  &p.index,
		"pane_width":  &p.width,
		"pane_height": &p.height,
	}

	for k, v := range record {
		if v == "" {
			continue
		}

		if field, ok := fieldsMap[k]; ok {
			*field = v
		}
	}

	p.state = stateApplied

	return nil
}

func (p *Pane) checkState() error {
	if p.IsClosed() {
		return ErrPaneNotApplied
	}

	if p.IsApplied() {
		return nil
	}

	return ErrPaneNotApplied
}

func (p *Pane) addPane(pane *Pane) {
	p.panes = append(p.panes, pane)
	p.win.addPane(p)
}

func (p *Pane) envArgs() []string {
	if p.pane != nil {
		return envArgs(p.sess.env, p.win.env, p.pane.env, p.env)
	}

	return envArgs(p.sess.env, p.win.env, p.env)
}

func (p *Pane) log(msg string, args ...any) {
	p.tmux.Log(msg, append(args,
		"session", p.sess.Name(), "window", p.win.Name(), "pane", p.Name(),
		"pane_width", p.width, "pane_height", p.height)...)
}

// dryRunRecord returns an output record string to use when running in dry-run
// mode.
func (p *Pane) dryRunRecord() string {
	return outputRecord{
		"pane_id":     fmt.Sprintf("%%%d", p.sess.NumPanes()+1),
		"pane_path":   p.path,
		"pane_index":  fmt.Sprintf("%d", p.win.NumPanes()+1),
		"pane_width":  "40",
		"pane_height": "12",
	}.String()
}

// PaneOption configures a [Pane].
type PaneOption func(*Pane) error

// PaneWithSize configures the [Pane] size.
//
// The size can be specified as a percentage of the available space or as a
// number of lines or columns.
func PaneWithSize(size string) PaneOption {
	return func(p *Pane) error {
		p.size = size
		return nil
	}
}

// PaneWithPath configures the [Pane] working directory.
//
// If a pane is not configured with a working directory, the window's working
// directory is used instead.
func PaneWithPath(s string) PaneOption {
	return func(p *Pane) error {
		p.path = s
		return nil
	}
}

// PaneWithCommands configures the [Pane] with an initial shell commands.
//
// The commands will run in the pane after it has been created.
//
// NOTE: Commands are appended to the list of commands, so applying this option
// multiple times will add to the list of commands.
func PaneWithCommands(cmds ...string) PaneOption {
	return func(p *Pane) error {
		p.cmds = append(p.cmds, cmds...)
		return nil
	}
}

// PaneWithHorizontalDirection configures the [Pane] to be horizontal instead of
// vertical.
func PaneWithHorizontalDirection() PaneOption {
	return func(p *Pane) error {
		p.horizontal = true
		return nil
	}
}

// PaneAsActive configures the [Pane] to be the active pane of its window.
func PaneAsActive() PaneOption {
	return func(p *Pane) error {
		p.active = true
		return nil
	}
}

// PaneWithEnv configures the [Pane] with environment variables.
//
// Environment variables are inherited from session to window to pane. If a
// an environment variable is is named the same as an inherited variable, it
// will take precedence.
func PaneWithEnv(env map[string]string) PaneOption {
	return func(p *Pane) error {
		p.env = env
		return nil
	}
}
