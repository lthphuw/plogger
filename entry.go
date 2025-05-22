package plogger

import (
	"fmt"
	"sync"
	"time"
)

// entryPool is a sync.Pool to reuse Entry objects for performance.
var entryPool = sync.Pool{
	New: func() any {
		return &Entry{
			FieldMap: make(map[string]any),
		}
	},
}

// Entry represents a log message with metadata like level, timestamp, and fields.
type Entry struct {
	Timestamp time.Time      `json:"timestamp"` // Log timestamp
	Level     Level          `json:"level"`     // Log level
	Msg       string         `json:"msg"`       // Log message
	FieldMap  map[string]any // Custom log fields

	caller bool // Whether to include caller info
}

// NewEntry creates a new log entry with initialized field map and current timestamp.
func NewEntry() *Entry {
	return &Entry{
		Timestamp: time.Now(),
		FieldMap:  make(map[string]any),
	}
}

// AcquireEntry returns a new or reused Entry from the pool.
func AcquireEntry() *Entry {
	e, ok := entryPool.Get().(*Entry)
	if !ok {
		e = NewEntry()
	}
	return reset(e)
}

// reset clears and reinitializes fields of an Entry.
func reset(e *Entry) *Entry {
	e.Timestamp = time.Now()
	e.Level = 0
	e.Msg = ""
	e.caller = false
	for k := range e.FieldMap {
		delete(e.FieldMap, k)
	}
	return e
}

// ReleaseEntry puts an Entry back to the pool.
func ReleaseEntry(e *Entry) {
	entryPool.Put(e)
}

// SetTimestamp sets the timestamp for the Entry.
func (e *Entry) SetTimestamp(t time.Time) *Entry {
	e.Timestamp = t
	return e
}

// SetMsg sets the log message for the Entry.
func (e *Entry) SetMsg(m string) *Entry {
	e.Msg = m
	return e
}

// SetFieldMap sets a full map of fields, replacing any existing fields.
func (e *Entry) SetFieldMap(fm map[string]any) *Entry {
	if e.FieldMap == nil {
		e.FieldMap = make(map[string]any, len(fm))
	} else {
		for k := range e.FieldMap {
			delete(e.FieldMap, k)
		}
	}
	for k, v := range fm {
		e.FieldMap[k] = v
	}
	return e
}

// AddField adds a single key-value field to the Entry.
func (e *Entry) AddField(key string, val any) *Entry {
	if e.FieldMap == nil {
		e.FieldMap = make(map[string]any)
	}
	e.FieldMap[key] = val
	return e
}

// AddFields adds multiple fields to the Entry.
func (e *Entry) AddFields(fields map[string]any) *Entry {
	if e.FieldMap == nil {
		e.FieldMap = make(map[string]any, len(fields))
	}
	for k, v := range fields {
		e.FieldMap[k] = v
	}
	return e
}

// String returns a string representation of the Entry.
func (e *Entry) String() string {
	return fmt.Sprintf("level=%s msg=%q", e.Level.String(), e.Msg)
}
