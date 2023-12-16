package cli

import (
	"context"
	"fmt"
)

// runCheck loads the configuration and validates it, logging any validation
// errors and returning [ErrInvalidConfig] if any are found.
func (a *App) runCheck(_ context.Context) error {
	a.initLogger()

	if err := a.loadConfig(); err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	err := a.cfg.Validate()
	if err != nil {
		return err
	}

	a.logger.Info("configuration file is valid")

	return nil
}
