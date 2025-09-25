package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Init(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Create logger
	logger, err := config.Build()
	if err != nil {
		return err
	}

	Logger = logger
	return nil
}

func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	} else {
		// Fallback if logger is not initialized
		panic(msg)
	}
}

func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}
