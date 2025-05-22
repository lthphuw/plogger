// Package color defines colors for log levels.
package color

import "fmt"

// Color represents to ANSI color
type Color string

const (
	// TraceColor use White color
	TraceColor Color = "\033[37m"
	// DebugColor use  Cyan color
	DebugColor Color = "\033[36m"
	// InfoColor use  Green color
	InfoColor Color = "\033[32m"
	//  WarnColor use Yellow color
	WarnColor Color = "\033[33m"
	//  ErrorColor use Red (foreground) color
	ErrorColor Color = "\033[31m"
	//  FatalColor use Red (background) color
	FatalColor Color = "\033[41m"
	//  ResetColor use Reset color
	ResetColor Color = "\033[0m"
)

// Add adds color to a string
func (c Color) Add(s string) string {
	return fmt.Sprintf("%s%s%s", c, s, ResetColor)
}
