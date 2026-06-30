package oauth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratePKCEReturnsRandomReadError(t *testing.T) {
	readErr := errors.New("random source failed")
	originalRandRead := randRead
	randRead = func([]byte) (int, error) {
		return 0, readErr
	}
	t.Cleanup(func() {
		randRead = originalRandRead
	})

	_, err := GeneratePKCE()

	require.ErrorIs(t, err, readErr)
	require.ErrorContains(t, err, "generate pkce verifier")
}
