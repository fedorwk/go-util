package buffer

import "errors"

var (
	ErrReadFromClosedBuffer = errors.New("Failed to read from closed buffer")
	ErrWriteToClosedBuffer  = errors.New("Failed to write to closed buffer")
	ErrRepeatCloseBuffer    = errors.New("Buffer already closed")
)
