package plogger_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lthphuw/plogger"
)

type hookMock struct{}

func (h *hookMock) Levels() []plogger.Level {
	return []plogger.Level{plogger.DebugLevel, plogger.InfoLevel}
}

func (h *hookMock) Run(entry *plogger.Entry) error {
	// Assume bug on info
	if entry.Level != plogger.DebugLevel {
		return errors.New("bug")
	}

	return nil
}

func TestHook(t *testing.T) {
	testcases := []struct {
		name string
		hook *hookMock
	}{
		{
			name: "Default",
			hook: &hookMock{},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			levelHook := make(plogger.LevelHook)
			entry := &plogger.Entry{
				Msg: "hello world",
			}
			levelHook.Add(tt.hook)
			fmt.Println("level hook: ", levelHook)
			entry.Level = plogger.DebugLevel
			err := levelHook.Run(plogger.DebugLevel, entry)
			if err != nil {
				t.Errorf("expected nil err, got %v", err)
			}

			entry.Level = plogger.InfoLevel
			err = levelHook.Run(plogger.InfoLevel, entry)
			if err == nil {
				t.Errorf("expected err, got nil")
			}

			// It actually dont run the hook
			entry.Level = plogger.TraceLevel
			err = levelHook.Run(plogger.TraceLevel, entry)
			if err != nil {
				t.Errorf("expected nil err, got %v", err)
			}

			entry.Level = plogger.WarnLevel
			err = levelHook.Run(plogger.WarnLevel, entry)
			if err != nil {
				t.Errorf("expected nil err, got %v", err)
			}

			entry.Level = plogger.ErrorLevel
			err = levelHook.Run(plogger.ErrorLevel, entry)
			if err != nil {
				t.Errorf("expected nil err, got %v", err)
			}

			entry.Level = plogger.FatalLevel
			err = levelHook.Run(plogger.FatalLevel, entry)
			if err != nil {
				t.Errorf("expected nil err, got %v", err)
			}
		})
	}
}
