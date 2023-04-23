package errs

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

type StatusCode uint

const (
	OK                 StatusCode = 0
	Cancelled          StatusCode = 1
	Unknown            StatusCode = 2
	InvalidArgument    StatusCode = 3
	DeadlineExceeded   StatusCode = 4
	NotFound           StatusCode = 5
	AlreadyExists      StatusCode = 6
	PermissionDenied   StatusCode = 7
	ResourceExhausted  StatusCode = 8
	FailedPrecondition StatusCode = 9
	Aborted            StatusCode = 10
	OutOfRange         StatusCode = 11
	Unimplemented      StatusCode = 12
	Internal           StatusCode = 13
	Unavailable        StatusCode = 14
	DataLoss           StatusCode = 15
	Unauthenticated    StatusCode = 16
)

func Code(err error) StatusCode {
	switch {
	case err == nil:
		return OK
	case errors.Is(err, context.Canceled):
		return Cancelled
	case errors.Is(err, context.DeadlineExceeded):
		return DeadlineExceeded
	default:
		if codeErr, ok := err.(codeI); ok {
			return codeErr.Code()
		} else {
			return Unknown
		}
	}
}

func Wrap(code StatusCode, err error) error {
	if err == nil {
		return nil
	}

	stackErr := errors.WithStack(err)
	return &codeErr{
		stackErr.(errorI),
		code,
	}
}

func Error(code StatusCode, msg string) error {
	return Wrap(code, errors.New(msg))
}

func Errorf(code StatusCode, format string, a ...interface{}) error {
	return Wrap(code, errors.Errorf(format, a...))
}

type errorI interface {
	error
	fmt.Formatter
}

type codeI interface {
	Code() StatusCode
}

type codeErr struct {
	errorI
	code StatusCode
}

func (c *codeErr) Unwrap() error {
	return c.errorI
}

func (c *codeErr) Cause() error {
	return c.errorI
}

func (c *codeErr) Code() StatusCode {
	return c.code
}
