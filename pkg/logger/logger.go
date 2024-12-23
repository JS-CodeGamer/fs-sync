package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

func createLogger(fname string) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	err := os.MkdirAll(filepath.Dir(fname), 0750)
	if err != nil {
		panic("error creating logfile")
	}
	file, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("error creating logfile")
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zapcore.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	return zap.New(core, zap.AddCaller())
}

func InitLogger(logPath string) {
	Logger = createLogger(logPath)
	Sugar = Logger.Sugar()
}

func Close() {
	Logger.Sync()
	Sugar.Sync()

	Logger.Sync()
	Sugar.Sync()
}
