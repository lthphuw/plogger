package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/term"
)

var (
	isWindows = strings.Contains(strings.ToLower(runtime.GOOS), "windows")
	isLinux   = strings.Contains(strings.ToLower(runtime.GOOS), "linux")
)

// IsWindows check if OS is windows
func IsWindows() bool {
	return isWindows
}

// IsLinux check if OS is linux
func IsLinux() bool {
	return isLinux
}

// IsTerminal detects terminal
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd())) || term.IsTerminal(int(os.Stderr.Fd()))
}

// ParseInt8 parse a string to int8
func ParseInt8(s string) (int8, error) {
	n, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(n), nil
}

// ParseUInt32 parse a string to unit32
func ParseUInt32(s string) (uint32, error) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(n), nil
}

// RandomHex random a string hex with n bytes
func RandomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
