package client

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

const targetViewerFixture = `{
	"viewer": {
		"id": "user-id",
		"name": "Omer",
		"displayName": "Omer",
		"email": "omer@example.com",
		"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
	}
}`

type strictTargetFixture struct {
	operation string
	variables string
	payload   string
}

type strictTargetGraphQLClient struct {
	fixtures []strictTargetFixture
	next     int
}

func (client *strictTargetGraphQLClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if client.next >= len(client.fixtures) {
		return fmt.Errorf("unexpected GraphQL request: operation=%s", request.OpName)
	}

	fixture := client.fixtures[client.next]
	variables, err := json.Marshal(request.Variables)
	if err != nil {
		return fmt.Errorf("marshal %s variables: %w", request.OpName, err)
	}
	if request.OpName != fixture.operation || string(variables) != fixture.variables {
		return fmt.Errorf(
			"unexpected GraphQL request: got operation=%s variables=%s want operation=%s variables=%s",
			request.OpName,
			variables,
			fixture.operation,
			fixture.variables,
		)
	}
	client.next++

	wrapped := []byte(`{"data":` + fixture.payload + `}`)
	if err := json.Unmarshal(wrapped, response); err != nil {
		return fmt.Errorf("decode %s fixture: %w", request.OpName, err)
	}

	return nil
}

func Test_ResolveTarget_key_only_mismatch_skips_direct_lookup_and_describes_team_key(t *testing.T) {
	graphqlClient := &strictTargetGraphQLClient{fixtures: []strictTargetFixture{
		{operation: "Viewer", variables: "null", payload: targetViewerFixture},
		{
			operation: "Teams",
			variables: `{"first":250,"after":null,"includeArchived":true}`,
			payload:   `{"teams":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		},
	}}

	_, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.ErrorContains(t, err, "cannot reach team_key=LIT")
	require.NotContains(t, err.Error(), "team_id=")
	require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
}

func Test_ResolveTarget_falls_back_to_scan_when_direct_lookup_returns_another_team(t *testing.T) {
	graphqlClient := &strictTargetGraphQLClient{fixtures: []strictTargetFixture{
		{operation: "Viewer", variables: "null", payload: targetViewerFixture},
		directTeamLookupFixture("other-team-id"),
		teamScanFixture("team-id"),
	}}

	target, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	require.NoError(t, err)
	require.Equal(t, "team-id", target.Team.ID)
	require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
}

func Test_ResolveTarget_refuses_scanned_team_with_matching_key_and_other_id(t *testing.T) {
	graphqlClient := &strictTargetGraphQLClient{fixtures: []strictTargetFixture{
		{operation: "Viewer", variables: "null", payload: targetViewerFixture},
		directTeamLookupFixture("other-team-id"),
		teamScanFixture("other-team-id"),
	}}

	_, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.ErrorContains(t, err, "cannot reach team_key=LIT team_id=team-id")
	require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
}

func Test_ResolveTarget_refuses_when_pinned_project_resolves_to_another_id(t *testing.T) {
	graphqlClient := &strictTargetGraphQLClient{fixtures: []strictTargetFixture{
		{operation: "Viewer", variables: "null", payload: targetViewerFixture},
		directTeamLookupFixture("team-id"),
		targetProjectFixture("null", "other-project-id", `[{"id": "team-id", "key": "LIT", "name": "linctl-it"}]`, false),
	}}

	_, err := ResolveTarget(context.Background(), graphqlClient, matchingTarget())

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.ErrorContains(t, err, "expected project_id=project-id resolved project_id=other-project-id")
	require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
}

func Test_ResolveTarget_refuses_when_later_project_page_resolves_to_another_id(t *testing.T) {
	graphqlClient := &strictTargetGraphQLClient{fixtures: []strictTargetFixture{
		{operation: "Viewer", variables: "null", payload: targetViewerFixture},
		directTeamLookupFixture("team-id"),
		targetProjectFixture("null", "project-id", `[]`, true),
		targetProjectFixture(
			`"project-cursor-1"`, "other-project-id", `[{"id": "team-id", "key": "LIT", "name": "linctl-it"}]`, false,
		),
	}}

	_, err := ResolveTarget(context.Background(), graphqlClient, matchingTarget())

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.ErrorContains(t, err, "expected project_id=project-id resolved project_id=other-project-id")
	require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
}

func Test_requireTargetMatch_rejects_unresolved_pinned_project(t *testing.T) {
	expected := config.Target{OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id", ProjectID: "project-id"}
	resolved := config.Target{OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id"}

	err := requireTargetMatch(expected, resolved)

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func directTeamLookupFixture(teamID string) strictTargetFixture {
	return strictTargetFixture{
		operation: "team",
		variables: `{"id":"team-id"}`,
		payload: fmt.Sprintf(`{"team": {
			"id": %q, "key": "LIT", "name": "linctl-it",
			"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
		}}`, teamID),
	}
}

func teamScanFixture(teamID string) strictTargetFixture {
	return strictTargetFixture{
		operation: "Teams",
		variables: `{"first":250,"after":null,"includeArchived":true}`,
		payload: fmt.Sprintf(`{"teams": {
			"nodes": [{
				"id": %q, "key": "LIT", "name": "linctl-it",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}],
			"pageInfo": {"hasNextPage": false, "endCursor": null}
		}}`, teamID),
	}
}

func targetProjectFixture(after string, projectID string, teams string, hasNextPage bool) strictTargetFixture {
	endCursor := "null"
	if hasNextPage {
		endCursor = `"project-cursor-1"`
	}

	return strictTargetFixture{
		operation: "TargetProject",
		variables: `{"id":"project-id","first":250,"after":` + after + `}`,
		payload: fmt.Sprintf(`{"project": {
			"id": %q, "name": "fixture",
			"teams": {"nodes": %s, "pageInfo": {"hasNextPage": %t, "endCursor": %s}}
		}}`, projectID, teams, hasNextPage, endCursor),
	}
}
