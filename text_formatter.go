package plogger

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/lthphuw/plogger/internal/color"
	"github.com/lthphuw/plogger/internal/utils"
	"github.com/lthphuw/plogger/trace"
)

// TextFormatterOptions configures the behavior of a text formatter.
type TextFormatterOptions struct {
	// DisableColor determines whether to disable colored level output
	DisableColor *bool `default:"false"` // optional; default: false

	// DisableSorting determines whether to disable sorting of fields.
	DisableSorting *bool `default:"false"` // optional; default: false

	// SortingFieldKeyFunc specifies a custom function to sort field keys.
	SortingFieldKeyFunc func([]string) // optional; default: sort.Strings

	// TimestampFormat specifies the format for timestamps, following Go's time format.
	TimestampFormat *string `default:"2006-01-02T15:04:05Z07:00"` // optional; default: "2006-01-02T15:04:05Z07:00"
}

// TextFormatterOptionsBuilder builds and applies a list of setters for TextFormatterOptions.
type TextFormatterOptionsBuilder struct {
	Options []Setter[TextFormatterOptions] // Options is the list of setter functions to configure TextFormatterOptions.
}

// TextFormatter formats text output based on options.
type TextFormatter struct {
	*TextFormatterOptions // TextFormatterOptions holds configuration for formatting.

	pool *sync.Pool // pool is used to reuse buffers and reduce allocations.
}

// NewTextFormatter creates a new TextFormatter
func NewTextFormatter(opts ...Lister[TextFormatterOptions]) (*TextFormatter, error) {
	args, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}
	if args.SortingFieldKeyFunc == nil {
		args.SortingFieldKeyFunc = sort.Strings
	}
	return &TextFormatter{
		TextFormatterOptions: args,
		pool: &sync.Pool{
			New: func() any { return &bytes.Buffer{} },
		},
	}, nil
}

// Format format entry to key=value, and return []byte, error
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	keys := make([]string, 3+len(entry.FieldMap))
	keys[0] = FieldKeyTime
	keys[1] = FieldKeyLevel
	keys[2] = FieldKeyMsg

	i := 3
	for k := range entry.FieldMap {
		keys[i] = k
		i++
	}

	if !*f.DisableSorting {
		f.SortingFieldKeyFunc(keys[3:])
	}

	// Caller
	var funcVal, fileVal string
	var lineVal int
	if entry.caller {
		f := trace.GetCaller()
		funcVal = f.Function
		fileVal = f.File
		lineVal = f.Line
		if funcVal != "" {
			keys = append(keys, FieldKeyFuncCaller)
		}
		if fileVal != "" {
			keys = append(keys, FieldKeyFileCaller)
		}
		if lineVal != 0 {
			keys = append(keys, FieldKeyLineCaller)
		}
	}

	// Create a buffer for append string
	buf, ok := f.pool.Get().(*bytes.Buffer)
	if !ok {
		buf = &bytes.Buffer{}
	}

	defer f.pool.Put(buf)
	buf.Reset()

	for _, k := range keys {
		var v any
		switch k {
		case FieldKeyTime:
			if entry.Timestamp.IsZero() {
				entry.Timestamp = time.Now()
			}
			v = entry.Timestamp.Format(*f.TimestampFormat)

		case FieldKeyLevel:
			if f.shouldColornize() {
				v = f.colornizeLevel(entry.Level)
			} else {
				v = entry.Level.String()
			}

		case FieldKeyMsg:
			v = entry.Msg

		case FieldKeyFuncCaller:
			v = funcVal

		case FieldKeyFileCaller:
			v = fileVal

		case FieldKeyLineCaller:
			v = lineVal

		default:
			v = entry.FieldMap[k]
		}
		f.appendKeyValue(buf, k, v)
	}
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

// shouldColornize determine whether we should color the level
func (f *TextFormatter) shouldColornize() bool {
	return !*f.DisableColor && utils.IsTerminal() && !utils.IsWindows()
}

// colornizeLevel colornizes the level.
func (f *TextFormatter) colornizeLevel(level Level) string {
	lvl := level.String()
	switch level {
	case TraceLevel:
		return color.TraceColor.Add(lvl)
	case DebugLevel:
		return color.DebugColor.Add(lvl)
	case InfoLevel:
		return color.InfoColor.Add(lvl)
	case WarnLevel:
		return color.WarnColor.Add(lvl)
	case ErrorLevel:
		return color.ErrorColor.Add(lvl)
	case FatalLevel:
		return color.FatalColor.Add(lvl)
	}
	return lvl
}

// appendKeyValue appends a key-value pair to the buffer in the format "key=value".
func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value any) {
	if value == nil {
		value = "nil"
	}
	if b.Len() > 0 {
		fmt.Fprintf(b, " %s=%v", key, value)
	} else {
		fmt.Fprintf(b, "%s=%v", key, value)
	}
}

// NewTextFormatterOptions creates a new *TextFormatterOptionsBuilder
func NewTextFormatterOptions() *TextFormatterOptionsBuilder {
	return &TextFormatterOptionsBuilder{}
}

// List lists all the setter function for TextFormatter
func (b *TextFormatterOptionsBuilder) List() []Setter[TextFormatterOptions] {
	return b.Options
}

// SetDisableColor sets disable color
func (b *TextFormatterOptionsBuilder) SetDisableColor(disable bool) *TextFormatterOptionsBuilder {
	b.Options = append(b.Options, func(tfo *TextFormatterOptions) error {
		tfo.DisableColor = &disable

		return nil
	})
	return b
}

// SetTimestampFormat sets timestamp format
func (b *TextFormatterOptionsBuilder) SetTimestampFormat(
	format string,
) *TextFormatterOptionsBuilder {
	b.Options = append(b.Options, func(tfo *TextFormatterOptions) error {
		tfo.TimestampFormat = &format

		return nil
	})
	return b
}

// SetDisableSorting sets disable sorting.
func (b *TextFormatterOptionsBuilder) SetDisableSorting(disable bool) *TextFormatterOptionsBuilder {
	b.Options = append(b.Options, func(tfo *TextFormatterOptions) error {
		tfo.DisableSorting = &disable

		return nil
	})
	return b
}

// SetSortingFieldKeyFunc sets the sorting key function.
func (b *TextFormatterOptionsBuilder) SetSortingFieldKeyFunc(
	sortingFieldKeyFunc func([]string),
) *TextFormatterOptionsBuilder {
	b.Options = append(b.Options, func(tfo *TextFormatterOptions) error {
		tfo.SortingFieldKeyFunc = sortingFieldKeyFunc

		return nil
	})
	return b
}
