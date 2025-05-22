package plogger

import (
	"os"
)

// ConsoleWriterOptions configures where to write logs (stdout or stderr).
type ConsoleWriterOptions struct {
	StdoutWriter *bool `default:"false"` // Write to stdout
	StderrWriter *bool `default:"true"`  // Write to stderr
}

// ConsoleWriterOptionsBuilder builds ConsoleWriterOptions.
type ConsoleWriterOptionsBuilder struct {
	Options []Setter[ConsoleWriterOptions]
}

// ConsoleWriter writes logs to stdout or stderr.
type ConsoleWriter struct {
	*ConsoleWriterOptions
	stdout Writer
	stderr Writer
}

// NewConsoleWriter creates a ConsoleWriter with given options.
func NewConsoleWriter(opts ...Lister[ConsoleWriterOptions]) (*ConsoleWriter, error) {
	args, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	return &ConsoleWriter{
		stdout:               os.Stdout,
		stderr:               os.Stderr,
		ConsoleWriterOptions: args,
	}, nil
}

// Write logs to stdout or stderr based on options.
func (w *ConsoleWriter) Write(p []byte) (n int, err error) {
	if w.StdoutWriter != nil && *w.StdoutWriter {
		return w.stdout.Write(p)
	}
	return w.stderr.Write(p)
}

// NewConsoleWriterOptions returns a new options builder.
func NewConsoleWriterOptions() *ConsoleWriterOptionsBuilder {
	return &ConsoleWriterOptionsBuilder{}
}

// List returns all configured setters.
func (b *ConsoleWriterOptionsBuilder) List() []Setter[ConsoleWriterOptions] {
	return b.Options
}

// SetStderr enables stderr, disables stdout.
func (b *ConsoleWriterOptionsBuilder) SetStderr(enable bool) *ConsoleWriterOptionsBuilder {
	b.Options = append(b.Options, func(cwo *ConsoleWriterOptions) error {
		notEnable := !enable
		cwo.StderrWriter = &enable
		cwo.StdoutWriter = &notEnable
		return nil
	})
	return b
}

// SetStdout enables stdout, disables stderr.
func (b *ConsoleWriterOptionsBuilder) SetStdout(enable bool) *ConsoleWriterOptionsBuilder {
	b.Options = append(b.Options, func(cwo *ConsoleWriterOptions) error {
		notEnable := !enable
		cwo.StdoutWriter = &enable
		cwo.StderrWriter = &notEnable
		return nil
	})
	return b
}
