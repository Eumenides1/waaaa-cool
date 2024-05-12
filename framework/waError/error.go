package waError

import (
	"errors"
	"google.golang.org/grpc/status"
)

type Error struct {
	Code int
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func NewError(code int, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func ToError(err error) *Error {
	fromError, _ := status.FromError(err)
	return NewError(int(fromError.Code()), errors.New(fromError.Message()))
}
