package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

func InitLogger() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	file, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(file)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, zapcore.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	Log = zap.New(core, zap.AddCaller())
	Sugar = Log.Sugar()
}

func Close() {
	Log.Sync()
	Sugar.Sync()
}
