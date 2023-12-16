package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/invopop/validation"
	"github.com/lmittmann/tint"

	"github.com/michenriksen/tmpl/config"
	"github.com/michenriksen/tmpl/internal/env"
)

// skipLogErrors contains errors that should not be logged.
var skipLogErrors = []error{ErrVersion, ErrHelp}

func (a *App) initLogger() {
	a.logger = a.newLogger()

	if a.tmux != nil {
		a.tmux.SetLogger(a.logger)
	}
}

// newLogger creates a new structured logger that writes to the provided writer
// configured with the provided options.
func (a *App) newLogger() *slog.Logger {
	hOpts := &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		ReplaceAttr: a.attrReplacer,
	}

	if a.opts != nil {
		if a.opts.Debug {
			hOpts.Level = slog.LevelDebug
		}

		if a.opts.Quiet {
			hOpts.Level = slog.LevelWarn
		}
	}

	var handler slog.Handler

	if a.opts != nil && a.opts.JSON {
		handler = a.newJSONHandler(hOpts)
	} else {
		handler = a.newTextHandler(hOpts)
	}

	return slog.New(handler)
}

func (a *App) newTextHandler(hOpts *slog.HandlerOptions) slog.Handler {
	tOpts := &tint.Options{
		Level:       hOpts.Level,
		ReplaceAttr: hOpts.ReplaceAttr,
		TimeFormat:  "15:04:05",
		NoColor:     envNoColor(),
	}

	handler := tint.NewHandler(a.out, tOpts)

	return handler
}

func (a *App) newJSONHandler(hOpts *slog.HandlerOptions) *slog.JSONHandler {
	return slog.NewJSONHandler(a.out, hOpts)
}

// handleErr logs and returns the provided error.
//
// If err is nil, this method is a no-op.
// If err is in [skipLogErrors], logging is skipped.
func (a *App) handleErr(err error) error {
	if err == nil {
		return nil
	}

	for _, skipErr := range skipLogErrors {
		if errors.Is(err, skipErr) {
			return err
		}
	}

	logger := a.logger
	if logger == nil {
		logger = a.newLogger()
	}

	var decodeErr config.DecodeError
	if errors.As(err, &decodeErr) {
		logger.Error("configuration file cannot be decoded", "path", decodeErr.Path())

		if wrapped := decodeErr.Unwrap(); wrapped != nil {
			logger.Warn(wrapped.Error())
		}

		return ErrInvalidConfig
	}

	var verrs validation.Errors
	if errors.As(err, &verrs) {
		logger.Error("configuration file is invalid", "errors", len(verrs))
		logValidationErrs(logger, verrs, "")

		return ErrInvalidConfig
	}

	logger.Error(err.Error())

	return err
}

// logValidationErrs logs validation errors recursively.
func logValidationErrs(logger *slog.Logger, verrs validation.Errors, fieldPrfx string) {
	for field, err := range verrs {
		field = fieldPrfx + field

		if verr, ok := err.(validation.Errors); ok {
			logValidationErrs(logger, verr, field+".")
			continue
		}

		logger.Warn(fmt.Sprintf("%s %v", field, err.Error()), "field", field)
	}
}

func envNoColor() bool {
	// See https://no-color.org/
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return true
	}

	// Check application specific environment variable.
	if _, ok := env.LookupEnv("NO_COLOR"); ok {
		return true
	}

	// See https://bixense.com/clicolors/
	if _, ok := os.LookupEnv("CLICOLOR_FORCE"); ok {
		return false
	}

	// $TERM is often set to `dumb` to indicate that the terminal is very basic
	// and sometimes if the current command output is redirected to a file or
	// piped to another command.
	if os.Getenv("TERM") == "dumb" {
		return true
	}

	return false
}
