package errcode

import (
	"errors"
	"fmt"
	"strconv"
)

type Error struct {
	code           int
	httpStatusCode int
	message        string
	args           []any
}

var codes = map[int]any{}

// NewError create new *errcode.Error with unique code
//
// NewError panic if found duplicate code
func NewError(code, httpStatusCode int, message string, args ...any) *Error {
	if _, ok := codes[code]; ok {
		panic("errcode: duplicate code " + strconv.Itoa(code))
	}

	codes[code] = nil

	return &Error{
		code:           code,
		httpStatusCode: httpStatusCode,
		message:        message,
		args:           args,
	}
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) HTTPStatusCode() int {
	return e.httpStatusCode
}

func (e *Error) Message() string {
	if len(e.args) > 0 {
		return fmt.Sprintf(e.message, e.args...)
	}

	return e.message
}

func (e *Error) WithArgs(args ...any) *Error {
	err := *e
	err.args = append(err.args, args...)

	return &err
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d - %s", e.Code(), e.Message())
}

// Unwrap recursion unwraps an error and returns *errcode.Error
//
// Unwrap returns inputted error if not found any *errcode.Error
func Unwrap(err error) error {
	for u := err; u != nil; u = errors.Unwrap(u) {
		if ec, ok := u.(*Error); ok {
			return ec
		}
	}

	return err
}

// HardUnwrap is errcode.Unwrap but returns *errcode.Error instead of error
func HardUnwrap(err error) (*Error, bool) {
	for u := err; u != nil; u = errors.Unwrap(u) {
		if ec, ok := u.(*Error); ok {
			return ec, true
		}
	}

	return nil, false
}
