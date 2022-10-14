package log

import (
	"testing"

	"go.uber.org/zap"
)

type service struct{
	log *Helper
}

func TestLoggerHelper(t *testing.T) {
	logger, _ := zap.NewProduction()
	zapLog := NewZapLogger(logger)
	s := service{
		log: NewHelper(zapLog),
	}
	s.log.Log(LevelInfo, "this is a test log")
}

func TestGlobalLog(t *testing.T) {
	Info("this is a test log")
}
