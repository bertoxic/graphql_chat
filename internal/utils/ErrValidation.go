package utils

//
//import (
//	"errors"
//	"fmt"
//	"net/http"
//	"path/filepath"
//	"runtime"
//	"strings"
//)
//
//var (
//	ErrValidation           = errors.New("validation error")
//	ErrInternalServer       = errors.New("internal server error")
//	ErrUserExist            = errors.New("this user already exists")
//	ErrAuthentication       = errors.New("authentication failed")
//	ErrUserNotFound         = errors.New("user not found")
//	ErrPermissionDenied     = errors.New("permission denied")
//	ErrResourceNotFound     = errors.New("resource not found")
//	ErrInvalidCredentials   = errors.New("invalid credentials")
//	ErrTokenExpired         = errors.New("token expired")
//	ErrTokenInvalid         = errors.New("invalid token")
//	ErrDatabase             = errors.New("database error")
//	ErrServiceUnavailable   = errors.New("service unavailable")
//	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
//	ErrConflict             = errors.New("conflict error")
//	ErrBadRequest           = errors.New("bad request")
//	ErrTimeout              = errors.New("request timeout")
//	ErrUnsupportedMediaType = errors.New("unsupported media type")
//	ErrMethodNotAllowed     = errors.New("method not allowed")
//	ErrTooManyRequests      = errors.New("too many requests")
//	ErrNotImplemented       = errors.New("not implemented")
//)
//
//// ErrorCode represents a unique error code
//type ErrorCode int
//
//// Define error codes
//const (
//	ErrCodeInternal ErrorCode = iota + 1000
//	ErrCodeBadRequest
//	ErrCodeUnauthorized
//	ErrCodeForbidden
//	ErrCodeNotFound
//	ErrCodeConflict
//	ErrCodeRateLimit
//	ErrCodeValidation
//	ErrCodeAuthentication
//	ErrCodeDatabase
//	ErrCodeExternalService
//)
//
//// AppError represents an application error
//type ErrDetails struct {
//	Message string
//	File    string
//	Line    int
//	Func    string
//}
//type AppError struct {
//	Code    ErrorCode
//	Message string
//	Details string
//	Err     error
//}
//
//func (e AppError) Error() string {
//	if e.Err != nil {
//		return fmt.Sprintf("error code %d: %s - %v", e.Code, e.Message, e.Err)
//	}
//	return fmt.Sprintf("error code %d: %s", e.Code, e.Message)
//}
//func (e AppError) Detail() string {
//	if e.Err != nil {
//		return fmt.Sprintf("error code %d: %s - %v", e.Code, e.Details, e.Err)
//	}
//	return fmt.Sprintf("error code %d: %s", e.Code, e.Details)
//}
//
//// Unwrap returns the wrapped error
//func (e AppError) Unwrap() error {
//	return e.Err
//}
//
//// Define error variables with codes
//var (
//	ErrInternalServer = NewAppError(ErrCodeInternal, "internal server error", nil)
//	ErrBadRequest     = NewAppError(ErrCodeBadRequest, "bad request", nil)
//	ErrUnauthorized   = NewAppError(ErrCodeUnauthorized, "unauthorized", nil)
//	ErrForbidden      = NewAppError(ErrCodeForbidden, "forbidden", nil)
//	ErrNotFound       = NewAppError(ErrCodeNotFound, "not found", nil)
//	ErrConflict       = NewAppError(ErrCodeConflict, "conflict", nil)
//	ErrRateLimit      = NewAppError(ErrCodeRateLimit, "rate limit exceeded", nil)
//)
//
//// NewAppError creates a new AppError
//func NewAppError(code ErrorCode, message string, err error) AppError {
//	// Get runtime info for the caller
//	pc, file, line, ok := runtime.Caller(1) // Get the caller of this function (1 stack frame up)
//	if !ok {
//		return AppError{
//			Code:    code,
//			Message: message,
//			Details: "could not get details of this error",
//			Err:     err,
//		}
//	}
//	fn := runtime.FuncForPC(pc)
//	funcName := "unknown"
//	if fn != nil {
//		fullFuncName := fn.Name()
//		funcName = filepath.Base(fullFuncName)
//		// If the function name still contains a package path, split by '/' and take the last part
//		parts := strings.Split(funcName, "/")
//		funcName = parts[len(parts)-1]
//		// Retrieve the function name
//	}
//	e := ErrDetails{
//		Message: message,
//		File:    file,
//		Line:    line,
//		Func:    funcName,
//	}
//	detailmsg := fmt.Sprintf(
//		"errmsg: %s,\nIn file: %s\nat line:%d, in function: %s\n",
//		e.Message, e.File, e.Line, e.Func)
//
//	return AppError{
//		Code:    code,
//		Message: message,
//		Details: detailmsg,
//		Err:     err,
//	}
//}
//
//// NewValidationError creates a new validation AppError
//func NewValidationError(field, message string) AppError {
//	return AppError{
//		Code:    ErrCodeValidation,
//		Message: fmt.Sprintf("validation error: %s - %s", field, message),
//	}
//}
//
//// NewAuthenticationError creates a new authentication AppError
//func NewAuthenticationError(message string) AppError {
//	return AppError{
//		Code:    ErrCodeAuthentication,
//		Message: fmt.Sprintf("authentication error: %s", message),
//	}
//}
//
//// NewResourceNotFoundError creates a new resource not found AppError
//func NewResourceNotFoundError(resource, id string) AppError {
//	return AppError{
//		Code:    ErrCodeNotFound,
//		Message: fmt.Sprintf("%s with ID %s not found", resource, id),
//	}
//}
//
//// NewDatabaseError creates a new database AppError
//func NewDatabaseError(operation string, err error) AppError {
//	return AppError{
//		Code:    ErrCodeDatabase,
//		Message: fmt.Sprintf("database error during %s", operation),
//		Err:     err,
//	}
//}
//
//// NewExternalServiceError creates a new external service AppError
//func NewExternalServiceError(service string, err error) AppError {
//	return AppError{
//		Code:    ErrCodeExternalService,
//		Message: fmt.Sprintf("external service error (%s)", service),
//		Err:     err,
//	}
//}
//
//// IsAppError checks if the error is an AppError
//func IsAppError(err error) bool {
//	var appErr AppError
//	return errors.As(err, &appErr)
//}
//
//// GetErrorCode returns the error code of an AppError
//func GetErrorCode(err error) ErrorCode {
//	var appErr AppError
//	if errors.As(err, &appErr) {
//		return appErr.Code
//	}
//	return ErrCodeInternal
//}
//
//// GetHTTPStatus returns an appropriate HTTP status code for an AppError
//func GetHTTPStatus(err error) int {
//	code := GetErrorCode(err)
//	switch code {
//	case ErrCodeBadRequest:
//		return http.StatusBadRequest
//	case ErrCodeUnauthorized:
//		return http.StatusUnauthorized
//	case ErrCodeForbidden:
//		return http.StatusForbidden
//	case ErrCodeNotFound:
//		return http.StatusNotFound
//	case ErrCodeConflict:
//		return http.StatusConflict
//	case ErrCodeRateLimit:
//		return http.StatusTooManyRequests
//	case ErrCodeValidation:
//		return http.StatusUnprocessableEntity
//	default:
//		return http.StatusInternalServerError
//	}
//}
