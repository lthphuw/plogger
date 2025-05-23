package plogger_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/lthphuw/plogger"
)

func TestFileWriterConcurrent(t *testing.T) {
	tmpDir := "./tmp"
	logFile := filepath.Join(tmpDir, "test.log")

	// We set SetBackupRandomHexSuffix true,
	// because the timestamp can not fast enough to generate unique file
	opts := plogger.NewFileWriterOptions().
		SetFilename(logFile).
		SetMaxSize(1).
		SetBackupRandomHexSuffix(true)
	writer, err := plogger.NewFileWriter(opts)
	if err != nil {
		t.Fatalf("NewFileWriter() error: %v", err)
	}

	const numGoroutines = 10_000
	const writesPerGoroutine = 100

	data := []byte(
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam ut euismod nunc. Mauris nec tristique lacus. Sed mollis, tellus id aliquam molestie, risus libero placerat felis, vel sagittis mauris lacus ut ex. Fusce placerat neque vel est rhoncus, non interdum libero hendrerit. Duis sit amet pharetra nisi. Phasellus a volutpat sem. Nullam pretium nisi a maximus volutpat. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculoustouch .git/hooks/pre-commit
 mus. Fusce interdum arcu nec scelerisque cursus.`,
	)

	// Expected Bytes that we want to write
	expectedBytes := numGoroutines * writesPerGoroutine * len(data)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for range numGoroutines {
		go func() {
			defer wg.Done()
			for j := 0; j < writesPerGoroutine; j++ {
				n, err := writer.Write(data)
				if err != nil {
					t.Errorf("Write() error: %v", err)
					return
				}
				if n != len(data) {
					t.Errorf("Write() wrote %d bytes, want %d", n, len(data))
					return
				}
			}
		}()
	}

	wg.Wait()

	// Close writer
	if err := writer.Close(); err != nil {
		t.Errorf("Close() error: %v", err)
	}

	// Check total bytes we writes
	var totalSize int64
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Walk() error: %v", err)
	}

	if totalSize != int64(expectedBytes) {
		t.Errorf("Total written bytes = %d, want %d", totalSize, expectedBytes)
	}

	// Clean
	if err := os.RemoveAll(tmpDir); err != nil {
	}
}
