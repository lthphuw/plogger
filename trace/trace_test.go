package trace

import (
	"strings"
	"testing"
)

// Test table-driven cho getPackageName
func TestGetPackageName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"runtime.main", "runtime"},
		{"github.com/lthphuw/plogger.Logger.Log", "github.com/lthphuw/plogger"},
		{"github.com/lthphuw/plogger.subpkg.Func", "github.com/lthphuw/plogger"},
		{"main", "main"},
		{"github.com/user/project/plogger.FuncName", "github.com/user/project/plogger"},
		{"github.com/user/project/plogger.subpkg.Func", "github.com/user/project/plogger"},
	}

	for _, tt := range tests {
		got := getPackageName(tt.input)
		if got != tt.want {
			t.Errorf("getPackageName(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

// Test isInternalPackage với package "github.com/lthphuw/plogger" làm tracePackage
func TestIsInternalPackage(t *testing.T) {
	// Giả sử tracePackage = "github.com/lthphuw/plogger"
	tracePackage = "github.com/lthphuw/plogger"

	cases := []struct {
		pkg  string
		want bool
	}{
		{"runtime", true},
		{"github.com/lthphuw/plogger", true},
		{"github.com/lthphuw/plogger.subpkg", true},
		{"someotherpkg", false},
		{"main", false},
	}

	for _, c := range cases {
		got := isInternalPackage(c.pkg)
		if got != c.want {
			t.Errorf("isInternalPackage(%q) = %v; want %v", c.pkg, got, c.want)
		}
	}
}

// Test GetCaller trả về frame bên ngoài package internal
func TestGetCaller(t *testing.T) {
	frame := GetCaller()
	if frame == nil {
		t.Fatal("GetCaller() returned nil")
	}

	pkg := getPackageName(frame.Function)
	if isInternalPackage(pkg) {
		t.Errorf("GetCaller() returned internal package frame %q", frame.Function)
	}
}

// Test GetStackTrace & FormatStackTrace
func TestGetStackTraceAndFormat(t *testing.T) {
	depth := 5
	frames := GetStackTrace(0, depth)

	if len(frames) == 0 {
		t.Fatal("GetStackTrace returned empty slice")
	}

	if len(frames) > depth {
		t.Errorf("GetStackTrace returned more frames (%d) than depth (%d)", len(frames), depth)
	}

	formatted := FormatStackTrace(frames)
	for _, f := range frames {
		if !strings.Contains(formatted, f.Function) {
			t.Errorf("FormatStackTrace output missing function %q", f.Function)
		}
		if !strings.Contains(formatted, f.File) {
			t.Errorf("FormatStackTrace output missing file %q", f.File)
		}
	}
}
