package utils

import (
	"fmt"
	"runtime"
)

// CustomError represents an enhanced error with file, line, and message
type CustomError struct {
	Message string
	File    string
	Line    int
	Func    string
}

// Error satisfies the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("%s\nIn file: %s at line: %d in function: %s", e.Message, e.File, e.Line, e.Func)
}

// New creates a new CustomError with file and line information
func NewCustomErr(message string) error {
	// Get runtime info for the caller
	pc, file, line, ok := runtime.Caller(1) // Get the caller of this function (1 stack frame up)
	if !ok {
		return fmt.Errorf("could not retrieve runtime information")
	}

	// Get function name from the program counter
	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = fn.Name() // Retrieve the function name
	}

	return &CustomError{
		Message: message,
		File:    file,
		Line:    line,
		Func:    funcName,
	}
}

//// Basic usage
//err := errors.New(errors.ErrCodeNotFound, "user not found", nil)
//
//// With options
//err := errors.New(errors.ErrCodeDatabase,
//"failed to update user",
//originalError,
//errors.WithRequestID("req-123"),
//errors.WithContext(map[string]interface{}{
//"userId": 123,
//"action": "update",
//}),
//errors.WithOperation("UpdateUser"))
//
//// Using convenience constructors
//err := errors.NewValidationError("email", "invalid format")
//
//// Checking errors
//if errors.Is(err, errors.ErrCodeNotFound) {
//// Handle not found error
//}
//
//// Getting HTTP status
//status := err.HTTPStatusCode()
