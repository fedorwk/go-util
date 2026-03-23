package buffer

import (
	"context"
	"sync/atomic"
)

type Buffer[T any] interface {
	WriteBuffer[T]
	ReadBuffer[T]
}

type ReadBuffer[T any] interface {
	Read(context.Context) (T, error)
}

type WriteBuffer[T any] interface {
	Write(T) error
	Close() error
}

// New creates a new stream with the given options.
// If no buffer size is provided, an unbuffered channel is used.
// If no overflow behavior is set, WithOverflowBlock is the default.
func New[T any](size int, opts ...option[T]) Buffer[T] {
	// Allocate the stream first!
	s := &bufferBase[T]{
		done: make(chan struct{}),
		buf:  make(chan T, size),
	}

	for _, opt := range opts {
		opt(s)
	}

	// Apply defaults
	if s.process == nil {
		WithOverflowBlock[T]()(s)
	}
	return s
}

const (
	bufferStateOpen = iota
	bufferStateClosed
)

type bufferBase[T any] struct {
	buf     chan T
	done    chan struct{}
	state   atomic.Uint32
	process func(T)
}

func (s *bufferBase[T]) Write(rec T) error {
	select {
	case <-s.done:
		return ErrWriteToClosedBuffer
	default:
	}

	s.process(rec)
	return nil
}

func (s *bufferBase[T]) Read(ctx context.Context) (T, error) {
	var dummy T
	select {
	case val, ok := <-s.buf:
		if ok {
			return val, nil
		}
		return dummy, ErrReadFromClosedBuffer
	case <-ctx.Done():
		return dummy, ctx.Err()
	}
}

func (s *bufferBase[T]) Close() error {
	if !s.state.CompareAndSwap(bufferStateOpen, bufferStateClosed) {
		return ErrRepeatCloseBuffer
	}

	close(s.done)
	close(s.buf)
	return nil
}

// option defines a functional option for configuring a stream.
type option[T any] func(*bufferBase[T])

// WithOverflowBlock makes Send block until the value can be sent or the stream is closed.
// This is the default behaviour.
func WithOverflowBlock[T any]() option[T] {
	return func(s *bufferBase[T]) {
		s.process = func(rec T) {
			select {
			case <-s.done:
				return
			case s.buf <- rec:
			}
		}
	}
}

// WithDropNewest makes Send drop the new value if the buffer is full.
func WithDropNewest[T any]() option[T] {
	return func(s *bufferBase[T]) {
		s.process = func(rec T) {
			select {
			case <-s.done:
				return
			default:
			}

			select {
			case s.buf <- rec:
			default:
				// buffer full – drop newest (do nothing)
			}
		}
	}
}

// WithDropOldest makes Send drop the oldest value in the buffer if it is full,
// then retry sending the new value. It repeats until success or the stream is closed.
func WithDropOldest[T any]() option[T] {
	return func(s *bufferBase[T]) {
		s.process = func(rec T) {
			for {
				select {
				case <-s.done:
					return
				case s.buf <- rec:
					return
				default:
					// Buffer is full – attempt to drop one element
				}

				// Try to receive one element (non‑blocking)
				select {
				case <-s.done:
					return
				case <-s.buf:
					// One element removed – loop and try to send again
				default:
					// No element to drop (buffer became empty between the
					// checks). Loop back to the top – the send will likely
					// succeed now.
				}
			}
		}
	}
}
