package logging

import (
	"fmt"
)

type decoratedLogger struct {
	logger Logger
	prefix string
	args   []interface{}
}

func (dl *decoratedLogger) Log(level LogLevel, text string, args ...interface{}) {
	passedArgs := append(dl.args, args...)
	dl.logger.Log(level, fmt.Sprintf("%s%s", dl.prefix, text), passedArgs...)
}

func (dl *decoratedLogger) MinLevel() LogLevel {
	return dl.logger.MinLevel()
}

func (dl *decoratedLogger) IsConcurrent() bool {
	return dl.logger.IsConcurrent()
}

func Decorate(logger Logger, prefix string, args ...interface{}) Logger {
	if logger == NilLogger {
		// optimization: keep logging nothing
		return logger
	}

	return &decoratedLogger{
		prefix: prefix,
		logger: logger,
		args:   args,
	}
}
