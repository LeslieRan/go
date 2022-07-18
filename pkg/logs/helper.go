package logs

import (
	"context"
)

func Info(format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Info(msg, fields...)
}

func InfoCtx(ctx context.Context, format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Ctx(ctx).Info(msg, fields...)
}

func Debug(format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Debug(msg, fields...)
}

func DebugCtx(ctx context.Context, format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Ctx(ctx).Debug(msg, fields...)
}

func Warn(format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Warn(msg, fields...)
}

func WarnCtx(ctx context.Context, format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Ctx(ctx).Warn(msg, fields...)
}

func Error(format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Error(msg, fields...)
}

func ErrorCtx(ctx context.Context, err error, format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Err(err).Ctx(ctx).Error(msg, fields...)
}

func Panic(format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Panic(msg, fields...)
}

func PanicCtx(ctx context.Context, format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Ctx(ctx).Panic(msg, fields...)
}

func Fatal(format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Fatal(msg, fields...)
}

func FatalCtx(ctx context.Context, format string, args ...interface{}) {
	f := GetFactory()
	msg, fields := f.parse(format, args...)

	f.logger.Ctx(ctx).Fatal(msg, fields...)
}
