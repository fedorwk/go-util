package errns

import (
	"fmt"
)

func NewErrorNamespace(name string) ErrorNamespace {
	return namespace(name)
}

type ErrorNamespace interface {
	NewType(name string) ErrorType
}

type ErrorType interface {
	New(msg string) error
	Wrap(err error, msg string) error
	WrapWithNoMessage(err error) error
	Error() string
}

type namespace string

func (ns namespace) NewType(name string) ErrorType {
	return errortype{
		base: fmt.Errorf("%s: %s", ns, name),
	}
}

type errortype struct {
	base error
}

func (et errortype) Error() string {
	return et.base.Error()
}

func (t errortype) New(msg string) error {
	return fmt.Errorf("%w: %s", t, msg)
}

func (t errortype) Wrap(err error, msg string) error {
	return fmt.Errorf("%w: %s: %w", t, msg, err)
}

func (t errortype) WrapWithNoMessage(err error) error {
	return fmt.Errorf("%w: %w", t, err)
}
