/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0

Refactored: 1
*/

package logging

import (
	"fmt"
	"io"
	"math"
	"os"
)

// Logger is minimal logging interface designed to be easily adaptable to any
// logging library.
type Logger interface {
	// Log is invoked with the log level, the log message, and key/value pairs
	// of any relevant log details. The keys are always strings, while the
	// values are unspecified.
	Log(level LogLevel, text string, args ...interface{})

	// MinLevel returns the minimum level at which the logging is performed.
	// This method must be thread-safe and wait-free even if IsConcurrent returns false.
	MinLevel() LogLevel

	// IsConcurrent returns true iff the logger can be safely concurrently accessed by multiple go-routines.
	IsConcurrent() bool
}

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelDisable // not supposed to be passed to the Log method.
)

// Simple logger writing log messages directly to an I/O writer.
type streamLogger struct {
	level  LogLevel
	writer io.Writer
}

// Log is invoked with the log level, the log message, and key/value pairs
// of any relevant log details. The keys are always strings, while the
// values are unspecified. If the level is greater of equal than this consoleLogger,
// Log() writes the log message to standard output.
func (l streamLogger) Log(level LogLevel, text string, args ...interface{}) {
	if level < l.level {
		return
	}

	fmt.Print(text)
	for i := 0; i < len(args); i++ {
		if i+1 < len(args) {
			switch args[i+1].(type) {
			case []byte:
				// Print byte arrays in base 16 encoding.
				fmt.Fprintf(l.writer, " %s=%x", args[i], args[i+1])
			default:
				// Print all other types using the Go default format.
				fmt.Fprintf(l.writer, " %s=%v", args[i], args[i+1])
			}
			i++
		} else {
			fmt.Fprintf(l.writer, " %s=%%MISSING%%", args[i])
		}
	}
	fmt.Fprintf(l.writer, "\n")
}

func (l streamLogger) MinLevel() LogLevel {
	return l.level
}

func (l streamLogger) IsConcurrent() bool {
	return false
}

// The nil logger drops all messages.
type nilLogger struct{}

// The Log method of the nilLogger does nothing, effectively dropping every log message.
func (nl nilLogger) Log(level LogLevel, text string, args ...interface{}) {
	// Do nothing.
}

func (nl nilLogger) MinLevel() LogLevel {
	return LevelDisable
}

func (nl nilLogger) IsConcurrent() bool {
	return true
}

func NewStreamLogger(level LogLevel, writer io.Writer) Logger {
	return &streamLogger{
		level:  level,
		writer: writer,
	}
}

var (
	// ConsoleTraceLogger implements Logger and writes all log messages to stdout.
	ConsoleTraceLogger = Synchronize(streamLogger{LevelTrace, os.Stdout})

	// ConsoleDebugLogger implements Logger and writes all LevelDebug and above messages to stdout.
	ConsoleDebugLogger = Synchronize(streamLogger{LevelDebug, os.Stdout})

	// ConsoleInfoLogger implements Logger and writes all LevelInfo and above log messages to stdout.
	ConsoleInfoLogger = Synchronize(streamLogger{LevelInfo, os.Stdout})

	// ConsoleWarnLogger implements Logger and writes all LevelWarn and above log messages to stdout.
	ConsoleWarnLogger = Synchronize(streamLogger{LevelWarn, os.Stdout})

	// ConsoleErrorLogger implements Logger and writes all LevelError log messages to stdout.
	ConsoleErrorLogger = Synchronize(streamLogger{LevelError, os.Stdout})

	// NilLogger drops all log messages.
	NilLogger = nilLogger{}
)

type multiLogger struct {
	loggers []Logger
}

func NewMultiLogger(loggers []Logger) Logger {
	return &multiLogger{loggers: loggers}
}

func (ml *multiLogger) Log(level LogLevel, text string, args ...interface{}) {
	for _, l := range ml.loggers {
		l.Log(level, text, args...)
	}
}

func (ml *multiLogger) MinLevel() LogLevel {
	min := LogLevel(math.MaxInt)
	for _, l := range ml.loggers {
		level := l.MinLevel()
		if level < min {
			min = level
		}
	}

	return min
}

func (ml *multiLogger) IsConcurrent() bool {
	res := true
	for _, l := range ml.loggers {
		res = res && l.IsConcurrent()
	}
	return res
}
