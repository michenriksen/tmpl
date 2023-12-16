package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/michenriksen/tmpl/tmux"
)

// Apply applies the provided tmux session configuration using the provided
// tmux command.
//
// If a session with the same name already exists, it is assumed to be in the
// correct state and the session is returned. Otherwise, a new session is
// created and returned.
//
// If the provided configuration is invalid, an error is returned. Caller can
// check for validity beforehand by calling [config.Config.Validate] if needed.
func Apply(ctx context.Context, cfg *Config, runner tmux.Runner) (*tmux.Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration file: %w", err)
	}

	sCfg := cfg.Session

	sessions, err := tmux.GetSessions(ctx, runner)
	if err != nil {
		return nil, fmt.Errorf("getting current tmux sessions: %w", err)
	}

	for _, s := range sessions {
		if s.Name() == sCfg.Name {
			return s, nil
		}
	}

	session, err := tmux.NewSession(runner, makeSessionOpts(cfg.Session)...)
	if err != nil {
		return fatalf(session, "creating session: %w", err)
	}

	if err := session.Apply(ctx); err != nil {
		return fatalf(session, "applying %s: %w", session, err)
	}

	for _, wCfg := range sCfg.Windows {
		if _, err := applyWindowCfg(ctx, runner, session, wCfg); err != nil {
			return fatalf(session, "applying window configuration: %w", err)
		}
	}

	if err := session.SelectActive(ctx); err != nil {
		return fatalf(session, "selecting active window: %w", err)
	}

	return session, nil
}

// fatalf is a helper function for [Apply] that constructs an error from
// provided format and args, and closes the provided tmux session if not nil.
func fatalf(sess *tmux.Session, format string, args ...any) (*tmux.Session, error) {
	err := fmt.Errorf(format, args...)

	if sess != nil {
		if closeErr := sess.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing failed session: %w", closeErr))
		}
	}

	return nil, err
}

// applyWindowCfg creates a new tmux window from the provided configuration on
// the provided tmux session.
func applyWindowCfg(ctx context.Context, r tmux.Runner, s *tmux.Session, cfg WindowConfig) (*tmux.Window, error) {
	win, err := tmux.NewWindow(r, s, makeWindowOpts(cfg)...)
	if err != nil {
		return nil, fmt.Errorf("creating window: %w", err)
	}

	if err := win.Apply(ctx); err != nil {
		return nil, fmt.Errorf("applying %s: %w", win, err)
	}

	for _, pCfg := range cfg.Panes {
		if _, err := applyPaneCfg(ctx, r, win, nil, pCfg); err != nil {
			return nil, err
		}
	}

	return win, nil
}

func applyPaneCfg(ctx context.Context, r tmux.Runner, w *tmux.Window, pp *tmux.Pane, cfg PaneConfig) (*tmux.Pane, error) { //nolint:revive // more readable in one line.
	pane, err := tmux.NewPane(r, w, pp, makePaneOpts(cfg)...)
	if err != nil {
		return nil, fmt.Errorf("creating pane: %w", err)
	}

	if err := pane.Apply(ctx); err != nil {
		return nil, fmt.Errorf("applying %s: %w", pane, err)
	}

	for _, pCfg := range cfg.Panes {
		if _, err := applyPaneCfg(ctx, r, w, pane, pCfg); err != nil {
			return nil, err
		}
	}

	return pane, nil
}

func makeSessionOpts(sCfg SessionConfig) []tmux.SessionOption {
	opts := []tmux.SessionOption{}

	if sCfg.Name != "" {
		opts = append(opts, tmux.SessionWithName(sCfg.Name))
	}

	if sCfg.Path != "" {
		opts = append(opts, tmux.SessionWithPath(sCfg.Path))
	}

	if sCfg.OnWindow != "" {
		opts = append(opts, tmux.SessionWithOnWindowCommand(sCfg.OnWindow))
	}

	if sCfg.OnPane != "" {
		opts = append(opts, tmux.SessionWithOnPaneCommand(sCfg.OnPane))
	}

	if sCfg.OnAny != "" {
		opts = append(opts, tmux.SessionWithOnAnyCommand(sCfg.OnAny))
	}

	if len(sCfg.Env) != 0 {
		opts = append(opts, tmux.SessionWithEnv(sCfg.Env))
	}

	return opts
}

func makeWindowOpts(wCfg WindowConfig) []tmux.WindowOption {
	opts := []tmux.WindowOption{tmux.WindowWithName(wCfg.Name)}

	if wCfg.Path != "" {
		opts = append(opts, tmux.WindowWithPath(wCfg.Path))
	}

	if wCfg.Command != "" {
		opts = append(opts, tmux.WindowWithCommands(wCfg.Command))
	}

	if len(wCfg.Commands) != 0 {
		opts = append(opts, tmux.WindowWithCommands(wCfg.Commands...))
	}

	if wCfg.Active {
		opts = append(opts, tmux.WindowAsActive())
	}

	if len(wCfg.Env) != 0 {
		opts = append(opts, tmux.WindowWithEnv(wCfg.Env))
	}

	return opts
}

func makePaneOpts(pCfg PaneConfig) []tmux.PaneOption {
	opts := []tmux.PaneOption{}

	if pCfg.Path != "" {
		opts = append(opts, tmux.PaneWithPath(pCfg.Path))
	}

	if pCfg.Size != "" {
		opts = append(opts, tmux.PaneWithSize(pCfg.Size))
	}

	if pCfg.Command != "" {
		opts = append(opts, tmux.PaneWithCommands(pCfg.Command))
	}

	if len(pCfg.Commands) != 0 {
		opts = append(opts, tmux.PaneWithCommands(pCfg.Commands...))
	}

	if pCfg.Horizontal {
		opts = append(opts, tmux.PaneWithHorizontalDirection())
	}

	if pCfg.Active {
		opts = append(opts, tmux.PaneAsActive())
	}

	if len(pCfg.Env) != 0 {
		opts = append(opts, tmux.PaneWithEnv(pCfg.Env))
	}

	return opts
}
