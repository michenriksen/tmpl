package tmux

import (
	"context"
	"fmt"
	"strings"
)

var sessionOutputFormat = outputFormat("session_id", "session_name", "session_path")

// Session represents a tmux session.
type Session struct {
	tmux    Runner
	id      string
	path    string
	name    string
	winCmd  string
	paneCmd string
	anyCmd  string
	env     map[string]string
	windows []*Window
	state   state
}

// NewSession creates a new Session instance configured with the provided
// options.
//
// NOTE: The session is not created until [Session.Apply] is called.
func NewSession(runner Runner, opts ...SessionOption) (*Session, error) {
	if runner == nil {
		return nil, ErrNilRunner
	}

	s := &Session{tmux: runner, state: stateNew}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, fmt.Errorf("applying session option: %w", err)
		}
	}

	return s, nil
}

// Apply creates the tmux session by invoking the new-session command using its
// internal [Runner] instance.
//
// If the session is already applied, this method is a no-op.
//
// https://man.archlinux.org/man/tmux.1#new-session
func (s *Session) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if s.IsClosed() {
		return ErrSessionClosed
	}

	if s.IsApplied() {
		return nil
	}

	args := []string{"new-session", "-d", "-P", "-F", sessionOutputFormat}

	if s.name != "" {
		args = append(args, "-s", s.name)
	}

	output, err := s.tmux.Run(ctx, args...)
	if err != nil {
		return fmt.Errorf("running new-session command: %w", err)
	}

	if s.tmux.IsDryRun() {
		output = []byte(s.dryRunRecord())
	}

	records, err := parseOutput(output)
	if err != nil {
		return fmt.Errorf("parsing new-session command output: %w", err)
	}

	if err := s.update(records[0]); err != nil {
		return fmt.Errorf("updating session data: %w", err)
	}

	s.log("session created")

	return nil
}

// Attach attaches the current client to the session by invoking a tmux command
// using its internal [Runner] instance.
//
// NOTE: [syscall.Exec] is used to replace the current process with the tmux
// client process. This means that this method will never return if it
// succeeds.
//
// If the TMUX environment variable is set, it is assumed that the current
// process is already attached to a tmux session. In this case, the
// switch-client command is used instead of the attach-session command.
//
// https://man.archlinux.org/man/tmux.1#attach-session
// https://man.archlinux.org/man/tmux.1#switch-client
func (s *Session) Attach(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := s.checkState(); err != nil {
		return fmt.Errorf("checking session state: %w", err)
	}

	if inTmux() {
		s.log("switching client to session")

		if err := s.tmux.Execve("switch-client", "-t", s.name); err != nil {
			return fmt.Errorf("running switch-client command: %w", err)
		}

		return nil
	}

	s.log("attaching client to session")

	if err := s.tmux.Execve("attach-session", "-t", s.name); err != nil {
		return fmt.Errorf("running attach-session command: %w", err)
	}

	return nil
}

// SelectActive selects the window configured as the active window by invoking
// the select-window command using its internal [Runner] instance.
//
// If no window is configured as active, the first window is selected.
//
// https://man.archlinux.org/man/tmux.1#select-window
func (s *Session) SelectActive(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := s.checkState(); err != nil {
		return fmt.Errorf("checking session state: %w", err)
	}

	if s.NumWindows() == 1 {
		return s.windows[0].selectPane(ctx)
	}

	activeWin := s.windows[0]

	for _, w := range s.windows[1:] {
		if w.IsActive() {
			activeWin = w
		}
	}

	return activeWin.Select(ctx)
}

// Name returns the session name.
func (s *Session) Name() string {
	return s.name
}

// NumWindows returns the number of windows in the session.
func (s *Session) NumWindows() int {
	return len(s.windows)
}

// NumPanes returns the number of panes across all windows in the session.
func (s *Session) NumPanes() int {
	n := 0

	for _, w := range s.windows {
		n += len(w.panes)
	}

	return n
}

// Close closes the session by invoking the kill-session command using its
// internal [Runner] instance.
//
// If the session is already closed or not applied, this method is a no-op.
//
// Any subsequent calls to command-invoking methods on the session or any of its
// windows or panes will return an [ErrSessionClosed] error.
//
// https://man.archlinux.org/man/tmux.1#kill-session
func (s *Session) Close() error {
	if s.IsClosed() || !s.IsApplied() {
		return nil
	}

	if _, err := s.tmux.Run(context.Background(), "kill-session", "-t", s.name); err != nil {
		return fmt.Errorf("running kill-session command: %w", err)
	}

	s.state = stateClosed
	s.windows = nil

	s.tmux.Debug("session closed", "session", s.name)

	return nil
}

// IsClosed returns true if the session has been closed with [Session.Close].
func (s *Session) IsClosed() bool {
	return s.state == stateClosed
}

// IsApplied returns true if the session has been applied with [Session.Apply].
func (s *Session) IsApplied() bool {
	return s.state == stateApplied
}

// String returns a string representation of the session.
func (s *Session) String() string {
	return fmt.Sprintf("session %s", s.Name())
}

// update updates the session's internal state from the provided output record.
func (s *Session) update(record outputRecord) error {
	fieldsMap := map[string]*string{
		"session_id":   &s.id,
		"session_name": &s.name,
		"session_path": &s.path,
	}

	for k, v := range record {
		if v == "" {
			continue
		}

		if field, ok := fieldsMap[k]; ok {
			*field = v
		}
	}

	s.state = stateApplied

	return nil
}

// addWindow adds a window to the session.
func (s *Session) addWindow(w *Window) {
	s.windows = append(s.windows, w)
}

// checkState checks that the session is applied and not closed.
//
// Returns [ErrSessionClosed] if the session is closed.
//
// Returns [ErrSessionNotApplied] if the session is not applied.
func (s *Session) checkState() error {
	if s.IsClosed() {
		return ErrSessionClosed
	}

	if !s.IsApplied() {
		return ErrSessionNotApplied
	}

	return nil
}

func (s *Session) onWindowCommands() []string {
	var cmds []string

	if s.anyCmd != "" {
		cmds = append(cmds, s.anyCmd)
	}

	if s.winCmd != "" {
		cmds = append(cmds, s.winCmd)
	}

	return cmds
}

func (s *Session) onPaneCommands() []string {
	var cmds []string

	if s.anyCmd != "" {
		cmds = append(cmds, s.anyCmd)
	}

	if s.paneCmd != "" {
		cmds = append(cmds, s.paneCmd)
	}

	return cmds
}

func (s *Session) log(msg string, args ...any) {
	if s.NumWindows() > 0 {
		args = append(args, "windows", s.NumWindows())
	}

	if s.NumPanes() > 0 {
		args = append(args, "panes", s.NumPanes())
	}

	s.tmux.Log(msg, append(args, "session", s.name)...)
}

// dryRunRecord returns an output record string to use when running in dry-run
// mode.
func (s *Session) dryRunRecord() string {
	return outputRecord{
		"session_id":   "$0",
		"session_name": s.name,
		"session_path": s.path,
	}.String()
}

// SessionOption configures a [Session].
type SessionOption func(*Session) error

// SessionWithName configures the [Session] name.
func SessionWithName(name string) SessionOption {
	return func(s *Session) error {
		s.name = name
		return nil
	}
}

// SessionWithPath configures the [Session] working directory.
func SessionWithPath(path string) SessionOption {
	return func(s *Session) error {
		s.path = path

		return nil
	}
}

// SessionWithOnWindowCommand configures a [Session] with a shell command that
// will be run in all created windows.
//
// If a command is also configured with [SessionWithOnAnyCommand], the window
// command will be run after the any command.
func SessionWithOnWindowCommand(cmd string) SessionOption {
	return func(s *Session) error {
		s.winCmd = strings.TrimSpace(cmd)
		return nil
	}
}

// SessionWithOnPaneCommand configures a [Session] with a shell command that
// will be run in all created panes.
//
// If a command is also configured with [SessionWithOnAnyCommand], the pane
// command will be run after the any command.
func SessionWithOnPaneCommand(cmd string) SessionOption {
	return func(s *Session) error {
		s.paneCmd = strings.TrimSpace(cmd)
		return nil
	}
}

// SessionWithOnAnyCommand configures a [Session] with a shell command that will
// be run in all created windows and panes.
//
// If a command is also configured with [SessionWithOnWindowCommand] or
// [SessionWithOnPaneCommand], the any command will be run first.
func SessionWithOnAnyCommand(cmd string) SessionOption {
	return func(s *Session) error {
		s.anyCmd = strings.TrimSpace(cmd)
		return nil
	}
}

// SessionWithEnv configures the [Session] environment variables.
//
// Environment variables are inherited from session to window to pane. If a
// window or pane is confiured with a similarly named environment variable, it
// will take precedence over the session environment variable.
func SessionWithEnv(env map[string]string) SessionOption {
	return func(s *Session) error {
		s.env = env
		return nil
	}
}

// GetSessions returns a list of current tmux sessions by invoking the
// list-sessions command using the provided [Runner] instance.
//
// https://man.archlinux.org/man/tmux.1#list-sessions
func GetSessions(ctx context.Context, runner Runner) ([]*Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if runner == nil {
		return nil, ErrNilRunner
	}

	output, err := runner.Run(ctx, "list-sessions", "-F", sessionOutputFormat)
	if err != nil {
		return nil, fmt.Errorf("running list-sessions command: %w", err)
	}

	records, err := parseOutput(output)
	if err != nil {
		return nil, fmt.Errorf("parsing list-sessions command output: %w", err)
	}

	sessions := make([]*Session, len(records))

	for i, record := range records {
		s := &Session{tmux: runner, state: stateApplied}
		if err := s.update(record); err != nil {
			return nil, fmt.Errorf("updating session data: %w", err)
		}

		sessions[i] = s
	}

	return sessions, nil
}
