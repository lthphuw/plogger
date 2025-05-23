package plogger_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lthphuw/plogger"
)

func TestLogger(t *testing.T) {
	testcases := []struct {
		name             string
		opts             plogger.Lister[plogger.LoggerOptions]
		level            plogger.Level
		entry            *plogger.Entry
		useJSONFormatter bool
		useTextFormatter bool
		wantNewErr       bool
		wantLogErr       bool
	}{
		{
			name:       "Empty entry",
			opts:       nil,
			entry:      nil,
			level:      plogger.InfoLevel,
			wantLogErr: true,
		},
		{
			name: "log trace",
			opts: nil,
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.TraceLevel,
			wantLogErr: false,
		},
		{
			name: "log debug",
			opts: nil,
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.DebugLevel,
			wantLogErr: false,
		},
		{
			name: "log info",
			opts: nil,
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.InfoLevel,
			wantLogErr: false,
		},
		{
			name: "log warn",
			opts: nil,
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.WarnLevel,
			wantLogErr: false,
		},
		{
			name: "log error",
			opts: nil,
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.ErrorLevel,
			wantLogErr: false,
		},
		{
			name: "log fatal",
			opts: plogger.NewLoggerOptions().SetExitFunc(func(i int) {}),
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.FatalLevel,
			wantLogErr: false,
		},
		{
			name: "with hooks (run only on info & debug), error on info",
			opts: plogger.NewLoggerOptions().SetHooks(&hookMock{}),
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.InfoLevel,
			wantLogErr: true,
		},
		{
			name: "with hooks (run only on info & debug), no error on debug",
			opts: plogger.NewLoggerOptions().SetHooks(&hookMock{}),
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.DebugLevel,
			wantLogErr: false,
		},
		{
			name: "set levels, only run on error or above",
			opts: plogger.NewLoggerOptions().SetLevel(plogger.ErrorLevel),
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.DebugLevel,
			wantLogErr: false,
		},
		{
			name: "set io.Discard",
			opts: plogger.NewLoggerOptions().SetWriter(io.Discard),
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.DebugLevel,
			wantLogErr: false,
		},
		{
			name: "set mock writer",
			opts: plogger.NewLoggerOptions().SetFormatter(&formatterMock{}),
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.DebugLevel,
			wantLogErr: false,
		},
		{
			name: "set mock writer with error on Format()",
			opts: plogger.NewLoggerOptions().SetFormatter(&formatterWithErrorMock{}),
			entry: &plogger.Entry{
				Msg: "hello world",
			},
			level:      plogger.DebugLevel,
			wantLogErr: true,
		},
		{
			name:             "enable caller with JSON Formatter",
			opts:             plogger.NewLoggerOptions().SetCaller(true),
			entry:            plogger.NewEntry().SetMsg("hello world"),
			level:            plogger.DebugLevel,
			useJSONFormatter: true,
			wantNewErr:       false,
			wantLogErr:       false,
		},
		{
			name:             "enable caller Text Formatter",
			opts:             plogger.NewLoggerOptions().SetCaller(true),
			entry:            plogger.NewEntry().SetMsg("hello world"),
			level:            plogger.DebugLevel,
			useTextFormatter: true,
			wantNewErr:       false,
			wantLogErr:       false,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.opts
			if tt.useTextFormatter {
				f, _ := plogger.NewTextFormatter()
				opts = plogger.NewLoggerOptions().SetCaller(true).SetFormatter(f)
			}
			if tt.useJSONFormatter {
				f, _ := plogger.NewJSONFormatter()
				opts = plogger.NewLoggerOptions().SetCaller(true).SetFormatter(f)
			}

			l, err := plogger.NewLogger(opts)
			if (err != nil) != tt.wantNewErr {
				t.Errorf("NewLogger got error %v, want %v", err, tt.wantNewErr)
			}
			switch tt.level {
			case plogger.TraceLevel:
				err = l.Trace(tt.entry)
			case plogger.DebugLevel:
				err = l.Debug(tt.entry)
			case plogger.InfoLevel:
				err = l.Info(tt.entry)
			case plogger.WarnLevel:
				err = l.Warn(tt.entry)
			case plogger.ErrorLevel:
				err = l.Error(tt.entry)
			case plogger.FatalLevel:
				err = l.Fatal(tt.entry)
			}
			if (err != nil) != tt.wantLogErr {
				t.Errorf("%s() got error %v, want %v", tt.level.String(), err, tt.wantNewErr)
			}
		})
	}
}

type formatterMock struct{}

func (f *formatterMock) Format(*plogger.Entry) ([]byte, error) {
	return []byte{}, nil
}

type formatterWithErrorMock struct{}

func (f *formatterWithErrorMock) Format(*plogger.Entry) ([]byte, error) {
	return []byte{}, errors.New("bug")
}
