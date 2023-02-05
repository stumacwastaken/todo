package errors

import (
	"fmt"
)

//TodoError is the generic catchall error for the todo service. This means that we're coupling our code to http statuses
//but to be honest I'm ok with that for now. We could refactor that later for true separation of concerns if needed
type TodoError struct {
	msg      string
	HttpCode int
	details  string
}

func (e *TodoError) Error() string {
	return fmt.Sprintf( //quickhack
		`{
			"message": "%s",
			"code": %d,
			"details": "%s"
		}`,
		e.msg, e.HttpCode, e.details,
	)
}

// ErrorWithCode is the general factory method to create new TodoErrors. Requires a message, additional details if needed
// and an http based status code.
func ErrorWithCode(msg, details string, code int) *TodoError {
	return &TodoError{
		msg:      msg,
		details:  details,
		HttpCode: code,
	}
}

func InternalError() *TodoError { //replace this to error compadre....
	return ErrorWithCode("internal error", "something's wrong on our end", 500)
}

func UnknownError() error {
	return ErrorWithCode("unknown", "unknown error occurred", 500)
}
