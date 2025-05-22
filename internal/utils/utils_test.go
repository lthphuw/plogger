package utils

import (
	"encoding/hex"
	"os"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/term"
)

// TestIsWindows checks if IsWindows returns correct OS detection
func TestIsWindows(t *testing.T) {
	if got := IsWindows(); got != strings.Contains(strings.ToLower(runtime.GOOS), "windows") {
		t.Errorf("IsWindows() = %v, want %v", got, runtime.GOOS == "windows")
	}
}

// TestIsLinux checks if IsLinux returns correct OS detection
func TestIsLinux(t *testing.T) {
	if got := IsLinux(); got != strings.Contains(strings.ToLower(runtime.GOOS), "linux") {
		t.Errorf("IsLinux() = %v, want %v", got, runtime.GOOS == "linux")
	}
}

// TestIsTerminal checks if IsTerminal detects terminal correctly
func TestIsTerminal(t *testing.T) {
	// Note: This test may depend on environment (CI vs local)
	got := IsTerminal()
	isTerm := term.IsTerminal(int(os.Stdout.Fd())) || term.IsTerminal(int(os.Stderr.Fd()))
	if got != isTerm {
		t.Errorf("IsTerminal() = %v, want %v", got, isTerm)
	}
}

// TestParseInt8 checks ParseInt8 for valid and invalid inputs
func TestParseInt8(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int8
		wantErr bool
	}{
		{"Valid int8", "127", 127, false},
		{"Invalid int8", "128", 0, true},
		{"Non-numeric", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInt8(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt8() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseInt8() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseUInt32 checks ParseUInt32 for valid and invalid inputs
func TestParseUInt32(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint32
		wantErr bool
	}{
		{"Valid uint32", "1000", 1000, false},
		{"Negative number", "-1", 0, true},
		{"Non-numeric", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUInt32(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUInt32() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseUInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRandomHex checks if RandomHex generates correct length hex string
func TestRandomHex(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"Length 0", 0},
		{"Length 4", 4},
		{"Length 10", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomHex(tt.length)
			// Hex string length is twice the input byte length
			if len(got) != tt.length*2 {
				t.Errorf("RandomHex(%d) length = %d, want %d", tt.length, len(got), tt.length*2)
			}
			// Check if output is valid hex
			if _, err := hex.DecodeString(got); err != nil {
				t.Errorf("RandomHex(%d) produced invalid hex: %v", tt.length, err)
			}
		})
	}
}
