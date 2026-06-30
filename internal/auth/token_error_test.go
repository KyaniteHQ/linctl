package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewError_defaults_code_and_formats_message(t *testing.T) {
	t.Parallel()

	err := NewError("", "run linctl auth login")
	require.Equal(t, ErrorCodeReauthRequired, err.Code)
	require.Equal(t, "AUTH_REAUTH_REQUIRED: run linctl auth login", err.Error())

	err = NewError(ErrorCodeNotConfigured, "")
	require.Equal(t, "AUTH_NOT_CONFIGURED", err.Error())

	var nilErr *AuthError
	require.Empty(t, nilErr.Error())
	require.NoError(t, nilErr.Unwrap())
}

func Test_WrapError_preserves_cause(t *testing.T) {
	t.Parallel()
	cause := errors.New("refresh failed")

	err := WrapError("", "refresh OAuth token", cause)

	require.Equal(t, ErrorCodeReauthRequired, err.Code)
	require.Equal(t, "AUTH_REAUTH_REQUIRED: refresh OAuth token", err.Error())
	require.ErrorIs(t, err, cause)
	require.Same(t, cause, err.Unwrap())
}

func Test_NewTokenEndpointError_defaults_code_and_formats_status(t *testing.T) {
	t.Parallel()

	err := NewTokenEndpointError("", 401, "invalid_grant")
	require.Equal(t, ErrorCodeReauthRequired, err.Code)
	require.Equal(t, 401, err.StatusCode)
	require.Equal(t, "invalid_grant", err.OAuthError)
	require.Equal(t, "AUTH_REAUTH_REQUIRED: oauth token endpoint failed (http status 401)", err.Error())

	err = NewTokenEndpointError(ErrorCodeRefreshFailed, 0, "")
	require.Equal(t, "AUTH_REFRESH_FAILED: oauth token endpoint failed", err.Error())

	var nilErr *TokenEndpointError
	require.Empty(t, nilErr.Error())
}
