package client

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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

func Test_resolveTeamDirect_rejects_mismatched_direct_lookup(t *testing.T) {
	graphqlClient := &strictTargetGraphQLClient{fixtures: []strictTargetFixture{{
		operation: "team",
		variables: `{"id":"team-id"}`,
		payload: `{
			"team": {
				"id": "other-team-id",
				"key": "LIT",
				"name": "other",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
	}}}

	team, ok := resolveTeamDirect(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	require.False(t, ok)
	require.Empty(t, team)
	require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
}

func Test_findResolvedTeam_skips_matching_key_with_mismatched_pinned_team_id(t *testing.T) {
	team, ok := findResolvedTeam([]gql.TeamsTeamsTeamConnectionNodesTeam{{
		Id:  "other-team-id",
		Key: "LIT",
		Organization: gql.TeamsTeamsTeamConnectionNodesTeamOrganization{
			Id: "org-id",
		},
	}}, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	require.False(t, ok)
	require.Empty(t, team)
}

func Test_resolveProject_rejects_mismatched_project_id(t *testing.T) {
	graphqlClient := strictProjectFixtureClient("other-project-id")

	project, ok, err := resolveProject(
		context.Background(),
		graphqlClient,
		config.Target{ProjectID: "project-id"},
		targetFixtureTeam(),
	)

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.ErrorContains(t, err, "expected project_id=project-id resolved project_id=other-project-id")
	require.False(t, ok)
	require.Empty(t, project)
	require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
}

func Test_resolveProject_rejects_empty_project_id(t *testing.T) {
	graphqlClient := strictProjectFixtureClient("")

	project, ok, err := resolveProject(
		context.Background(),
		graphqlClient,
		config.Target{ProjectID: "project-id"},
		targetFixtureTeam(),
	)

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.ErrorContains(t, err, "expected project_id=project-id resolved project_id=")
	require.False(t, ok)
	require.Empty(t, project)
	require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
}

func Test_resolveProject_rejects_mismatched_project_id_on_later_page(t *testing.T) {
	for _, projectID := range []string{"", "other-project-id"} {
		t.Run(projectID, func(t *testing.T) {
			graphqlClient := strictLaterProjectPageFixtureClient(projectID)

			project, ok, err := resolveProject(
				context.Background(),
				graphqlClient,
				config.Target{ProjectID: "project-id"},
				targetFixtureTeam(),
			)

			require.ErrorIs(t, err, ErrTargetMismatch)
			require.ErrorContains(t, err, "expected project_id=project-id resolved project_id="+projectID)
			require.False(t, ok)
			require.Empty(t, project)
			require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
		})
	}
}

func Test_resolvedTargetConfig_uses_resolved_project_id(t *testing.T) {
	resolved := resolvedTargetConfig(
		"org-id",
		targetFixtureTeam(),
		ResolvedProject{ID: "resolved-project-id"},
		true,
	)

	require.Equal(t, "resolved-project-id", resolved.ProjectID)
}

func Test_requireTargetMatch_rejects_unresolved_pinned_project(t *testing.T) {
	expected := config.Target{OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id", ProjectID: "project-id"}
	resolved := resolvedTargetConfig("org-id", targetFixtureTeam(), ResolvedProject{}, false)

	err := requireTargetMatch(expected, resolved)

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func strictProjectFixtureClient(projectID string) *strictTargetGraphQLClient {
	return &strictTargetGraphQLClient{fixtures: []strictTargetFixture{{
		operation: "TargetProject",
		variables: `{"id":"project-id","first":250,"after":null}`,
		payload: fmt.Sprintf(`{
			"project": {
				"id": %q,
				"name": "fixture",
				"teams": {
					"nodes": [{"id": "team-id", "key": "LIT", "name": "linctl-it"}],
					"pageInfo": {"hasNextPage": false, "endCursor": null}
				}
			}
		}`, projectID),
	}}}
}

func strictLaterProjectPageFixtureClient(projectID string) *strictTargetGraphQLClient {
	return &strictTargetGraphQLClient{fixtures: []strictTargetFixture{
		{
			operation: "TargetProject",
			variables: `{"id":"project-id","first":250,"after":null}`,
			payload: `{
				"project": {
					"id": "project-id",
					"name": "fixture",
					"teams": {
						"nodes": [],
						"pageInfo": {"hasNextPage": true, "endCursor": "project-cursor-1"}
					}
				}
			}`,
		},
		{
			operation: "TargetProject",
			variables: `{"id":"project-id","first":250,"after":"project-cursor-1"}`,
			payload: fmt.Sprintf(`{
				"project": {
					"id": %q,
					"name": "fixture",
					"teams": {
						"nodes": [{"id": "team-id", "key": "LIT", "name": "linctl-it"}],
						"pageInfo": {"hasNextPage": false, "endCursor": null}
					}
				}
			}`, projectID),
		},
	}}
}

func targetFixtureTeam() gql.TeamsTeamsTeamConnectionNodesTeam {
	return gql.TeamsTeamsTeamConnectionNodesTeam{
		Id:   "team-id",
		Key:  "LIT",
		Name: "linctl-it",
		Organization: gql.TeamsTeamsTeamConnectionNodesTeamOrganization{
			Id:     "org-id",
			Name:   "Kyanite",
			UrlKey: "kyanite",
		},
	}
}
