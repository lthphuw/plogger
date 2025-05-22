package plogger

// DiscardWriter is a Writer implementation that discards all written data.
type DiscardWriter struct{}

// NewDiscardWriter creates a new DiscardWriter.
func NewDiscardWriter() (*DiscardWriter, error) {
	return &DiscardWriter{}, nil
}

// Write implements the io.Writer interface.
// It pretends to write the data by returning its length without doing anything.
func (DiscardWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
