package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func Initialize() error {
	baseLogger, err := zap.NewProduction()

	if err != nil {
		return err
	}

	Log = baseLogger.Sugar()

	return nil
}
