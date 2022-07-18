package logs

import (
	"context"

	"go.uber.org/zap"
)

const defaultNameSpace = "content"

var defaultOptions = &Options{
	namespace: defaultNameSpace,
	errHandler: func(err error) []zap.Field {
		return []zap.Field{zap.Error(err)}
	},
}

type Options struct {
	namespace  string
	ctxHandler func(ctx context.Context) []zap.Field
	errHandler func(err error) []zap.Field
}

type Option func(options *Options)

func WithNamespace(ns string) Option {
	return func(options *Options) {
		options.namespace = ns
	}
}

func WithCtxHandler(fn func(ctx context.Context) []zap.Field) Option {
	return func(options *Options) {
		options.ctxHandler = fn
	}
}

func WithErrHandler(fn func(err error) []zap.Field) Option {
	return func(options *Options) {
		options.errHandler = fn
	}
}
