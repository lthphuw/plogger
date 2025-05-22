package plogger

import (
	"bytes"
	"fmt"
	"maps"
	"sync"
	"time"

	"github.com/lthphuw/plogger/trace"
	"github.com/segmentio/encoding/json"
)

// JSONFormatterOptions configures JSON formatter behavior.
type JSONFormatterOptions struct {
	TimestampFormat *string `default:"2006-01-02T15:04:05Z07:00"` // time format for timestamps
	PrettyPrint     *bool   `default:"true"`                      // enable pretty print
	EscapeHTML      *bool   `default:"true"`                      // escape HTML characters
}

// JSONFormatterOptionsBuilder builds JSONFormatterOptions using setters.
type JSONFormatterOptionsBuilder struct {
	Options []Setter[JSONFormatterOptions]
}

// JSONFormatter formats log entries as JSON.
type JSONFormatter struct {
	*JSONFormatterOptions

	pool     *sync.Pool
	dataPool *sync.Pool
}

// NewJSONFormatter creates a new JSONFormatter with options.
func NewJSONFormatter(opts ...Lister[JSONFormatterOptions]) (*JSONFormatter, error) {
	args, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}
	return &JSONFormatter{
		JSONFormatterOptions: args,
		pool: &sync.Pool{
			New: func() any { return &bytes.Buffer{} },
		},
		dataPool: &sync.Pool{
			New: func() any { return make(map[string]any) },
		},
	}, nil
}

// Format converts a log Entry to JSON bytes.
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	data := f.dataFields()
	defer f.dataPool.Put(data)

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	data[FieldKeyTime] = entry.Timestamp.Format(*f.TimestampFormat)
	data[FieldKeyLevel] = entry.Level
	data[FieldKeyMsg] = entry.Msg

	if entry.caller {
		f := trace.GetCaller()
		if f.Function != "" {
			data[FieldKeyFuncCaller] = f.Function
		}
		if f.File != "" {
			data[FieldKeyFileCaller] = f.File
		}
		if f.Line != 0 {
			data[FieldKeyLineCaller] = f.Line
		}
	}
	maps.Copy(data, entry.FieldMap)

	buf, ok := f.pool.Get().(*bytes.Buffer)
	if !ok {
		buf = &bytes.Buffer{}
	}
	defer f.pool.Put(buf)
	buf.Reset()

	enc := json.NewEncoder(buf)
	if f.EscapeHTML != nil {
		enc.SetEscapeHTML(*f.EscapeHTML)
	}
	if f.PrettyPrint != nil && *f.PrettyPrint {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}
	return buf.Bytes(), nil
}

func (f *JSONFormatter) dataFields() map[string]any {
	data, ok := f.dataPool.Get().(map[string]any)
	if !ok {
		return make(map[string]any)
	}
	for k := range data {
		delete(data, k)
	}
	return data
}

// NewJSONFormatterOptions creates a builder for JSONFormatterOptions.
func NewJSONFormatterOptions() *JSONFormatterOptionsBuilder {
	return &JSONFormatterOptionsBuilder{}
}

// List returns the list of option setters.
func (b *JSONFormatterOptionsBuilder) List() []Setter[JSONFormatterOptions] {
	return b.Options
}

// SetTimestampFormat sets the timestamp format.
func (b *JSONFormatterOptionsBuilder) SetTimestampFormat(
	format string,
) *JSONFormatterOptionsBuilder {
	b.Options = append(b.Options, func(jfo *JSONFormatterOptions) error {
		jfo.TimestampFormat = &format
		return nil
	})
	return b
}

// SetPrettyPrint enables or disables pretty print.
func (b *JSONFormatterOptionsBuilder) SetPrettyPrint(pretty bool) *JSONFormatterOptionsBuilder {
	b.Options = append(b.Options, func(jfo *JSONFormatterOptions) error {
		jfo.PrettyPrint = &pretty
		return nil
	})
	return b
}

// SetEscapeHTML enables or disables HTML escaping.
func (b *JSONFormatterOptionsBuilder) SetEscapeHTML(escape bool) *JSONFormatterOptionsBuilder {
	b.Options = append(b.Options, func(jfo *JSONFormatterOptions) error {
		jfo.EscapeHTML = &escape
		return nil
	})
	return b
}
