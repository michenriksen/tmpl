package cli

import (
	"context"
	"fmt"

	"github.com/michenriksen/tmpl/config"
	"github.com/michenriksen/tmpl/tmux"
)

// runApply loads the configuration, applies it to a new tmux session and
// attaches it.
func (a *App) runApply(ctx context.Context) error {
	a.initLogger()

	if err := a.loadConfig(); err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	runner, err := a.newTmux()
	if err != nil {
		return fmt.Errorf("creating tmux runner: %w", err)
	}

	a.sess, err = config.Apply(ctx, a.cfg, runner)
	if err != nil {
		return fmt.Errorf("applying configuration: %w", err)
	}

	if err := a.sess.Attach(ctx); err != nil {
		return fmt.Errorf("attaching session: %w", err)
	}

	return nil
}

func (a *App) newTmux() (tmux.Runner, error) {
	if a.tmux != nil {
		return a.tmux, nil
	}

	cmdOpts := []tmux.RunnerOption{tmux.WithLogger(a.logger)}

	if a.cfg.Tmux != "" {
		cmdOpts = append(cmdOpts, tmux.WithTmux(a.cfg.Tmux))
	}

	if len(a.cfg.TmuxOptions) > 0 {
		cmdOpts = append(cmdOpts, tmux.WithTmuxOptions(a.cfg.TmuxOptions...))
	}

	if a.opts.DryRun {
		a.logger.Info("DRY-RUN MODE ENABLED: no tmux commands will be executed and output is simulated")

		cmdOpts = append(cmdOpts, tmux.WithDryRunMode(true))
	}

	cmd, err := tmux.NewRunner(cmdOpts...)
	if err != nil {
		return nil, err //nolint:wrapcheck // Wrapping is done by caller.
	}

	return cmd, nil
}
