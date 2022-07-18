package logs

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var (
	once     sync.Once
	_factory *Factory
)

func init() {
	_factory = NewFactory(&Config{
		Level:  zap.DebugLevel.String(),
		Output: OutputConsole,
	})
}

type Factory struct {
	logger *logger
	opts   *Options
}

func NewFactory(cfg *Config, opts ...Option) *Factory {
	defaultOpts := defaultOptions

	for _, opt := range opts {
		opt(defaultOpts)
	}

	logger := newLogger(cfg)

	factory := &Factory{
		logger: logger,
		opts:   defaultOpts,
	}

	logger.factory = factory

	return factory
}

func (f *Factory) parse(format string, args ...interface{}) (string, []zap.Field) {
	var (
		fields []zap.Field
		a      []interface{}
		msg    string = format
	)

	for _, arg := range args {
		if field, ok := arg.(zap.Field); ok {
			fields = append(fields, field)
		} else {
			a = append(a, arg)
		}
	}

	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}

	if len(fields) > 0 {
		fields = append([]zap.Field{zap.Namespace(f.opts.namespace)}, fields...)
	}

	return msg, fields
}

func SetFactory(factory *Factory) {
	once.Do(func() {
		_factory = factory
	})
}

func GetFactory() *Factory {
	return _factory
}
