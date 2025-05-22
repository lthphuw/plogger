package plogger_test

import (
	"testing"

	"github.com/lthphuw/plogger"
)

func TestDiscardWriter(t *testing.T) {
	testcases := []struct {
		name         string
		input        []byte
		wantNewErr   bool
		wantWriteErr bool
		wantLength   int
	}{
		{
			name:         "Test Discard Writer",
			input:        []byte("Hello Discard Writer"),
			wantNewErr:   false,
			wantWriteErr: false,
			wantLength:   len("Hello Discard Writer"),
		},
		{
			name:         "Write Empty slice",
			input:        nil,
			wantNewErr:   false,
			wantWriteErr: false,
			wantLength:   0,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			writer, err := plogger.NewDiscardWriter()
			if (err != nil) != tt.wantNewErr {
				t.Errorf("TestDiscardWriter() error = %v, wantErr = %v", err, tt.wantNewErr)
				return
			}
			if err != nil {
				return
			}

			n, err := writer.Write(tt.input)
			if (err != nil) != tt.wantWriteErr {
				t.Errorf("Write() error = %v, wantErr = %v", err, tt.wantWriteErr)
			}
			if n != tt.wantLength {
				t.Errorf("Write() length = %d, want = %d", n, tt.wantLength)
			}
		})
	}
}
