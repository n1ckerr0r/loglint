package logtest

import "go.uber.org/zap"

func zapLoggerExample() {
	logger, _ := zap.NewProduction()
	logger.Info("Bad Message") // want "lowercase"
}
