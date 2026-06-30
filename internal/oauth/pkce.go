package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const pkceVerifierBytes = 32

var randRead = rand.Read

// PKCE carries the verifier and S256 challenge used by OAuth login.
type PKCE struct {
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

// GeneratePKCE returns a fresh PKCE verifier and S256 challenge.
func GeneratePKCE() (PKCE, error) {
	random := make([]byte, pkceVerifierBytes)
	if _, err := randRead(random); err != nil {
		return PKCE{}, fmt.Errorf("generate pkce verifier: %w", err)
	}

	verifier := base64.RawURLEncoding.EncodeToString(random)

	return PKCE{
		CodeVerifier:        verifier,
		CodeChallenge:       CodeChallengeS256(verifier),
		CodeChallengeMethod: "S256",
	}, nil
}

// CodeChallengeS256 returns the RFC 7636 S256 challenge for a verifier.
func CodeChallengeS256(verifier string) string {
	digest := sha256.Sum256([]byte(verifier))

	return base64.RawURLEncoding.EncodeToString(digest[:])
}
