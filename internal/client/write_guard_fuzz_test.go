package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// FuzzRequireTargetMatch is the property form of the fail-closed invariant from
// ADR-0001: requireTargetMatch must accept a resolved target only when its
// org, team id, and team key all equal the pinned target, and otherwise refuse
// with ErrTargetMismatch. The example-based guard tests cover specific
// scenarios; this proves the invariant across the whole input space.
func FuzzRequireTargetMatch(f *testing.F) {
	seeds := []struct {
		expectedOrg, expectedTeamID, expectedTeamKey string
		resolvedOrg, resolvedTeamID, resolvedTeamKey string
	}{
		{"org", "team", "LIT", "org", "team", "LIT"},
		{"org", "team", "LIT", "other", "team", "LIT"},
		{"org", "team", "LIT", "org", "other", "LIT"},
		{"org", "team", "LIT", "org", "team", "OTHER"},
		{"", "", "", "", "", ""},
		{"o", "t", "k", "", "", ""},
	}
	for _, seed := range seeds {
		f.Add(
			seed.expectedOrg, seed.expectedTeamID, seed.expectedTeamKey,
			seed.resolvedOrg, seed.resolvedTeamID, seed.resolvedTeamKey,
		)
	}

	f.Fuzz(func(
		t *testing.T,
		expectedOrg, expectedTeamID, expectedTeamKey string,
		resolvedOrg, resolvedTeamID, resolvedTeamKey string,
	) {
		expected := config.Target{OrgID: expectedOrg, TeamID: expectedTeamID, TeamKey: expectedTeamKey}
		resolved := config.Target{OrgID: resolvedOrg, TeamID: resolvedTeamID, TeamKey: resolvedTeamKey}

		err := requireTargetMatch(expected, resolved)

		matches := resolvedOrg == expectedOrg &&
			resolvedTeamID == expectedTeamID &&
			resolvedTeamKey == expectedTeamKey
		if matches {
			require.NoError(t, err)

			return
		}

		require.ErrorIs(t, err, ErrTargetMismatch)
	})
}

func Test_guardedMutation_returns_guard_resolution_error(t *testing.T) {
	_, err := guardedMutation(
		context.Background(),
		errorGraphQLClient{err: errors.New("resolve failed")},
		matchingTarget(),
		func(writeGuard) (IssueSummary, error) {
			t.Fatal("mutation ran before guard resolved")
			return IssueSummary{}, nil
		},
	)

	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve failed")
}
