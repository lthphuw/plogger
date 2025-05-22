// Package trace provides utilities to retrieve and format runtime stack trace
// and caller information, excluding internal packages.
package trace

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

const (
	maximumCallerDepth = 25
)

var (
	tracePackage   string
	callerInitOnce sync.Once
)

// CallerInfo represent for a Frame
type CallerInfo struct {
	File     string
	Line     int
	Function string
}

// getPackageName reduces a fully qualified function name to the package name.
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

// isInternalPackage checks if a package is considered internal.
func isInternalPackage(pkg string) bool {
	internalPackages := []string{
		"runtime",
		"github.com/lthphuw/plogger",
		tracePackage,
	}
	for _, internal := range internalPackages {
		if strings.HasPrefix(pkg, internal) {
			return true
		}
	}
	return false
}

// GetCaller returns the first frame not in the current/internal packages.
func GetCaller() *runtime.Frame {
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		n := runtime.Callers(0, pcs)
		for i := 0; i < n; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "GetCaller") {
				tracePackage = getPackageName(funcName)
				break
			}
		}
	})

	pcs := make([]uintptr, maximumCallerDepth)
	n := runtime.Callers(0, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	// Skip until we see our package
	for {
		f, more := frames.Next()
		if !more {
			break
		}
		pkg := getPackageName(f.Function)
		if !isInternalPackage(pkg) {
			return &f
		}
	}
	return nil
}

// GetStackTrace returns the stack trace starting from `skip` and limited to `depth` frames.
func GetStackTrace(skip int, depth int) []runtime.Frame {
	pcs := make([]uintptr, maximumCallerDepth)
	n := runtime.Callers(skip, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	var stack []runtime.Frame
	for range depth {
		f, more := frames.Next()
		if !more {
			break
		}
		stack = append(stack, f)
	}
	return stack
}

// FormatStackTrace formats a slice of runtime.Frame into a readable stack trace string.
func FormatStackTrace(frames []runtime.Frame) string {
	var b strings.Builder
	for _, f := range frames {
		b.WriteString(f.Function)
		b.WriteString("\n\t")
		b.WriteString(f.File)
		b.WriteString(":")
		b.WriteString(fmt.Sprintf("%d", f.Line))
		b.WriteString("\n")
	}
	return b.String()
}
