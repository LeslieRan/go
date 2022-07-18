package logs

import (
	"context"
	goLog "log"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type logger struct {
	factory   *Factory
	zapLogger *zap.Logger
}

func newLogger(cfg *Config) *logger {
	if cfg == nil {
		goLog.Panicln("configurations must be given")
	}

	var zapLogger *zap.Logger
	switch cfg.Output {
	case OutputFile:
		if cfg.Lumberjack == nil {
			goLog.Panicln("lumberjack configuration must be given")
		}
		zapLogger = newFileLogger(cfg)
	case OutputConsole:
		zapLogger = newConsoleLogger(cfg)
	default:
		goLog.Panicf("unsupported output way: %v", cfg.Output)
	}

	return &logger{
		zapLogger: zapLogger,
	}
}

func newFileLogger(cfg *Config) *zap.Logger {
	writer := zapcore.AddSync(cfg.Lumberjack)
	lv := level(cfg.Level)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(newFileEncoderConfig()), writer, lv)

	return zap.New(core, zap.WithCaller(true), zap.AddCallerSkip(2))
}

func newConsoleLogger(cfg *Config) *zap.Logger {
	lv := level(cfg.Level)
	loggerConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(lv),
		Encoding:         "console",
		EncoderConfig:    newConsoleEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := loggerConfig.Build(zap.AddCallerSkip(2))
	if err != nil {
		goLog.Panicf("new console logger: %s", err.Error())
	}

	return logger
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *logger) Panic(msg string, fields ...zap.Field) {
	l.zapLogger.Panic(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *logger) With(fields ...zap.Field) *logger {
	if len(fields) == 0 {
		return l
	}

	ll := l.clone()
	ll.zapLogger = l.zapLogger.With(fields...)
	ll.factory = l.factory

	return ll
}

func (l *logger) Ctx(ctx context.Context) *logger {
	var (
		opts   = l.factory.opts
		fields []zap.Field
	)

	if opts.ctxHandler != nil {
		fields = append(fields, opts.ctxHandler(ctx)...)
	}

	return l.With(fields...)
}

func (l *logger) Err(err error) *logger {
	var (
		opts   = l.factory.opts
		fields []zap.Field
	)

	if opts.errHandler != nil {
		fields = append(fields, opts.errHandler(err)...)
	}

	return l.With(fields...)
}

func (l *logger) clone() *logger {
	c := *l

	return &c
}

func level(lv string) zapcore.Level {
	var zapLv zapcore.Level
	switch strings.ToLower(lv) {
	case "info":
		zapLv = zapcore.InfoLevel
	case "debug":
		zapLv = zapcore.DebugLevel
	case "warn":
		zapLv = zapcore.WarnLevel
	case "error":
		zapLv = zapcore.ErrorLevel
	case "panic":
		zapLv = zapcore.PanicLevel
	case "fatal":
		zapLv = zapcore.FatalLevel
	default:
		zapLv = zapcore.InfoLevel
	}

	return zapLv
}

func newFileEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newConsoleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
