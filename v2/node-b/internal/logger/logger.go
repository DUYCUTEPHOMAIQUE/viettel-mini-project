package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	writer := &zapcore.BufferedWriteSyncer{
		WS:   zapcore.AddSync(os.Stdout),
		Size: 256 * 1024,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		writer,
		zapcore.InfoLevel,
	)
	return zap.New(core)
}
