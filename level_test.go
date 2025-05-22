package plogger_test

import (
	"testing"

	"github.com/lthphuw/plogger"
)

func TestLevel(t *testing.T) {
	testcases := []struct {
		name           string
		level          plogger.Level
		lvlString      string
		wantErrParse   bool
		wantLevelParse string
		wantParseLevel plogger.Level
	}{
		{
			name:           "Test Level Info",
			level:          plogger.InfoLevel,
			lvlString:      "info",
			wantLevelParse: "info",
			wantParseLevel: plogger.InfoLevel,
		},
		{
			name:           "Test Level Trace",
			level:          plogger.TraceLevel,
			lvlString:      "trace",
			wantLevelParse: "trace",
			wantParseLevel: plogger.TraceLevel,
		},
		{
			name:           "Test Level Debug",
			level:          plogger.DebugLevel,
			lvlString:      "debug",
			wantLevelParse: "debug",
			wantParseLevel: plogger.DebugLevel,
		},
		{
			name:           "Test Level Warn",
			level:          plogger.WarnLevel,
			lvlString:      "warn",
			wantLevelParse: "warn",
			wantParseLevel: plogger.WarnLevel,
		},
		{
			name:           "Test Level Error",
			level:          plogger.ErrorLevel,
			lvlString:      "error",
			wantLevelParse: "error",
			wantParseLevel: plogger.ErrorLevel,
		},
		{
			name:           "Test Level Fatal",
			level:          plogger.FatalLevel,
			lvlString:      "fatal",
			wantLevelParse: "fatal",
			wantParseLevel: plogger.FatalLevel,
		},
		{
			name:           "Test Level Unknown",
			level:          -100,
			lvlString:      "unknown",
			wantErrParse:   true,
			wantLevelParse: "unknown",
			wantParseLevel: plogger.DebugLevel,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			// Parse Level
			parseLevel, err := plogger.ParseLevel(tt.lvlString)
			if (err != nil) != tt.wantErrParse {
				t.Errorf("ParseLevel got error %v, want %v", err, tt.wantErrParse)
			}
			if parseLevel != tt.wantParseLevel {
				t.Errorf("ParseLevel got %v, want %v", parseLevel, tt.level)
			}

			// String
			str := tt.level.String()
			if str != tt.wantLevelParse {
				t.Errorf("String() got %v, want %v", str, tt.lvlString)
			}
		})
	}
}
