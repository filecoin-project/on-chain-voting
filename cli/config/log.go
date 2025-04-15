package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// InitLogger initializes the logger with custom configurations.
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
}

// getEncoder creates and returns a configured encoder for the logger.
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter returns a WriteSyncer that writes logs to os.Stdout.
func getLogWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}
