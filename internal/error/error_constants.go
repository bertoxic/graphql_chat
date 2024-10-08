package errorx

var (
	// System/Internal Errors
	ErrInternal           = New(ErrCodeInternal, "internal server error", nil, WithSeverity(SeverityCritical))
	ErrDatabase           = New(ErrCodeDatabase, "database operation failed", nil, WithSeverity(SeverityCritical))
	ErrTimeout            = New(ErrCodeTimeout, "operation timed out", nil, WithSeverity(SeverityHigh))
	ErrNotImplemented     = New(ErrCodeNotImplemented, "functionality not implemented", nil, WithSeverity(SeverityMedium))
	ErrServiceUnavailable = New(ErrCodeServiceUnavailable, "service temporarily unavailable", nil, WithSeverity(SeverityCritical))
	ErrSystemOverload     = New(ErrCodeSystemOverload, "system is currently overloaded", nil, WithSeverity(SeverityCritical))

	// Authentication/Authorization Errors
	ErrUnauthorized       = New(ErrCodeUnauthorized, "unauthorized access", nil, WithSeverity(SeverityHigh))
	ErrForbidden          = New(ErrCodeForbidden, "forbidden access", nil, WithSeverity(SeverityHigh))
	ErrTokenExpired       = New(ErrCodeTokenExpired, "authentication token has expired", nil, WithSeverity(SeverityMedium))
	ErrInvalidToken       = New(ErrCodeInvalidToken, "invalid authentication token", nil, WithSeverity(SeverityHigh))
	ErrInvalidCredentials = New(ErrCodeInvalidCredentials, "invalid credentials provided", nil, WithSeverity(SeverityMedium))
	ErrSessionExpired     = New(ErrCodeSessionExpired, "user session has expired", nil, WithSeverity(SeverityMedium))
	ErrAccountLocked      = New(ErrCodeAccountLocked, "account has been locked", nil, WithSeverity(SeverityHigh))
	ErrMFARequired        = New(ErrCodeMFARequired, "multi-factor authentication required", nil, WithSeverity(SeverityMedium))

	// Input/Validation Errors
	ErrBadRequest           = New(ErrCodeBadRequest, "invalid request", nil, WithSeverity(SeverityLow))
	ErrValidation           = New(ErrCodeValidation, "validation failed", nil, WithSeverity(SeverityLow))
	ErrUnsupportedMediaType = New(ErrCodeUnsupportedMedia, "unsupported media type", nil, WithSeverity(SeverityLow))
	ErrMethodNotAllowed     = New(ErrCodeMethodNotAllowed, "method not allowed", nil, WithSeverity(SeverityLow))
	ErrMissingField         = New(ErrCodeMissingField, "required field is missing", nil, WithSeverity(SeverityLow))
	ErrInvalidFormat        = New(ErrCodeInvalidFormat, "invalid data format", nil, WithSeverity(SeverityLow))

	// Resource Errors
	ErrNotFound        = New(ErrCodeNotFound, "resource not found", nil, WithSeverity(SeverityLow))
	ErrConflict        = New(ErrCodeConflict, "resource conflict", nil, WithSeverity(SeverityMedium))
	ErrAlreadyExists   = New(ErrCodeAlreadyExists, "resource already exists", nil, WithSeverity(SeverityLow))
	ErrResourceLocked  = New(ErrCodeResourceLocked, "resource is locked", nil, WithSeverity(SeverityMedium))
	ErrResourceExpired = New(ErrCodeResourceExpired, "resource has expired", nil, WithSeverity(SeverityMedium))
	ErrQuotaExceeded   = New(ErrCodeQuotaExceeded, "resource quota exceeded", nil, WithSeverity(SeverityHigh))

	// Rate Limiting Errors
	ErrRateLimit        = New(ErrCodeRateLimit, "rate limit exceeded", nil, WithSeverity(SeverityMedium))
	ErrTooManyRequests  = New(ErrCodeTooManyRequests, "too many requests", nil, WithSeverity(SeverityMedium))
	ErrConcurrencyLimit = New(ErrCodeConcurrencyLimit, "concurrency limit exceeded", nil, WithSeverity(SeverityHigh))
	ErrAPIQuotaExceeded = New(ErrCodeAPIQuotaExceeded, "API quota exceeded", nil, WithSeverity(SeverityHigh))

	// Business Logic Errors
	ErrBusinessRule      = New(ErrCodeBusinessRule, "business rule violation", nil, WithSeverity(SeverityMedium))
	ErrInsufficientFunds = New(ErrCodeInsufficientFunds, "insufficient funds", nil, WithSeverity(SeverityHigh))
	ErrWorkflowViolation = New(ErrCodeWorkflowViolation, "workflow violation", nil, WithSeverity(SeverityMedium))
	ErrDataIntegrity     = New(ErrCodeDataIntegrity, "data integrity violation", nil, WithSeverity(SeverityHigh))

	// External Service Errors
	ErrThirdPartyService  = New(ErrCodeThirdPartyService, "third-party service error", nil, WithSeverity(SeverityHigh))
	ErrUpstreamTimeout    = New(ErrCodeUpstreamTimeout, "upstream service timeout", nil, WithSeverity(SeverityHigh))
	ErrServiceMaintenance = New(ErrCodeServiceMaintenance, "service under maintenance", nil, WithSeverity(SeverityMedium))

	// Security Errors
	ErrSecurityViolation = New(ErrCodeSecurityViolation, "security violation detected", nil, WithSeverity(SeverityCritical))
	ErrCSRF              = New(ErrCodeCSRF, "CSRF token validation failed", nil, WithSeverity(SeverityHigh))
	ErrWeakPassword      = New(ErrCodeWeakPassword, "password does not meet security requirements", nil, WithSeverity(SeverityMedium))
)
