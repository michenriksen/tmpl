// Package cli is the application entry point.
package cli

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/michenriksen/tmpl/config"
	"github.com/michenriksen/tmpl/internal/env"
	"github.com/michenriksen/tmpl/tmux"
)

const (
	cmdInit  = "init"
	cmdCheck = "check"
)

// ErrInvalidConfig is returned when a configuration is invalid.
var ErrInvalidConfig = fmt.Errorf("invalid configuration")

// App is the main command-line application.
//
// The application orchestrates the loading of options and configuration, and
// applying the configuration to a new tmux session.
//
// The application can be configured with options to facilitate testing, such as
// providing a writer for output instead of os.Stdout, or providing a pre-loaded
// configuration instead of loading it from a configuration file.
type App struct {
	opts         *options
	cfg          *config.Config
	tmux         tmux.Runner
	sess         *tmux.Session
	logger       *slog.Logger
	attrReplacer func([]string, slog.Attr) slog.Attr
	out          io.Writer
}

// NewApp creates a new command-line application.
func NewApp(opts ...AppOption) (*App, error) {
	app := &App{out: os.Stdout}

	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, fmt.Errorf("applying application option: %w", err)
		}
	}

	return app, nil
}

// Run runs the command-line application.
//
// Different logic is performed dependening on the sub-command provided as the
// first command-line argument.
func (a *App) Run(ctx context.Context, args ...string) error {
	var (
		cmd string
		err error
	)

	if len(args) > 0 {
		cmd = args[0]
	}

	switch cmd {
	case cmdCheck:
		if a.opts == nil {
			if a.opts, err = parseCheckOptions(args[1:], a.out); err != nil {
				return a.handleErr(err)
			}
		}

		return a.handleErr(a.runCheck(ctx))
	case cmdInit:
		if a.opts == nil {
			if a.opts, err = parseInitOptions(args[1:], a.out); err != nil {
				return a.handleErr(err)
			}
		}

		return a.handleErr(a.runInit(ctx))
	default:
		if a.opts == nil {
			if a.opts, err = parseApplyOptions(args, a.out); err != nil {
				return a.handleErr(err)
			}
		}

		return a.handleErr(a.runApply(ctx))
	}
}

func (a *App) loadConfig() (err error) {
	if a.cfg != nil {
		return nil
	}

	wd, err := env.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %w", err)
	}

	if a.opts.ConfigPath == "" {
		a.opts.ConfigPath, err = config.FindConfigFile(wd)
		if err != nil {
			return fmt.Errorf("finding configuration file: %w", err)
		}
	}

	if a.cfg, err = config.FromFile(a.opts.ConfigPath); err != nil {
		return err //nolint:wrapcheck // Wrapping is done by caller.
	}

	a.logger.Info("configuration file loaded", "path", a.opts.ConfigPath)

	return nil
}

// AppOptions configures an [App].
type AppOption func(*App) error

// WithOptions configures the [App] to use provided options instead of
// parsing command-line flags.
//
// This option is intended for testing purposes only.
func WithOptions(opts *options) AppOption {
	return func(a *App) error {
		a.opts = opts
		return nil
	}
}

// WithConfig configures the [App] to use provided configuration instead of
// loading it from a configuration file.
//
// This option is intended for testing purposes only.
func WithConfig(cfg *config.Config) AppOption {
	return func(a *App) error {
		a.cfg = cfg
		return nil
	}
}

// WithOutputWriter configures the [App] to use provided writer for output
// instead of os.Stdout.
//
// This option is intended for testing purposes only.
func WithOutputWriter(w io.Writer) AppOption {
	return func(a *App) error {
		a.out = w
		return nil
	}
}

// WithTmux configures the [App] to use provided tmux runner instead of the
// constructing a new one from configuration.
//
// This option is intended for testing purposes only.
func WithTmux(tm tmux.Runner) AppOption {
	return func(a *App) error {
		a.tmux = tm
		return nil
	}
}

// WithSlogAttrReplacer configures the [App] to use provided function for
// replacing slog attributes.
//
// This option is intended for testing purposes only.
func WithSlogAttrReplacer(f func([]string, slog.Attr) slog.Attr) AppOption {
	return func(a *App) error {
		a.attrReplacer = f
		return nil
	}
}
