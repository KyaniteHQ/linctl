package auth

import "fmt"

// ErrorCode is a stable machine-readable OAuth failure reason.
type ErrorCode string

// Error codes classify OAuth auth/readiness failures for command output.
const (
	ErrorCodeNotConfigured  ErrorCode = "AUTH_NOT_CONFIGURED"
	ErrorCodeTokenExpired   ErrorCode = "AUTH_TOKEN_EXPIRED" //nolint:gosec // OAuth error code label, not a credential.
	ErrorCodeRefreshFailed  ErrorCode = "AUTH_REFRESH_FAILED"
	ErrorCodeReauthRequired ErrorCode = "AUTH_REAUTH_REQUIRED"
	ErrorCodeMissingScope   ErrorCode = "MISSING_SCOPE"
	ErrorCodeActorMismatch  ErrorCode = "AUTH_ACTOR_MISMATCH"
	ErrorCodeTargetMismatch ErrorCode = "AUTH_TARGET_MISMATCH"
)

// TokenEndpointError reports a token endpoint failure without secret values.
type TokenEndpointError struct {
	Code       ErrorCode
	StatusCode int
	OAuthError string
}

// AuthError reports a stable OAuth auth/readiness failure without secrets.
//
//nolint:revive // AuthError is the current errors.As type used by CLI callers.
type AuthError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// NewError returns a structured OAuth auth/readiness error.
func NewError(code ErrorCode, message string) *AuthError {
	if code == "" {
		code = ErrorCodeReauthRequired
	}

	return &AuthError{Code: code, Message: message}
}

// WrapError returns a structured OAuth auth/readiness error with a cause.
func WrapError(code ErrorCode, message string, err error) *AuthError {
	authErr := NewError(code, message)
	authErr.Err = err

	return authErr
}

func (err *AuthError) Error() string {
	if err == nil {
		return ""
	}
	if err.Message == "" {
		return string(err.Code)
	}

	return fmt.Sprintf("%s: %s", err.Code, err.Message)
}

func (err *AuthError) Unwrap() error {
	if err == nil {
		return nil
	}

	return err.Err
}

// NewTokenEndpointError returns a redacted structured token endpoint error.
func NewTokenEndpointError(code ErrorCode, statusCode int, oauthError string) *TokenEndpointError {
	if code == "" {
		code = ErrorCodeReauthRequired
	}

	return &TokenEndpointError{
		Code:       code,
		StatusCode: statusCode,
		OAuthError: oauthError,
	}
}

func (err *TokenEndpointError) Error() string {
	if err == nil {
		return ""
	}
	message := fmt.Sprintf("%s: oauth token endpoint failed", err.Code)
	if err.StatusCode != 0 {
		message += fmt.Sprintf(" (http status %d)", err.StatusCode)
	}

	return message
}
