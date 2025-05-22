package plogger_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lthphuw/plogger"
)

func TestJSONFormatter(t *testing.T) {
	testcases := []struct {
		name            string
		timestamp       time.Time
		opts            plogger.Lister[plogger.JSONFormatterOptions]
		entry           *plogger.Entry
		wantBytes       []byte
		wantNewError    bool
		wantFormatError bool
	}{
		{
			name: "Test JSON formatter",
			opts: plogger.NewJSONFormatterOptions().
				SetPrettyPrint(false).
				SetEscapeHTML(false).
				SetTimestampFormat("2006-01-02T15:04:05Z07:00"),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level:     plogger.InfoLevel,
				Timestamp: time.Now(),
				Msg:       "Hello world",
			},
			wantBytes: fmt.Appendf(
				nil,
				`{"level":"info","msg":"Hello world","timestamp":"%s"}`+"\n",
				time.Now().Format("2006-01-02T15:04:05Z07:00"),
			),
		},

		{
			name: "Test JSON formatter no timestamp",
			opts: plogger.NewJSONFormatterOptions().
				SetPrettyPrint(false).
				SetEscapeHTML(false).
				SetTimestampFormat("2006-01-02T15:04:05Z07:00"),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level: plogger.InfoLevel,
				Msg:   "Hello world",
			},
			wantBytes: fmt.Appendf(
				nil,
				`{"level":"info","msg":"Hello world","timestamp":"%s"}`+"\n",
				time.Now().Format("2006-01-02T15:04:05Z07:00"),
			),
		},
		{
			name: "No escape HTML",
			opts: plogger.NewJSONFormatterOptions().
				SetPrettyPrint(false).
				SetEscapeHTML(false).
				SetTimestampFormat("2006-01-02T15:04:05Z07:00"),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level:     plogger.InfoLevel,
				Timestamp: time.Now(),
				Msg:       "Hello world ><",
			},
			wantBytes: fmt.Appendf(
				nil,
				`{"level":"info","msg":"Hello world ><","timestamp":"%s"}`+"\n",
				time.Now().Format("2006-01-02T15:04:05Z07:00"),
			),
		},
		{
			name: "Escape HTML",
			opts: plogger.NewJSONFormatterOptions().
				SetPrettyPrint(false).
				SetEscapeHTML(true).
				SetTimestampFormat("2006-01-02T15:04:05Z07:00"),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level:     plogger.InfoLevel,
				Timestamp: time.Now(),
				Msg:       "Hello world ><",
			},
			wantBytes: fmt.Appendf(
				nil,
				`{"level":"info","msg":"Hello world \u003e\u003c","timestamp":"%s"}`+"\n",
				time.Now().Format("2006-01-02T15:04:05Z07:00"),
			),
		},
		{
			name: "Set format timestamp",
			opts: plogger.NewJSONFormatterOptions().
				SetPrettyPrint(false).
				SetEscapeHTML(true).
				SetTimestampFormat("2006-01-02_15-04-05"),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level:     plogger.InfoLevel,
				Timestamp: time.Now(),
				Msg:       "Hello world ><",
			},
			wantBytes: fmt.Appendf(
				nil,
				`{"level":"info","msg":"Hello world \u003e\u003c","timestamp":"%s"}`+"\n",
				time.Now().Format("2006-01-02_15-04-05"),
			),
		},
		{
			name: "Pretty print",
			opts: plogger.NewJSONFormatterOptions().
				SetPrettyPrint(true).
				SetEscapeHTML(false).
				SetTimestampFormat("2006-01-02T15:04:05Z07:00"),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level:     plogger.InfoLevel,
				Timestamp: time.Now(),
				Msg:       "Hello world",
			},
			wantBytes: fmt.Appendf(nil,
				`{
  "level": "info",
  "msg": "Hello world",
  "timestamp": "%s"
}`+"\n", time.Now().Format("2006-01-02T15:04:05Z07:00")),
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			f, err := plogger.NewJSONFormatter(tt.opts)
			if (err != nil) != tt.wantNewError {
				t.Errorf("NewJSONFormatter() got error %v, want %v", err, tt.wantNewError)
			}
			b, err := f.Format(tt.entry)
			if (err != nil) != tt.wantFormatError {
				t.Errorf("Format() got error %v, want %v", err, tt.wantFormatError)
			}
			if string(b) != string(tt.wantBytes) {
				t.Errorf("Format() got bytes %v, want %v", string(b), string(tt.wantBytes))
			}
		})
	}
}
