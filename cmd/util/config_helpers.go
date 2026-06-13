package util

import (
	"fmt"

	"github.com/ogilcher/lunar-deploy-agent/internal/config"
)

// LoadAndValidateConfig loads a YAML config file and validates it before use.
//
// It returns the loaded config when successful.
// It returns a wrapped error when loading or validation fails so callers can
// log or display a clear failure reason.
func LoadAndValidateConfig(
	configPath string,
) (*config.Config, error) {
	appConfig, err := config.LoadConfig(configPath)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to load config %q: %w",
			configPath,
			err,
		)
	}

	if err := config.ValidateConfig(appConfig); err != nil {
		return nil, fmt.Errorf(
			"invalid config %q: %w",
			configPath,
			err,
		)
	}

	return appConfig, nil
}
