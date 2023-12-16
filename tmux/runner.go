package tmux

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// DefaultTmux is the default name of the tmux executable.
const DefaultTmux = "tmux"

// Runner runs tmux commands and provide logging functionality.
type Runner interface {
	// Run runs the tmux command in a context-aware manner with the provided
	// arguments and return the output.
	Run(ctx context.Context, args ...string) ([]byte, error)
	// Execve runs the tmux command with the provided arguments using the execve
	// syscall ([syscall.Exec]) to replace the current process with the tmux
	// process.
	Execve(args ...string) error
	// IsDryRun returns true if the runner is in dry-run mode.
	//
	// Dry-run mode means that the runner will not actually run any tmux commands
	// but only log the commands that would have been run and return empty output.
	IsDryRun() bool
	// Debug writes a debug message using a [slog.Logger].
	Debug(msg string, args ...any)
	// Log writes an info message using a [slog.Logger].
	Log(msg string, args ...any)
	// SetLogger sets the logger used by the runner.
	SetLogger(logger *slog.Logger)
}

// OSCommandRunner runs a command with the provided name and arguments and
// returns the output.
type OSCommandRunner func(ctx context.Context, name string, args ...string) (output []byte, err error)

// SyscallExecRunner runs a command with the provided arguments using the execve
// syscall ([syscall.Exec]) to replace the current process with the new
// process.
type SyscallExecRunner func(string, []string, []string) error

// defaultOSCmdRunner is the default [OSCommandRunner] implementation.
var defaultOSCmdRunner OSCommandRunner = func(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err //nolint:wrapcheck // Wrapping is done by caller.
	}

	return output, nil
}

// state represents the state of a tmux session, window, or pane.
type state int

const (
	stateNew     state = iota // New and not yet applied.
	stateApplied              // Applied with a tmux command.
	stateClosed               // Closed or removed with a tmux command.
)

// DefaultRunner is the default [Runner] implementation.
type DefaultRunner struct {
	logger       *slog.Logger
	cmdRunner    OSCommandRunner
	execveRunner SyscallExecRunner
	tmux         string
	tmuxOpts     []string
	dryRun       bool
}

// NewRunner creates a new [Runner] with the provided options.
func NewRunner(opts ...RunnerOption) (*DefaultRunner, error) {
	c := &DefaultRunner{
		tmux:         DefaultTmux,
		cmdRunner:    defaultOSCmdRunner,
		execveRunner: syscall.Exec,
		logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}

	return c, nil
}

// Run runs the tmux command in a context-aware manner with the provided
// arguments and returns the output.
func (c *DefaultRunner) Run(ctx context.Context, args ...string) ([]byte, error) {
	start := time.Now()

	args = append(c.tmuxOpts, args...)

	msg := "command successful"

	if c.dryRun {
		c.Debug(msg, "name", c.tmux, "args", args, "dur", time.Since(start))
		return []byte{}, nil
	}

	output, err := c.cmdRunner(ctx, c.tmux, args...)
	if err != nil {
		msg = "command failed"
	}

	c.Debug(msg, "name", c.tmux, "args", args, "output", strings.TrimSpace(string(output)), "dur", time.Since(start))

	return output, err
}

// Execve runs the tmux command with the provided arguments using the execve
// syscall ([syscall.Exec]) to replace the current process with the tmux
// process.
func (c *DefaultRunner) Execve(args ...string) error {
	start := time.Now()

	args = append(c.tmuxOpts, args...)

	absPath, err := exec.LookPath(c.tmux)
	if err != nil {
		return fmt.Errorf("looking up absolute path for %s executable: %w", c.tmux, err)
	}

	args = append([]string{absPath}, args...)

	msg := "execve successful"

	if c.dryRun {
		c.Debug(msg, "path", c.tmux, "args", args[1:], "dur", time.Since(start))
		return nil
	}

	err = c.execveRunner(absPath, args, os.Environ())
	if err != nil {
		msg = "execve failed"
	}

	c.Debug(msg, "path", absPath, "args", args[1:], "dur", time.Since(start))

	return err
}

// IsDryRun returns true if the runner is in dry-run mode.
//
// Dry-run mode means that the runner will not actually run any tmux commands
// but only log the commands that would have been run and return empty output.
func (c *DefaultRunner) IsDryRun() bool {
	return c.dryRun
}

// Debug logs a debug message using a [slog.Logger].
func (c *DefaultRunner) Debug(msg string, args ...any) {
	if c.dryRun {
		args = append(args, "dry_run", true)
	}

	c.logger.Debug(msg, args...)
}

// Log logs an info message using a [slog.Logger].
func (c *DefaultRunner) Log(msg string, args ...any) {
	if c.dryRun {
		args = append(args, "dry_run", true)
	}

	c.logger.Info(msg, args...)
}

// SetLogger sets the logger used by the runner.
func (c *DefaultRunner) SetLogger(logger *slog.Logger) {
	c.logger = logger
}

// RunnerOption configures a [DefaultRunner].
type RunnerOption func(*DefaultRunner) error

// WithTmux configures the runner to use provided name as the tmux executable.
//
// The default is "tmux".
func WithTmux(name string) RunnerOption {
	return func(c *DefaultRunner) error {
		c.tmux = name
		return nil
	}
}

// WithTmuxOptions configures the runner with additional tmux options to be
// added to all tmux command invocations.
func WithTmuxOptions(opts ...string) RunnerOption {
	return func(c *DefaultRunner) error {
		c.tmuxOpts = opts
		return nil
	}
}

// WithOSCommandRunner configures the runner to use the provided
// [OSCommandRunner] for running tmux commands.
//
// This option is intended for testing purposes only.
func WithOSCommandRunner(runner OSCommandRunner) RunnerOption {
	return func(c *DefaultRunner) error {
		c.cmdRunner = runner
		return nil
	}
}

// WithSyscallExecRunner configures the runner to use the provided
// [SyscallExecRunner] for running tmux commands.
//
// This option is intended for testing purposes only.
func WithSyscallExecRunner(runner SyscallExecRunner) RunnerOption {
	return func(c *DefaultRunner) error {
		c.execveRunner = runner
		return nil
	}
}

// WithLogger configures the runner to use the provided [slog.Logger] for
// logging.
//
// The default is a no-op logger writing to [os.Discard].
func WithLogger(logger *slog.Logger) RunnerOption {
	return func(c *DefaultRunner) error {
		c.logger = logger
		return nil
	}
}

// WithDryRunMode configures the runner to run in dry-run mode.
//
// Dry-run mode means that the runner will not actually run any tmux commands
// but only log the commands that would have been run and return empty output.
//
// The default is false.
func WithDryRunMode(enable bool) RunnerOption {
	return func(c *DefaultRunner) error {
		c.dryRun = enable
		return nil
	}
}
