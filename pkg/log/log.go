package log

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Initializes the log depending on the environment
func InitGlobal() error {
	var core zapcore.Core

	var logger *zap.Logger
	if strings.ToLower(os.Getenv("ENV")) == "prod" {
		zapconf := zap.NewProductionConfig()
		zapconf.EncoderConfig.FunctionKey = "func"
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(zapconf.EncoderConfig),
			zapcore.Lock(os.Stdout),
			zapcore.InfoLevel,
		)

		logger = zap.New(
			core,
			zap.AddCaller(),
		).Named("prod")
	} else {
		// if we are not in gcp use a console logger
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
		}

		logLevel := zapcore.InfoLevel
		debug := strings.ToLower(os.Getenv("DEBUG"))
		if debug == "1" || debug == "true" {
			logLevel = zapcore.DebugLevel
		}

		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(config.EncoderConfig),
			zapcore.Lock(os.Stdout),
			logLevel,
		)

		logger = zap.New(
			core,
			zap.AddCaller(),
		).Named("dev")
	}

	zap.ReplaceGlobals(logger)
	return nil
}
