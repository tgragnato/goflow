package debug

import (
	"errors"
)

var PanicError = errors.New("panic")

type PanicErrorMessage struct {
	Msg        interface{}
	Inner      string
	Stacktrace []byte
}

func (e *PanicErrorMessage) Error() string {
	return e.Inner
}

func (e *PanicErrorMessage) Unwrap() []error {
	return []error{PanicError}
}
