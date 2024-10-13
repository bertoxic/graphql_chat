package errorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ErrorCode represents a unique error code
type ErrorCode int

// ErrorSeverity represents the severity level of an error
type ErrorSeverity int

// Define severity levels
const (
	SeverityLow ErrorSeverity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// Define error codes with proper spacing for custom codes
const (
	// System/Internal Errors (1000-1099)
	ErrCodeInternal           ErrorCode = 1000
	ErrCodeDatabase           ErrorCode = 1001
	ErrCodeExternalService    ErrorCode = 1002
	ErrCodeServiceUnavailable ErrorCode = 1003
	ErrCodeTimeout            ErrorCode = 1004
	ErrCodeNotImplemented     ErrorCode = 1005
	ErrCodeSystemOverload     ErrorCode = 1006
	ErrCodeDataCorruption     ErrorCode = 1007
	ErrCodeCacheFailure       ErrorCode = 1008
	ErrCodeConfigError        ErrorCode = 1009

	// Authentication/Authorization Errors (1100-1199)
	ErrCodeUnauthorized       ErrorCode = 1100
	ErrCodeForbidden          ErrorCode = 1101
	ErrCodeTokenExpired       ErrorCode = 1102
	ErrCodeInvalidToken       ErrorCode = 1103
	ErrCodeInvalidCredentials ErrorCode = 1104
	ErrCodeSessionExpired     ErrorCode = 1105
	ErrCodeAccountLocked      ErrorCode = 1106
	ErrCodeInvalidScope       ErrorCode = 1107
	ErrCodeMFARequired        ErrorCode = 1108
	ErrCodePasswordExpired    ErrorCode = 1109
	ErrCodeNoUserIdInContext  ErrorCode = 1110

	// Input/Validation Errors (1200-1299)
	ErrCodeBadRequest       ErrorCode = 1200
	ErrCodeValidation       ErrorCode = 1201
	ErrCodeUnsupportedMedia ErrorCode = 1202
	ErrCodeMethodNotAllowed ErrorCode = 1203
	ErrCodeInvalidFormat    ErrorCode = 1204
	ErrCodeMissingField     ErrorCode = 1205
	ErrCodeInvalidLength    ErrorCode = 1206
	ErrCodeInvalidRange     ErrorCode = 1207
	ErrCodeInvalidEnum      ErrorCode = 1208
	ErrCodeInvalidPattern   ErrorCode = 1209

	// Resource Errors (1300-1399)
	ErrCodeNotFound          ErrorCode = 1300
	ErrCodeConflict          ErrorCode = 1301
	ErrCodeResourceExhausted ErrorCode = 1302
	ErrCodeAlreadyExists     ErrorCode = 1303
	ErrCodeResourceLocked    ErrorCode = 1304
	ErrCodeResourceDeleted   ErrorCode = 1305
	ErrCodeResourceDisabled  ErrorCode = 1306
	ErrCodeQuotaExceeded     ErrorCode = 1307
	ErrCodeResourceExpired   ErrorCode = 1308
	ErrCodeResourceMoved     ErrorCode = 1309

	// Rate Limiting/Throttling Errors (1400-1499)
	ErrCodeRateLimit        ErrorCode = 1400
	ErrCodeTooManyRequests  ErrorCode = 1401
	ErrCodeConcurrencyLimit ErrorCode = 1402
	ErrCodeBandwidthLimit   ErrorCode = 1403
	ErrCodeAPIQuotaExceeded ErrorCode = 1404
	ErrCodeIPBlocked        ErrorCode = 1405
	ErrCodeThrottled        ErrorCode = 1406
	ErrCodeBurst            ErrorCode = 1407

	// Business Logic Errors (1500-1599)
	ErrCodeBusinessRule      ErrorCode = 1500
	ErrCodeInsufficientFunds ErrorCode = 1501
	ErrCodeDependencyFailed  ErrorCode = 1502
	ErrCodeStateMismatch     ErrorCode = 1503
	ErrCodeWorkflowViolation ErrorCode = 1504
	ErrCodeDataIntegrity     ErrorCode = 1505
	ErrCodeOperationOrder    ErrorCode = 1506
	ErrCodeUnsupportedOption ErrorCode = 1507

	// External Service Errors (1600-1699)
	ErrCodeThirdPartyService  ErrorCode = 1600
	ErrCodeUpstreamTimeout    ErrorCode = 1601
	ErrCodeServiceMaintenance ErrorCode = 1606

	// Data Errors (1700-1799)
	ErrCodeDataInconsistency ErrorCode = 1700
	ErrCodeDataMigration     ErrorCode = 1701
	ErrCodeDataArchived      ErrorCode = 1702
	ErrCodeDataPurged        ErrorCode = 1703
	ErrCodeBackupFailure     ErrorCode = 1704
	ErrCodeRestoreFailure    ErrorCode = 1705
	ErrCodeReplicationLag    ErrorCode = 1706

	// Security Errors (1800-1899)
	ErrCodeSecurityViolation ErrorCode = 1800
	ErrCodeCSRF              ErrorCode = 1801
	ErrCodeWeakPassword      ErrorCode = 1806
)

// AppError represents a detailed application error
type AppError struct {
	Code      ErrorCode     `json:"code"`
	Message   string        `json:"message"`
	Details   *ErrorDetails `json:"details,omitempty"`
	Severity  ErrorSeverity `json:"severity"`
	Timestamp time.Time     `json:"timestamp"`
	RequestID string        `json:"request_id,omitempty"`
	Err       error         `json:"-"` // Internal error, not exposed in JSON
}

// ErrorDetails contains detailed information about the error
type ErrorDetails struct {
	File      string                 `json:"file,omitempty"`
	Line      int                    `json:"line,omitempty"`
	Function  string                 `json:"function,omitempty"`
	Stack     string                 `json:"stack,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Operation string                 `json:"operation,omitempty"`
}

// Option represents a functional option for creating errors
type Option func(*AppError)

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap implements the errors.Wrapper interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// MarshalJSON implements the json.Marshaler interface
func (e *AppError) MarshalJSON() ([]byte, error) {
	type Alias AppError
	return json.Marshal(&struct {
		*Alias
		Timestamp string `json:"timestamp"`
	}{
		Alias:     (*Alias)(e),
		Timestamp: e.Timestamp.Format(time.RFC3339),
	})
}

// WithRequestID adds a request ID to the error
func WithRequestID(requestID string) Option {
	return func(e *AppError) {
		e.RequestID = requestID
	}
}

// WithSeverity sets the error severity
func WithSeverity(severity ErrorSeverity) Option {
	return func(e *AppError) {
		e.Severity = severity
	}
}

// WithContext adds context information to the error
func WithContext(ctx map[string]interface{}) Option {
	return func(e *AppError) {
		if e.Details == nil {
			e.Details = &ErrorDetails{}
		}
		e.Details.Context = ctx
	}
}

// WithOperation adds operation information to the error
func WithOperation(operation string) Option {
	return func(e *AppError) {
		if e.Details == nil {
			e.Details = &ErrorDetails{}
		}
		e.Details.Operation = operation
	}
}

// New creates a new AppError with stack trace and options
func New(code ErrorCode, message string, err error, opts ...Option) *AppError {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	details := &ErrorDetails{
		File:     filepath.Base(file),
		Line:     line,
		Function: getFunctionName(fn),
		Stack:    getStackTrace(2), // Skip New and caller
	}

	appErr := &AppError{
		Code:      code,
		Message:   message,
		Details:   details,
		Severity:  getSeverityForCode(code),
		Timestamp: time.Now(),
		Err:       err,
	}

	// Apply options
	for _, opt := range opts {
		opt(appErr)
	}

	return appErr
}

// Helper functions
func getFunctionName(fn *runtime.Func) string {
	if fn == nil {
		return "unknown"
	}
	name := filepath.Base(fn.Name())
	return strings.TrimPrefix(name, ".")
}

func getStackTrace(skip int) string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	lines := strings.Split(stack, "\n")
	if len(lines) > skip*2 {
		lines = lines[skip*2:]
	}
	return strings.Join(lines, "\n")
}

func getSeverityForCode(code ErrorCode) ErrorSeverity {
	switch {
	case code >= 1000 && code < 1100:
		return SeverityCritical
	case code >= 1100 && code < 1200:
		return SeverityHigh
	case code >= 1200 && code < 1300:
		return SeverityMedium
	default:
		return SeverityLow
	}
}

// HTTP status code mapping
func (e *AppError) HTTPStatusCode() int {
	switch e.Code {
	case ErrCodeInternal, ErrCodeDatabase, ErrCodeExternalService:
		return http.StatusInternalServerError
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeBadRequest, ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeRateLimit, ErrCodeTooManyRequests:
		return http.StatusTooManyRequests
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeUnsupportedMedia:
		return http.StatusUnsupportedMediaType
	case ErrCodeMethodNotAllowed:
		return http.StatusMethodNotAllowed
	default:
		return http.StatusInternalServerError
	}
}

// Convenience constructors for common errors
func NewValidationError(field, message string, opts ...Option) *AppError {
	return New(ErrCodeValidation,
		fmt.Sprintf("validation error: %s - %s", field, message),
		nil, opts...)
}

func NewAuthenticationError(message string, opts ...Option) *AppError {
	return New(ErrCodeUnauthorized,
		fmt.Sprintf("authentication error: %s", message),
		nil, opts...)
}

func NewDatabaseError(operation string, err error, opts ...Option) *AppError {
	return New(ErrCodeDatabase,
		fmt.Sprintf("database error during %s", operation),
		err,
		append(opts, WithOperation(operation))...)
}

// Error checking helpers
func Is(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

func GetCode(err error) ErrorCode {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return ErrCodeInternal
}
