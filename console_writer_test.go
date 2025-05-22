package plogger_test

import (
	"testing"

	"github.com/lthphuw/plogger"
)

func TestConsoleWriter(t *testing.T) {
	testcases := []struct {
		name         string
		opts         *plogger.ConsoleWriterOptionsBuilder
		input        []byte
		wantNewErr   bool
		wantWriteErr bool
		wantLength   int
	}{
		{
			name:         "No options - default (write to stderr)",
			opts:         nil,
			input:        []byte("Hello stderr"),
			wantWriteErr: false,
			wantNewErr:   false,
			wantLength:   len("Hello stderr"),
		},
		{
			name:         "Write Empty slice",
			opts:         nil,
			input:        nil,
			wantNewErr:   false,
			wantWriteErr: false,
			wantLength:   0,
		},
		{
			name:         "Write to stdout only",
			opts:         plogger.NewConsoleWriterOptions().SetStdout(true),
			input:        []byte("Hello stdout"),
			wantNewErr:   false,
			wantWriteErr: false,
			wantLength:   len("Hello stdout"),
		},
		{
			name:         "Write to stderr only",
			opts:         plogger.NewConsoleWriterOptions().SetStderr(true),
			input:        []byte("Hello stderr again"),
			wantNewErr:   false,
			wantWriteErr: false,
			wantLength:   len("Hello stderr again"),
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			writer, err := plogger.NewConsoleWriter(tt.opts)
			if (err != nil) != tt.wantNewErr {
				t.Errorf("NewConsoleWriter() error = %v, wantNewErr = %v", err, tt.wantNewErr)
				return
			}
			if err != nil {
				return
			}

			n, err := writer.Write(tt.input)
			if (err != nil) != tt.wantWriteErr {
				t.Errorf("Write() error = %v, wantWriteErr = %v", err, tt.wantWriteErr)
			}
			if n != tt.wantLength {
				t.Errorf("Write() length = %d, want = %d", n, tt.wantLength)
			}
		})
	}
}
