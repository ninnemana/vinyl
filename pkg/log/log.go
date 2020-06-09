package log

import (
	"os"
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	timeFormat   = "2006-01-02T15:04:05Z07:00"
	defaultLevel = -1
)

var (
	customTimeFormat string
	Logger           *zap.Logger
)

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(customTimeFormat))
}

func Init() (*zap.Logger, error) {
	globalLevel := zapcore.Level(defaultLevel)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	// It is usefull for Kubernetes deployment.
	// Kubernetes interprets os.Stdout log items as INFO and os.Stderr log items
	// as ERROR by default.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= globalLevel && lvl < zapcore.ErrorLevel
	})
	consoleInfos := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Configure console output.
	var useCustomTimeFormat bool
	ecfg := zap.NewProductionEncoderConfig()
	if len(timeFormat) > 0 {
		customTimeFormat = timeFormat
		ecfg.EncodeTime = customTimeEncoder
		useCustomTimeFormat = true
	}
	consoleEncoder := zapcore.NewJSONEncoder(ecfg)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleInfos, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	Logger = zap.New(core)
	zap.RedirectStdLog(Logger)

	if !useCustomTimeFormat {
		Logger.Warn("time format for logger is not provided - use zap default")
	}

	grpc_zap.ReplaceGrpcLogger(Logger)

	return Logger, nil
}
