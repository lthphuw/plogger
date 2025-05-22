package plogger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lthphuw/plogger/internal/utils"
)

const (
	megabyte = 1024 * 1024
	fileMode = 0o644
	dirMode  = 0o755
)

// FileWriterOptions configures the behavior of FileWriter.
type FileWriterOptions struct {
	// Filename is the path to the log file.
	Filename *string `default:"/tmp/test.log"`

	// MaxSize defines the maximum size (in MB) before rotating the log file.
	MaxSize *uint32 `default:"1024"`

	// BackupFormat specifies the timestamp format used in backup filenames.
	// Format follows Go's time format (e.g., "20060102_150405").
	BackupFormat *string `default:"20060102_150405"`

	// BackupRandomHexSuffix adds a random hex string to backup filenames to avoid conflicts.
	BackupRandomHexSuffix *bool `default:"false"`
}

// FileWriterOptionsBuilder builds FileWriterOptions using setters.
type FileWriterOptionsBuilder struct {
	Options []Setter[FileWriterOptions]
}

// FileWriter writes logs to file and rotate it.
type FileWriter struct {
	// FileWriterOptions holds configuration for FileWriter.
	*FileWriterOptions

	// file holds current file to append
	file *os.File
	// size holds current size of the fiel
	size int64
	// mu is a mutex
	mu sync.Mutex
}

// NewFileWriter create a FileWriter
func NewFileWriter(opts ...Lister[FileWriterOptions]) (*FileWriter, error) {
	args, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}
	return &FileWriter{
		FileWriterOptions: args,
	}, nil
}

// Write appends bytes to current file (rotate if needed)
func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	writeLen := int64(len(p))

	// Write byte's length exceed maximum file size
	if writeLen > w.maxSize() {
		return 0, fmt.Errorf(
			"write length %d exceeds maximum file size %d", writeLen, w.maxSize(),
		)
	}
	// Open file if not yet
	if err := w.openExistOrNewFile(); err != nil {
		return 0, err
	}

	// When byte's length exceed current file size,
	// rename current file as backup & create new file.
	if w.size+writeLen > w.maxSize() {
		if err := w.openNewFile(); err != nil {
			return 0, err
		}
	}

	// Write bytes
	n, err = w.file.Write(p)
	w.size += int64(n)
	return n, err
}

// openExistOrNewFile open a existing file or create a new file
func (w *FileWriter) openExistOrNewFile() error {
	// If already open the file, do nothing
	if w.isValidFile() {
		return nil
	}

	file, err := os.OpenFile(*w.Filename, os.O_APPEND|os.O_WRONLY, fileMode)
	if os.IsNotExist(err) {
		return w.openNewFile()
	}
	if err != nil {
		return fmt.Errorf("error opening log file: %s", err)
	}

	info, err := file.Stat()
	if err != nil {
		if err2 := file.Close(); err2 != nil {
			return fmt.Errorf("error getting log file info: %s, also close got err: %s", err, err2)
		}
		return fmt.Errorf("error getting log file info: %s", err)
	}

	w.file = file
	w.size = info.Size()
	return nil
}

// openNewFile rename current file name as a backup file (if any), and create a new file.
func (w *FileWriter) openNewFile() error {
	mode := os.FileMode(fileMode)
	// If current file is valid, rename it.
	if w.isValidFile() {
		if err := w.close(); err != nil {
			return err
		}
		info, err := os.Stat(*w.Filename)
		if err == nil {
			mode = info.Mode()
			// move the existing file
			newname := w.backupLogFilename()
			if err := os.Rename(*w.Filename, newname); err != nil {
				return fmt.Errorf("can't rename log file: %s", err)
			}
		}
	}

	dir := filepath.Dir(*w.Filename)

	// Create directories if not exists.
	if err := os.MkdirAll(dir, dirMode); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	file, err := os.OpenFile(*w.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, mode)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	info, err := os.Stat(*w.Filename)
	if err != nil {
		return fmt.Errorf("failed to read file stat: %w", err)
	}
	w.file = file
	w.size = info.Size()

	return nil
}

// backupLogFilename create a new log file name.
func (w *FileWriter) backupLogFilename() string {
	ext := filepath.Ext(*w.Filename)
	base := strings.TrimSuffix(filepath.Base(*w.Filename), ext)
	timestamp := time.Now().Format(*w.BackupFormat)
	if w.BackupRandomHexSuffix == nil || !(*w.BackupRandomHexSuffix) {
		return filepath.Join(
			filepath.Dir(*w.Filename),
			fmt.Sprintf("%s_%s%s", base, timestamp, ext),
		)
	}
	return filepath.Join(
		filepath.Dir(*w.Filename),
		fmt.Sprintf("%s_%s_%s%s", base, timestamp, utils.RandomHex(6), ext),
	)
}

// Close closes the file
func (w *FileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.close()
}

// close closes the file if it is open.
func (w *FileWriter) close() error {
	if w.file == nil || w.file.Fd() == ^uintptr(0) {
		w.file = nil
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}

// isValidFile checks if current file is valid
func (w *FileWriter) isValidFile() bool {
	return w.file != nil && w.file.Fd() != ^uintptr(0)
}

// maxSize gets max size of the file in bytes.
func (w *FileWriter) maxSize() int64 {
	return int64(*w.MaxSize) * int64(megabyte)
}

// NewFileWriterOptions create a new *FileWriterOptionsBuilder
func NewFileWriterOptions() *FileWriterOptionsBuilder {
	return &FileWriterOptionsBuilder{}
}

// List lists all setter for FileWriterOptions
func (b *FileWriterOptionsBuilder) List() []Setter[FileWriterOptions] {
	return b.Options
}

// SetFilename sets filename
func (b *FileWriterOptionsBuilder) SetFilename(name string) *FileWriterOptionsBuilder {
	b.Options = append(b.Options, func(fwo *FileWriterOptions) error {
		fwo.Filename = &name
		return nil
	})
	return b
}

// SetBackupFormat sets timestamp format for backup filename
func (b *FileWriterOptionsBuilder) SetBackupFormat(format string) *FileWriterOptionsBuilder {
	b.Options = append(b.Options, func(fwo *FileWriterOptions) error {
		fwo.BackupFormat = &format
		return nil
	})
	return b
}

// SetBackupRandomHexSuffix sets enable adding random hex to suffix file name
func (b *FileWriterOptionsBuilder) SetBackupRandomHexSuffix(enable bool) *FileWriterOptionsBuilder {
	b.Options = append(b.Options, func(fwo *FileWriterOptions) error {
		fwo.BackupRandomHexSuffix = &enable
		return nil
	})
	return b
}

// SetMaxSize sets max size per log file
func (b *FileWriterOptionsBuilder) SetMaxSize(size uint32) *FileWriterOptionsBuilder {
	b.Options = append(b.Options, func(fwo *FileWriterOptions) error {
		fwo.MaxSize = &size
		return nil
	})
	return b
}
