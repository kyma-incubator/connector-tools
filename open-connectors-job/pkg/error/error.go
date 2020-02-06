package error

import (
	"fmt"
	"runtime/debug"
)

type Error struct {
	Inner error
	Message string
	StackTrace string
	Misc map[string]interface{}
}

func (err Error) Error() string {
	return err.Message
}

func WrapError(err error, messagef string, msgArgs ...interface{}) Error {
	return Error{
		Inner: err,
		Message: fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc: make(map[string]interface{}),
	}
}