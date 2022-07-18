package test_test

import (
	"context"
	"fmt"
	"testing"

	"gitbub.com/LeslieRan/go/pkg/logs"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

var ctx = context.Background()

func TestConsole(t *testing.T) {
	logs.Debug("hello: %s", "world")
	logs.DebugCtx(ctx, "hello: %s", "world")
}

func TestFile(t *testing.T) {
	cfg := &logs.Config{
		Level:  "info",
		Output: logs.OutputFile,
		Lumberjack: &lumberjack.Logger{
			Filename:   "./app.log",
			MaxAge:     1,
			MaxSize:    20,
			MaxBackups: 3,
			Compress:   true,
		},
	}
	factory := logs.NewFactory(cfg)
	logs.SetFactory(factory)

	logs.ErrorCtx(ctx, fmt.Errorf("error"), "hello: %s", "world", zap.String("string", "zap.field"))
}
