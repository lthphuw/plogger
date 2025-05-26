package plogger

// Formatter defines how a log entry is serialized (e.g., to JSON or text).
type Formatter interface {
	Format(*Entry) ([]byte, error)
}

// Default key names for the default fields
const (
	DefaultFieldKeyMsg        = "msg"
	DefaultFieldKeyLevel      = "level"
	DefaultFieldKeyTime       = "timestamp"
	DefaultFieldKeyFileCaller = "file"
	DefaultFieldKeyFuncCaller = "func"
	DefaultFieldKeyLineCaller = "line"
)

// Field key names used for serialization (can be customized)
var (
	FieldKeyMsg        = DefaultFieldKeyMsg
	FieldKeyLevel      = DefaultFieldKeyLevel
	FieldKeyTime       = DefaultFieldKeyTime
	FieldKeyFileCaller = DefaultFieldKeyFileCaller
	FieldKeyFuncCaller = DefaultFieldKeyFuncCaller
	FieldKeyLineCaller = DefaultFieldKeyLineCaller
)

// SetFieldKeyMsg sets the field key used for the log message (default: "msg").
func SetFieldKeyMsg(key string) {
	FieldKeyMsg = key
}

// SetFieldKeyLevel sets the field key used for the log level (default: "level").
func SetFieldKeyLevel(key string) {
	FieldKeyLevel = key
}

// SetFieldKeyTimestamp sets the field key used for the log timestamp (default: "timestamp").
func SetFieldKeyTimestamp(key string) {
	FieldKeyTime = key
}

// SetFieldKeyFileCaller sets the field key used for the caller's file path (default: "file").
func SetFieldKeyFileCaller(key string) {
	FieldKeyFileCaller = key
}

// SetFieldKeyFuncCaller sets the field key used for the caller's function name (default: "func").
func SetFieldKeyFuncCaller(key string) {
	FieldKeyFuncCaller = key
}

// SetFieldKeyLineCaller sets the field key used for the caller's line number (default: "line").
func SetFieldKeyLineCaller(key string) {
	FieldKeyLineCaller = key
}

// SetFieldKeyAsDefault resets all field key names to their default values.
func SetFieldKeyAsDefault() {
	FieldKeyMsg = DefaultFieldKeyMsg
	FieldKeyLevel = DefaultFieldKeyLevel
	FieldKeyTime = DefaultFieldKeyTime
	FieldKeyFileCaller = DefaultFieldKeyFileCaller
	FieldKeyFuncCaller = DefaultFieldKeyFuncCaller
	FieldKeyLineCaller = DefaultFieldKeyLineCaller
}
