package color

import (
	"strings"
	"testing"
)

func TestColorAdd(t *testing.T) {
	tests := []struct {
		name      string
		color     Color
		input     string
		wantStart string
		wantEnd   string
	}{
		{
			name:      "Info color with color enabled",
			color:     InfoColor,
			input:     "Hello",
			wantStart: string(InfoColor),
			wantEnd:   string(ResetColor),
		},
		{
			name:      "Trace color with color enabled",
			color:     TraceColor,
			input:     "Failed",
			wantStart: string(TraceColor),
			wantEnd:   string(ResetColor),
		},
		{
			name:      "Debug color with color disable",
			color:     DebugColor,
			input:     "Hello",
			wantStart: string(DebugColor),
			wantEnd:   string(ResetColor),
		},
		{
			name:      "Warn color with color disable",
			color:     WarnColor,
			input:     "Failed",
			wantStart: string(WarnColor),
			wantEnd:   string(ResetColor),
		},
		{
			name:      "Error color with color disable",
			color:     ErrorColor,
			input:     "Failed",
			wantStart: string(ErrorColor),
			wantEnd:   string(ResetColor),
		},
		{
			name:      "Fatal color with color disable",
			color:     FatalColor,
			input:     "Failed",
			wantStart: string(FatalColor),
			wantEnd:   string(ResetColor),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.color.Add(tt.input)

			expectedLen := len(tt.wantStart) + len(tt.input) + len(tt.wantEnd)
			if len(actual) != expectedLen {
				t.Errorf(
					"result length mismatch: got %d, want %d, actual: %q",
					len(actual),
					expectedLen,
					actual,
				)
			}
			if !strings.HasPrefix(actual, tt.wantStart) {
				t.Errorf(
					"start color mismatch: got %q, want %q",
					actual[:len(tt.wantStart)],
					tt.wantStart,
				)
			}
			if !strings.HasSuffix(actual, tt.wantEnd) {
				t.Errorf(
					"reset color mismatch: got %q, want %q",
					actual[len(actual)-len(tt.wantEnd):],
					tt.wantEnd,
				)
			}
		})
	}
}
