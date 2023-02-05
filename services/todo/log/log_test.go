package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogLevelGeneration(t *testing.T) {
	type test struct {
		createLevel string
		zapLevel    zapcore.Level
	}
	tests := []test{
		{
			createLevel: "debug",
			zapLevel:    zap.DebugLevel,
		},
		{
			createLevel: "info",
			zapLevel:    zap.InfoLevel,
		},
		{
			createLevel: "WARN",
			zapLevel:    zap.WarnLevel,
		},
		{
			createLevel: "ErRoR",
			zapLevel:    zap.ErrorLevel,
		},
		{
			createLevel: "some-invalid-value",
			zapLevel:    zap.InfoLevel,
		},
	}
	for _, c := range tests {
		log, err := New(c.createLevel)
		assert.Nil(t, err)
		assert.True(t, log.Core().Enabled(c.zapLevel))
	}
}

func TestSetDefaul(t *testing.T) {
	log, _ := New("debug")
	SetDefault(log)
	assert.Same(t, log, Default(), "set default should provide the logger at the pointer level")
}
