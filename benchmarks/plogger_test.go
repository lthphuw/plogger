package benchmarks

import (
	"io"

	"github.com/lthphuw/plogger"
)

func newDisabledPLogger() *plogger.Logger {
	writer, _ := plogger.NewDiscardWriter()
	formatter, _ := plogger.NewJSONFormatter()
	opts := plogger.NewLoggerOptions().
		SetFormatter(formatter).
		SetWriter(writer).
		SetLevel(plogger.ErrorLevel)
	logger, _ := plogger.NewLogger(opts)
	return logger
}

func newPLogger() *plogger.Logger {
	// writer, _ := plogger.NewDiscardWriter()
	formatter, _ := plogger.NewJSONFormatter(
		plogger.NewJSONFormatterOptions().SetEscapeHTML(false).SetPrettyPrint(false),
	)

	opts := plogger.NewLoggerOptions().
		SetFormatter(formatter).
		SetWriter(io.Discard).
		SetCaller(false).
		SetLevel(plogger.DebugLevel)
	logger, _ := plogger.NewLogger(opts)
	return logger
}

func ploggerFieldMap() map[string]any {
	return map[string]any{
		"int":     _tenInts[0],
		"ints":    _tenInts,
		"string":  _tenStrings[0],
		"strings": _tenStrings,
		"time":    _tenTimes[0],
		"times":   _tenTimes,
		"user1":   _oneUser,
		"user2":   _oneUser,
		"users":   _tenUsers,
		"error":   errExample,
	}
}
