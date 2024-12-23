package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewChiLoggerWithZap(log *zap.Logger) middleware.LogFormatter {
	return &ChiZapLogger{log: log}
}

type ChiZapLogger struct {
	log *zap.Logger
}

func (l *ChiZapLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &ChiZapLogEntry{
		log: l.log,
		req: r,
	}
}

type ChiZapLogEntry struct {
	log *zap.Logger
	req *http.Request
}

func (l *ChiZapLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.log.Info("request completed",
		zap.String("method", l.req.Method),
		zap.String("path", l.req.URL.Path),
		zap.Int("status", status),
		zap.Int("bytes", bytes),
		zap.Duration("elapsed", elapsed),
		zap.String("ip", l.req.RemoteAddr),
	)
}

func (l *ChiZapLogEntry) Panic(v interface{}, stack []byte) {
	l.log.Error("request panic",
		zap.Any("panic", v),
		zap.String("stack", string(stack)),
	)
}
