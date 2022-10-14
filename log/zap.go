package log

import (
	"fmt"

	"go.uber.org/zap"
)

// var _ log.Logger = (*Logger)(nil)

type ZapLogger struct {
	log *zap.Logger
}

func NewZapLogger(zlog *zap.Logger) *ZapLogger {
	return &ZapLogger{zlog}
}

func (l *ZapLogger) Log(level Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case LevelDebug:
		l.log.Debug("", data...)
	case LevelInfo:
		l.log.Info("", data...)
	case LevelWarn:
		l.log.Warn("", data...)
	case LevelError:
		l.log.Error("", data...)
	case LevelFatal:
		l.log.Fatal("", data...)
	}
	return nil
}

func (l *ZapLogger) Sync() error {
	return l.log.Sync()
}

func (l *ZapLogger) Close() error {
	return l.Sync()
}