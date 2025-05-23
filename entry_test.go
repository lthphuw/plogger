package plogger_test

import (
	"testing"
	"time"

	"github.com/lthphuw/plogger"
)

func TestSetTimestamp(t *testing.T) {
	now := time.Now().Add(-time.Hour)

	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{"Set timestamp to past", now, now},
		{"Set timestamp to now", time.Now(), time.Now()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := plogger.NewEntry().SetTimestamp(tt.input)
			if !e.Timestamp.Equal(tt.expected) && e.Timestamp.Sub(tt.expected) > time.Second {
				t.Errorf("expected %v, got %v", tt.expected, e.Timestamp)
			}
		})
	}
}

func TestSetMsg(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Set normal msg", "hello", "hello"},
		{"Set empty msg", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := plogger.NewEntry().SetMsg(tt.input)
			if e.Msg != tt.expected {
				t.Errorf("expected msg %q, got %q", tt.expected, e.Msg)
			}
		})
	}
}

func TestSetFieldMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []map[string]any
		expected int
	}{
		{"Set non-empty map", []map[string]any{{"a": 1, "b": "x"}}, 2},
		{"Set empty map", []map[string]any{{}}, 0},
		{"Set nil map", []map[string]any{nil}, 0},
		{
			"Set multi map",
			[]map[string]any{nil, {"a": 1, "b": "x"}, {"aab": 32, "ab": "12x", "123": "#2"}},
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := plogger.NewEntry()
			for _, inp := range tt.input {
				e.SetFieldMap(inp)
			}
			if len(e.FieldMap) != tt.expected {
				t.Errorf("expected field map size %d, got %d", tt.expected, len(e.FieldMap))
			}

			for k, v := range tt.input[len(tt.input)-1] {
				if e.FieldMap[k] != v {
					t.Errorf("expected key %s to have value %v", k, v)
				}
			}
		})
	}
}

func TestAddField(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected any
	}{
		{"Add int field", "count", 42, 42},
		{"Add string field", "name", "John", "John"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := plogger.NewEntry().AddField(tt.key, tt.value)
			if got := e.FieldMap[tt.key]; got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestAddFields(t *testing.T) {
	tests := []struct {
		name     string
		fields   map[string]any
		expected int
	}{
		{"Add two fields", map[string]any{"x": 1, "y": 2}, 2},
		{"Add no field", map[string]any{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := plogger.NewEntry().AddFields(tt.fields)
			if len(e.FieldMap) != tt.expected {
				t.Errorf("expected %d fields, got %d", tt.expected, len(e.FieldMap))
			}
		})
	}
}

func TestAcquireAndReleaseEntry(t *testing.T) {
	t.Run("Reset fields after release", func(t *testing.T) {
		e := plogger.AcquireEntry()
		e.SetMsg("hello").AddField("x", 1).SetTimestamp(time.Unix(0, 0))
		plogger.ReleaseEntry(e)

		e2 := plogger.AcquireEntry()
		if e2.Msg != "" || len(e2.FieldMap) != 0 {
			t.Errorf("Entry not reset properly: got Msg=%q, FieldMap=%v", e2.Msg, e2.FieldMap)
		}
		plogger.ReleaseEntry(e2)
	})
}

func TestEntryString(t *testing.T) {
	tests := []struct {
		name     string
		level    plogger.Level
		msg      string
		expected string
	}{
		{"Basic info log", plogger.InfoLevel, "message", `level=info msg="message"`},
		{"Empty log", 0, "", `level=debug msg=""`}, // default level
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := plogger.NewEntry()
			e.Level = tt.level
			e.Msg = tt.msg
			out := e.String()
			if out != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, out)
			}
		})
	}
}

// Using Pooling is more recommended
// goos: darwin
// goarch: arm64
// pkg: plogger
// cpu: Apple M3 Pro
// BenchmarkNewEntry/NewEntry_(No_Pooling)-11         	 2713657	       446.9 ns/op	    1368 B/op	      12 allocs/op
// BenchmarkNewEntry/NewEntry_(Pooling)-11            	 3893079	       311.3 ns/op	     752 B/op	       9 allocs/op
func BenchmarkNewEntry(b *testing.B) {
	b.Run("NewEntry (No Pooling)", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				e := plogger.NewEntry()
				e.SetMsg("Hello world").SetFieldMap(map[string]any{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
					"key4": "value4",
					"key5": "value5",
				})
			}
		})
	})

	b.Run("NewEntry (Pooling)", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				e := plogger.AcquireEntry()
				e.SetMsg("Hello world").SetFieldMap(map[string]any{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
					"key4": "value4",
					"key5": "value5",
				})
				plogger.ReleaseEntry(e)
			}
		})
	})
}
