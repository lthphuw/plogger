package plogger_test

import (
	"testing"
	"time"

	"github.com/lthphuw/plogger"
)

func TestTextFormatter(t *testing.T) {
	testcases := []struct {
		name            string
		opts            plogger.Lister[plogger.TextFormatterOptions]
		timestamp       time.Time
		entry           *plogger.Entry
		wantNewError    bool
		wantFormatError bool
		wantBytes       []byte
	}{
		{
			name:      "Test TextFormatter",
			opts:      plogger.NewTextFormatterOptions(),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level:     plogger.InfoLevel,
				Timestamp: time.Now(),
				Msg:       "Hello world",
			},
			wantBytes: []byte(
				"timestamp=" + time.Now().
					Format("2006-01-02T15:04:05Z07:00") +
					" level=info" + " msg=Hello world" + "\n",
			),
		},
		{
			name:      "Test TextFormatter disable color",
			opts:      plogger.NewTextFormatterOptions().SetDisableColor(true),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level: plogger.FatalLevel,
				Msg:   "Hello world",
			},
			wantBytes: []byte(
				"timestamp=" + time.Now().
					Format("2006-01-02T15:04:05Z07:00") +
					" level=fatal" + " msg=Hello world" + "\n",
			),
		},

		{
			name:      "Test TextFormatter with sorting keys",
			opts:      plogger.NewTextFormatterOptions().SetDisableSorting(false),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level:     plogger.WarnLevel,
				Timestamp: time.Now(),
				Msg:       "Hello world",
				FieldMap: map[string]any{
					"key2": 1,
					"key1": 2,
					"key3": 3,
					"key0": nil,
				},
			},
			wantBytes: []byte(
				"timestamp=" + time.Now().
					Format("2006-01-02T15:04:05Z07:00") +
					" level=warn" + " msg=Hello world" + " key0=nil key1=2 key2=1 key3=3" + "\n",
			),
		},
		{
			name:      "Test TextFormatter without sorting keys",
			opts:      plogger.NewTextFormatterOptions().SetDisableSorting(false),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level:     plogger.TraceLevel,
				Timestamp: time.Now(),
				Msg:       "Hello world",
				FieldMap: map[string]any{
					"key2": 1,
					"key1": 2,
					"key3": 3,
				},
			},
			wantBytes: []byte(
				"timestamp=" + time.Now().
					Format("2006-01-02T15:04:05Z07:00") +
					" level=trace" + " msg=Hello world" + " key1=2 key2=1 key3=3" + "\n",
			),
		},
		{
			name: "Test TextFormatter with custom sorting keys",
			opts: plogger.NewTextFormatterOptions().SetDisableSorting(false).SetSortingFieldKeyFunc(
				func(s []string) {
					s[0] = "key3"
					s[1] = "key1"
					s[2] = "key2"
				}),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level:     plogger.DebugLevel,
				Timestamp: time.Now(),
				Msg:       "Hello world",
				FieldMap: map[string]any{
					"key2": 1,
					"key1": 2,
					"key3": 3,
				},
			},
			wantBytes: []byte(
				"timestamp=" + time.Now().
					Format("2006-01-02T15:04:05Z07:00") +
					" level=debug" + " msg=Hello world" + " key3=3 key1=2 key2=1" + "\n",
			),
		},
		{
			name:      "Test TextFormatter with format timestamp",
			opts:      plogger.NewTextFormatterOptions().SetTimestampFormat("2006-01-02_15-04-05"),
			timestamp: time.Now(),
			entry: &plogger.Entry{
				Level:     plogger.ErrorLevel,
				Timestamp: time.Now(),
				Msg:       "Hello world",
			},
			wantBytes: []byte(
				"timestamp=" + time.Now().
					Format("2006-01-02_15-04-05") +
					" level=error" + " msg=Hello world" + "\n",
			),
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			f, err := plogger.NewTextFormatter(tt.opts)
			if (err != nil) != tt.wantNewError {
				t.Errorf("NewTextFormatter() got error %v, want %v", err, tt.wantNewError)
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
