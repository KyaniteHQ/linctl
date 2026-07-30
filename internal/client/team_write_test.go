package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func teamJSON(orgID string) string {
	return `{
		"id": "created-team-id",
		"key": "OPS",
		"name": "Operations",
		"description": "ops",
		"archivedAt": null,
		"organization": {"id": "` + orgID + `", "name": "Kyanite", "urlKey": "kyanite"}
	}`
}

func Test_CreateTeam_creates_team_when_org_wide_passed(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"TeamCreate": `{"teamCreate":{"success":true,"team":` + teamJSON("org-id") + `}}`,
	})

	team, err := CreateTeam(context.Background(), graphqlClient, matchingTarget(), TeamCreateRequest{
		Name: "Operations", Key: "OPS", Description: "ops", Private: true, OrgWide: true,
	})

	require.NoError(t, err)
	require.Equal(t, "created-team-id", team.ID)
	require.Equal(t, "OPS", team.Key)
	require.Equal(t, "org-id", team.OrgID)
}

func Test_CreateTeam_refuses_without_org_wide(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := CreateTeam(context.Background(), recorder, matchingTarget(), TeamCreateRequest{
		Name: "Operations",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("TeamCreate"))
}

func Test_CreateTeam_requires_name(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := CreateTeam(context.Background(), recorder, matchingTarget(), TeamCreateRequest{
		OrgWide: true,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("TeamCreate"))
}

func Test_CreateTeam_refuses_when_target_unresolved(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := CreateTeam(context.Background(), recorder, config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, TeamCreateRequest{Name: "Operations", OrgWide: true})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("TeamCreate"))
}

func Test_CreateTeam_wraps_mutation_error(t *testing.T) {
	_, err := CreateTeam(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		TeamCreateRequest{Name: "Operations", OrgWide: true},
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateTeam_fails_when_mutation_reports_no_success(t *testing.T) {
	_, err := CreateTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"TeamCreate": `{"teamCreate":{"success":false,"team":` + teamJSON("org-id") + `}}`,
	}), matchingTarget(), TeamCreateRequest{Name: "Operations", OrgWide: true})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_CreateTeam_fails_when_mutation_returns_no_team(t *testing.T) {
	// A payload reporting success without the thing it created is not a
	// successful write, and reading it would dereference nil.
	_, err := CreateTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"TeamCreate": `{"teamCreate":{"success":true,"team":null}}`,
	}), matchingTarget(), TeamCreateRequest{Name: "Operations", OrgWide: true})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_CreateTeam_refuses_a_team_created_in_another_organization(t *testing.T) {
	// The Org-Scoped Write comparison is made against what Linear returned, not
	// against target resolution alone, so a team landing elsewhere is a hard stop.
	_, err := CreateTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"TeamCreate": `{"teamCreate":{"success":true,"team":` + teamJSON("other-org-id") + `}}`,
	}), matchingTarget(), TeamCreateRequest{Name: "Operations", OrgWide: true})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateTeam_sends_every_field_it_was_given(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"TeamCreate": `{"teamCreate":{"success":true,"team":` + teamJSON("org-id") + `}}`,
	})}

	_, err := CreateTeam(context.Background(), recorder, matchingTarget(), TeamCreateRequest{
		Name: "Operations", Key: "OPS", Description: "ops", Private: true, OrgWide: true,
	})

	require.NoError(t, err)
	require.JSONEq(t, `{
		"input": {
			"name": "Operations",
			"key": "OPS",
			"description": "ops",
			"private": true
		}
	}`, string(recorder.variablesFor(t, "TeamCreate")))
}

func Test_CreateTeam_omits_the_optional_fields_it_was_not_given(t *testing.T) {
	// Linear derives a key from the name, and an explicit null would override that
	// derivation rather than defer to it, so an unset field must not be sent.
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"TeamCreate": `{"teamCreate":{"success":true,"team":` + teamJSON("org-id") + `}}`,
	})}

	_, err := CreateTeam(context.Background(), recorder, matchingTarget(), TeamCreateRequest{
		Name: "Operations", OrgWide: true,
	})

	require.NoError(t, err)
	require.JSONEq(t, `{"input": {"name": "Operations"}}`, string(recorder.variablesFor(t, "TeamCreate")))
}

func Test_CreateTeam_proceeds_when_pinned_project_present(t *testing.T) {
	// matchingTarget() pins project-id. A Team create compares organization only,
	// so a pinned project must not block it.
	team, err := CreateTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"TeamCreate": `{"teamCreate":{"success":true,"team":` + teamJSON("org-id") + `}}`,
	}), matchingTarget(), TeamCreateRequest{Name: "Operations", OrgWide: true})

	require.NoError(t, err)
	require.Equal(t, "created-team-id", team.ID)
}
