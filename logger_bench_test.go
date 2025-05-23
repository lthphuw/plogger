package plogger_test

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/lthphuw/plogger"
)

func BenchmarkAddingFields(b *testing.B) {
	b.Logf("")
	b.Logf("Logging with additional context at each log site.")
	b.Run("plogger", func(b *testing.B) {
		logger := newPLogger()
		entry := plogger.NewEntry().SetMsg(getMessage(0)).SetFieldMap(ploggerFieldMap())
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(entry)
			}
		})
	})
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

var (
	errExample = errors.New("fail")

	_messages   = fakeMessages(1000)
	_tenInts    = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	_tenStrings = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	_tenTimes   = []time.Time{
		time.Unix(0, 0),
		time.Unix(1, 0),
		time.Unix(2, 0),
		time.Unix(3, 0),
		time.Unix(4, 0),
		time.Unix(5, 0),
		time.Unix(6, 0),
		time.Unix(7, 0),
		time.Unix(8, 0),
		time.Unix(9, 0),
	}
	_oneUser = &user{
		Name:      "Jane Doe",
		Email:     "jane@test.com",
		CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	_tenUsers = users{
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
	}
)

func fakeMessages(n int) []string {
	messages := make([]string, n)
	for i := range messages {
		messages[i] = fmt.Sprintf(
			"Test logging, but use a somewhat realistic message length. (#%v)",
			i,
		)
	}
	return messages
}

func getMessage(iter int) string {
	return _messages[iter%1000]
}

type users []*user

type user struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
