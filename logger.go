// Package plogger provides a lightweight and efficient logging framework.
package plogger

import (
	"errors"
	"os"
	"sync"
)

// LoggerOptions configures the logger's behavior.
type LoggerOptions struct {
	Writer    Writer    // Destination for log output
	Formatter Formatter // Format for log entries
	Level     Level     // Minimum log level to output
	Hooks     LevelHook // Hooks to run on log events
	ExitFunc  func(int) // Function to call on fatal logs (e.g., os.Exit)
	Caller    *bool     `default:"true"` // Whether to include caller info (default: true)
}

// LoggerOptionsBuilder builds and applies a list of setters for LoggerOptions.
type LoggerOptionsBuilder struct {
	Options []Setter[LoggerOptions]
}

// Logger logs
type Logger struct {
	*LoggerOptions // LoggerOptions holds configuration for logger.

	mu sync.Mutex // Mutex
}

// NewLogger creates a new logger
func NewLogger(opts ...Lister[LoggerOptions]) (*Logger, error) {
	args, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	// Set default exit function
	if args.ExitFunc == nil {
		args.ExitFunc = os.Exit
	}

	// Set default Writer (console)
	if args.Writer == nil {
		args.Writer, err = NewDiscardWriter()
		if err != nil {
			return nil, err
		}
	}

	// Set default Formatter (text)
	if args.Formatter == nil {
		args.Formatter, err = NewTextFormatter()
		if err != nil {
			return nil, err
		}
	}
	return &Logger{
		LoggerOptions: args,
	}, nil
}

// Log writes the log entry at the specified level if enabled.
func (l *Logger) Log(level Level, entry *Entry) error {
	if entry == nil {
		return errors.New("entry cant be nil")
	}
	if l.shouldLog(level) {
		l.mu.Lock()
		defer l.mu.Unlock()
		entry.Level = level
		entry.caller = *l.Caller
		if err := l.runHooks(level, entry); err != nil {
			return err
		}
		bytes, err := l.Formatter.Format(entry)
		if err != nil {
			return err
		}

		_, err = l.Writer.Write(bytes)

		return err
	}
	return nil
}

// Trace logs an entry at TraceLevel.
func (l *Logger) Trace(entry *Entry) error {
	return l.Log(TraceLevel, entry)
}

// Debug logs an entry at DebugLevel.
func (l *Logger) Debug(entry *Entry) error {
	return l.Log(DebugLevel, entry)
}

// Info logs an entry at InfoLevel.
func (l *Logger) Info(entry *Entry) error {
	return l.Log(InfoLevel, entry)
}

// Warn logs an entry at WarnLevel.
func (l *Logger) Warn(entry *Entry) error {
	return l.Log(WarnLevel, entry)
}

// Error logs an entry at ErrorLevel.
func (l *Logger) Error(entry *Entry) error {
	return l.Log(ErrorLevel, entry)
}

// Fatal logs an entry at FatalLevel and exits the application.
func (l *Logger) Fatal(entry *Entry) error {
	err := l.Log(FatalLevel, entry)
	if err == nil {
		l.ExitFunc(1)
	}
	return err
}

// runHooks runs all the hooks
func (l *Logger) runHooks(level Level, entry *Entry) error {
	if l.Hooks != nil {
		if err := l.Hooks.Run(level, entry); err != nil {
			return err
		}
	}
	return nil
}

// shouldLog checks if level passes the logger's threshold.
func (l *Logger) shouldLog(level Level) bool {
	return level >= l.Level
}

// NewLoggerOptions creates a new *LoggerOptionsBuilder
func NewLoggerOptions() *LoggerOptionsBuilder {
	return &LoggerOptionsBuilder{}
}

// List lists all the setter function for Logger
func (b *LoggerOptionsBuilder) List() []Setter[LoggerOptions] {
	return b.Options
}

// SetWriter sets writer (with Writer interface)
func (b *LoggerOptionsBuilder) SetWriter(writer Writer) *LoggerOptionsBuilder {
	b.Options = append(b.Options, func(lo *LoggerOptions) error {
		lo.Writer = writer
		return nil
	})
	return b
}

// SetFormatter sets formatter (with Formatter interface)
func (b *LoggerOptionsBuilder) SetFormatter(formatter Formatter) *LoggerOptionsBuilder {
	b.Options = append(b.Options, func(lo *LoggerOptions) error {
		lo.Formatter = formatter
		return nil
	})
	return b
}

// SetLevel sets the level that logger should log (or above)
func (b *LoggerOptionsBuilder) SetLevel(level Level) *LoggerOptionsBuilder {
	b.Options = append(b.Options, func(lo *LoggerOptions) error {
		lo.Level = level

		return nil
	})
	return b
}

// SetHooks adds hooks to logger
func (b *LoggerOptionsBuilder) SetHooks(hook Hook) *LoggerOptionsBuilder {
	b.Options = append(b.Options, func(lo *LoggerOptions) error {
		if lo.Hooks == nil {
			lo.Hooks = make(LevelHook)
		}
		lo.Hooks.Add(hook)

		return nil
	})
	return b
}

// SetExitFunc sets custom exit func when Fatal is called
func (b *LoggerOptionsBuilder) SetExitFunc(exitFunc func(int)) *LoggerOptionsBuilder {
	b.Options = append(b.Options, func(lo *LoggerOptions) error {
		lo.ExitFunc = exitFunc
		return nil
	})
	return b
}

// SetCaller sets enable caller
func (b *LoggerOptionsBuilder) SetCaller(enable bool) *LoggerOptionsBuilder {
	b.Options = append(b.Options, func(lo *LoggerOptions) error {
		lo.Caller = &enable
		return nil
	})
	return b
}
