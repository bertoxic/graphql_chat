package utils

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrValidation           = errors.New("validation error")
	ErrInternalServer       = errors.New("internal server error")
	ErrUserExist            = errors.New("this user already exists")
	ErrAuthentication       = errors.New("authentication failed")
	ErrUserNotFound         = errors.New("user not found")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrResourceNotFound     = errors.New("resource not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrTokenExpired         = errors.New("token expired")
	ErrTokenInvalid         = errors.New("invalid token")
	ErrDatabase             = errors.New("database error")
	ErrServiceUnavailable   = errors.New("service unavailable")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrConflict             = errors.New("conflict error")
	ErrBadRequest           = errors.New("bad request")
	ErrTimeout              = errors.New("request timeout")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	ErrMethodNotAllowed     = errors.New("method not allowed")
	ErrTooManyRequests      = errors.New("too many requests")
	ErrNotImplemented       = errors.New("not implemented")
)

// ErrorCode represents a unique error code
type ErrorCode int

// Define error codes
const (
	ErrCodeInternal ErrorCode = iota + 1000
	ErrCodeBadRequest
	ErrCodeUnauthorized
	ErrCodeForbidden
	ErrCodeNotFound
	ErrCodeConflict
	ErrCodeRateLimit
	ErrCodeValidation
	ErrCodeAuthentication
	ErrCodeDatabase
	ErrCodeExternalService
)

// AppError represents an application error
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("error code %d: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("error code %d: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e AppError) Unwrap() error {
	return e.Err
}

// Define error variables with codes
//var (
//	ErrInternalServer = NewAppError(ErrCodeInternal, "internal server error", nil)
//	ErrBadRequest     = NewAppError(ErrCodeBadRequest, "bad request", nil)
//	ErrUnauthorized   = NewAppError(ErrCodeUnauthorized, "unauthorized", nil)
//	ErrForbidden      = NewAppError(ErrCodeForbidden, "forbidden", nil)
//	ErrNotFound       = NewAppError(ErrCodeNotFound, "not found", nil)
//	ErrConflict       = NewAppError(ErrCodeConflict, "conflict", nil)
//	ErrRateLimit      = NewAppError(ErrCodeRateLimit, "rate limit exceeded", nil)
//)

// NewAppError creates a new AppError
func NewAppError(code ErrorCode, message string, err error) AppError {
	return AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewValidationError creates a new validation AppError
func NewValidationError(field, message string) AppError {
	return AppError{
		Code:    ErrCodeValidation,
		Message: fmt.Sprintf("validation error: %s - %s", field, message),
	}
}

// NewAuthenticationError creates a new authentication AppError
func NewAuthenticationError(message string) AppError {
	return AppError{
		Code:    ErrCodeAuthentication,
		Message: fmt.Sprintf("authentication error: %s", message),
	}
}

// NewResourceNotFoundError creates a new resource not found AppError
func NewResourceNotFoundError(resource, id string) AppError {
	return AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s with ID %s not found", resource, id),
	}
}

// NewDatabaseError creates a new database AppError
func NewDatabaseError(operation string, err error) AppError {
	return AppError{
		Code:    ErrCodeDatabase,
		Message: fmt.Sprintf("database error during %s", operation),
		Err:     err,
	}
}

// NewExternalServiceError creates a new external service AppError
func NewExternalServiceError(service string, err error) AppError {
	return AppError{
		Code:    ErrCodeExternalService,
		Message: fmt.Sprintf("external service error (%s)", service),
		Err:     err,
	}
}

// IsAppError checks if the error is an AppError
func IsAppError(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr)
}

// GetErrorCode returns the error code of an AppError
func GetErrorCode(err error) ErrorCode {
	var appErr AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return ErrCodeInternal
}

// GetHTTPStatus returns an appropriate HTTP status code for an AppError
func GetHTTPStatus(err error) int {
	code := GetErrorCode(err)
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrCodeValidation:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
