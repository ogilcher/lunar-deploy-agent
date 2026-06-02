package logger

import (
	"go.uber.org/zap"
)

// Log is the shared structured logger used across the agent.
var Log *zap.SugaredLogger

// Initialize configures the shared production logger.
func Initialize() error {
	baseLogger, err := zap.NewProduction()

	if err != nil {
		return err
	}

	Log = baseLogger.Sugar()

	return nil
}
