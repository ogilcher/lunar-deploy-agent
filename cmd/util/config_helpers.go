package util

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// LoadAndValidateConfig loads the agent config and validates it before
// command-specific logic runs.
//
// This keeps Cobra commands small and prevents duplicate config-loading logic.
func LoadAndValidateConfig(
	configPath string,
) *config.Config {
	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Log.Errorw("Failed to load config.", "error", err)
		return nil
	}

	if err := config.ValidateConfig(appConfig); err != nil {
		logger.Log.Errorw("Invalid deployment config.", "error", err)
		return nil
	}

	return appConfig
}
