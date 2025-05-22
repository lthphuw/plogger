package plogger

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Level defines log severity levels.
type Level int8

const (
	// TraceLevel is the most detailed log level.
	TraceLevel Level = iota - 1
	// DebugLevel is for debugging messages.
	DebugLevel
	// InfoLevel is for informational messages.
	InfoLevel
	// WarnLevel is for warnings.
	WarnLevel
	// ErrorLevel is for errors.
	ErrorLevel
	// FatalLevel is for fatal errors that usually exit the app.
	FatalLevel
)

// ParseLevel converts a string to a Level. Returns error if invalid.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "trace":
		return TraceLevel, nil
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "fatal":
		return FatalLevel, nil
	}
	return Level(0), fmt.Errorf("not a valid Level: %q", lvl)
}

// String returns the string representation of the log Level.
func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	}
	return "unknown"
}

// MarshalJSON marshal Level to string
func (l Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}
