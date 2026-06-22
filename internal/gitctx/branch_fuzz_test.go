package gitctx

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

// FuzzParseIssueIdentifier asserts the parser never panics and that any
// identifier it returns is a real, well-formed substring of the input that
// re-parses to itself.
func FuzzParseIssueIdentifier(f *testing.F) {
	seeds := []string{
		"",
		"LIT-123",
		"feature/LIT-42-add-thing",
		"no issue reference here",
		"lit-99",
		"ABC-1 DEF-2",
		"LIT-",
		"-5",
		"A1-2",
		"Ω-3",
		"LIT-123-456",
		"\n\tLIT-7\n",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	wellFormed := regexp.MustCompile(`^[A-Z][A-Z0-9]+-[0-9]+$`)

	f.Fuzz(func(t *testing.T, text string) {
		identifier, ok := ParseIssueIdentifier(text)
		if !ok {
			require.Empty(t, identifier)

			return
		}

		require.NotEmpty(t, identifier)
		require.Contains(t, text, identifier)
		require.Regexp(t, wellFormed, identifier)

		again, againOK := ParseIssueIdentifier(identifier)
		require.True(t, againOK)
		require.Equal(t, identifier, again)
	})
}
