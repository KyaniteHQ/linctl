package oauth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/oauth"
)

func Test_GeneratePKCE_always_returns_s256_challenge(t *testing.T) {
	t.Parallel()

	pkce := oauth.GeneratePKCE()

	require.Len(t, pkce.CodeVerifier, 43)
	require.NotEqual(t, pkce.CodeVerifier, pkce.CodeChallenge)
	require.Equal(t, "S256", pkce.CodeChallengeMethod)
	require.Equal(t, oauth.CodeChallengeS256(pkce.CodeVerifier), pkce.CodeChallenge)
	require.NotContains(t, pkce.CodeVerifier, "=")
	require.NotContains(t, pkce.CodeChallenge, "=")
}

func Test_CodeChallengeS256_matches_rfc7636_vector(t *testing.T) {
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"

	challenge := oauth.CodeChallengeS256(verifier)

	require.Equal(t, "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM", challenge)
	require.NotContains(t, challenge, "=")
}
