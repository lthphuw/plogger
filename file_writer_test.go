package plogger_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lthphuw/plogger"
	"github.com/lthphuw/plogger/testdata"
)

func TestFileWriter(t *testing.T) {
	testcases := []struct {
		name          string
		opts          *plogger.FileWriterOptionsBuilder
		input         []byte
		writeTimes    int
		removeDir     bool
		invalidDir    bool
		openExistFile bool
		wantNewErr    bool
		wantWriteErr  bool
		wantLength    int
	}{
		{
			name:         "No options",
			input:        []byte("Hello world"),
			opts:         nil,
			writeTimes:   1,
			wantNewErr:   false,
			wantWriteErr: false,
			wantLength:   len("Hello world"),
		},
		{
			name:         "Empty Slice",
			input:        nil,
			opts:         nil,
			writeTimes:   1,
			wantNewErr:   false,
			wantWriteErr: false,
			wantLength:   0,
		},
		{
			name:         "Filename without extension",
			input:        []byte("Hello world"),
			opts:         plogger.NewFileWriterOptions().SetFilename("./logs"),
			wantNewErr:   false,
			wantWriteErr: false,
			writeTimes:   1,
			wantLength:   len("Hello world"),
		},
		{
			name:         "Set backup format",
			input:        []byte("Hello world"),
			opts:         plogger.NewFileWriterOptions().SetBackupFormat("2006-01-02_15-04-05"),
			writeTimes:   1,
			wantNewErr:   false,
			wantWriteErr: false,
			wantLength:   len("Hello world"),
		},
		{
			name:          "Open exists file",
			input:         []byte("Hello world"),
			opts:          plogger.NewFileWriterOptions().SetFilename("./testdata/temp.log"),
			openExistFile: true,
			writeTimes:    1,
			wantNewErr:    false,
			wantWriteErr:  false,
			wantLength:    len("Hello world"),
		},
		{
			name:         "Write byte too long",
			input:        []byte(testdata.MoreThan1MBBytes),
			opts:         plogger.NewFileWriterOptions().SetMaxSize(1).SetFilename("./tmp.log"),
			writeTimes:   1,
			wantNewErr:   false,
			wantWriteErr: true,
			wantLength:   0,
		},
		{
			name:  "Write bytes total exceeds file size",
			input: []byte(testdata.ManyBytes),
			opts: plogger.NewFileWriterOptions().
				SetMaxSize(1).
				SetFilename("./temp/tmp.log").
				SetBackupRandomHexSuffix(false),
			removeDir:    true,
			writeTimes:   100,
			wantNewErr:   false,
			wantWriteErr: false,
			wantLength:   len(testdata.ManyBytes) * 100,
		},
		{
			name:  "Write bytes total exceeds file size",
			input: []byte(testdata.ManyBytes),
			opts: plogger.NewFileWriterOptions().
				SetMaxSize(1).
				SetFilename("./temp/tmp.log").
				SetBackupRandomHexSuffix(true),
			removeDir:    true,
			writeTimes:   100,
			wantNewErr:   false,
			wantWriteErr: false,
			wantLength:   len(testdata.ManyBytes) * 100,
		},
		{
			name:         "Invalid directories",
			input:        []byte("hello"),
			opts:         plogger.NewFileWriterOptions().SetFilename("./temp/tmp.log"),
			removeDir:    true,
			invalidDir:   true,
			writeTimes:   1,
			wantNewErr:   false,
			wantWriteErr: true,
			wantLength:   0,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			writer, err := plogger.NewFileWriter(tt.opts)
			if (err != nil) != tt.wantNewErr {
				t.Errorf("TestFileWriter() error = %v, wantNewErr = %v", err, tt.wantNewErr)
				return
			}
			if err != nil {
				return
			}
			err = nil

			// Create new file
			if tt.openExistFile {
				dir := filepath.Dir(*writer.Filename)
				// Create directories if not exists.
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Errorf("failed to create directories: %v", err)
				}
				_, err := os.OpenFile(*writer.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
				if err != nil {
					t.Errorf("failed to create log file: %v", err)
				}
			}

			if tt.invalidDir {
				if err := createFile(*writer.Filename); err != nil {
				}
			}

			n := 0
			nn := 0
			for range tt.writeTimes {
				nn, err = writer.Write(tt.input)
				n += nn
				if nn == 0 || err != nil {
					n = 0
					break
				}
			}
			if (err != nil) != tt.wantWriteErr {
				t.Errorf("Write() error = %v, wantWriteErr = %v", err, tt.wantWriteErr)
			}
			if n != tt.wantLength {
				t.Errorf("Write() length = %d, want = %d", n, tt.wantLength)
			}
			if err = writer.Close(); err != nil {
				t.Errorf("Close() got error: ")
			}
			if tt.removeDir {
				if err := os.RemoveAll(filepath.Dir(*writer.Filename)); err != nil {
				}
			} else {
				if err := os.Remove(*writer.Filename); err != nil {
				}
			}
		})
	}
}

func createFile(inputPath string) error {
	// Get the parent directory of the input path (abc/cc)
	parentDir := filepath.Dir(inputPath) // Returns "abc/cc"

	// Get the parent of the parent directory (abc)
	targetDir := filepath.Dir(parentDir) // Returns "abc"

	// Get the file name to create (cc)
	fileName := filepath.Base(parentDir) // Returns "cc"

	// Create the target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	// Create the file path (abc/cc)
	targetFile := filepath.Join(targetDir, fileName)

	// Create the file
	f, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", targetFile, err)
	}
	f.Close() // Close the file after creation
	fmt.Printf("Created file: %s\n", targetFile)
	return nil
}
