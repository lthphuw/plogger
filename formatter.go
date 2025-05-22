package plogger

// Formatter defines how a log entry is serialized (e.g., to JSON or text).
type Formatter interface {
	Format(*Entry) ([]byte, error)
}

// Default key names for the default fields
const (
	FieldKeyMsg        = "msg"
	FieldKeyLevel      = "level"
	FieldKeyTime       = "timestamp"
	FieldKeyFileCaller = "file"
	FieldKeyFuncCaller = "func"
	FieldKeyLineCaller = "line"
)
