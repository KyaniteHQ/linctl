package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_ResolveTarget_confirms_expected_team_and_project_when_token_matches(t *testing.T) {
	// Given
	graphqlClient := fakeGraphQLClient{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"Teams": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}`,
		"TargetProject": `{
			"project": {
				"id": "project-id",
				"name": "fixture",
				"teams": {
					"nodes": [{
						"id": "team-id",
						"key": "LIT",
						"name": "linctl-it",
						"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
					}]
				}
			}
		}`,
	}

	// When
	target, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:     "org-id",
		TeamKey:   "LIT",
		TeamID:    "team-id",
		ProjectID: "project-id",
	})

	// Then
	require.NoError(t, err)
	require.True(t, target.Confirmed)
	require.Equal(t, "Omer", target.Viewer.Name)
	require.Equal(t, "org-id", target.Org.ID)
	require.Equal(t, "LIT", target.Team.Key)
	require.Equal(t, "project-id", target.Project.ID)
}

func Test_ResolveTarget_refuses_when_expected_team_is_missing(t *testing.T) {
	// Given
	graphqlClient := fakeGraphQLClient{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"Teams": `{"teams":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	// When
	_, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_ResolveTarget_finds_expected_team_on_later_page(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"Teams": `{"teams":{"nodes":[],"pageInfo":{"hasNextPage":true,"endCursor":"cursor-1"}}}`,
		"Teams:cursor-1": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}`,
	}

	target, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	require.NoError(t, err)
	require.Equal(t, "team-id", target.Team.ID)
}

func Test_ResolveTarget_refuses_next_page_without_cursor(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"Teams": `{"teams":{"nodes":[],"pageInfo":{"hasNextPage":true,"endCursor":null}}}`,
	}

	_, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
	require.ErrorContains(t, err, "without endCursor")
}

func Test_ResolveTarget_finds_project_team_on_later_page(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"Teams": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}`,
		"TargetProject": `{
			"project": {
				"id": "project-id",
				"name": "fixture",
				"teams": {
					"nodes": [],
					"pageInfo": {"hasNextPage": true, "endCursor": "project-cursor-1"}
				}
			}
		}`,
		"TargetProject:project-cursor-1": `{
			"project": {
				"id": "project-id",
				"name": "fixture",
				"teams": {
					"nodes": [{
						"id": "team-id",
						"key": "LIT",
						"name": "linctl-it",
						"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
					}],
					"pageInfo": {"hasNextPage": false, "endCursor": null}
				}
			}
		}`,
	}

	target, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:     "org-id",
		TeamKey:   "LIT",
		TeamID:    "team-id",
		ProjectID: "project-id",
	})

	require.NoError(t, err)
	require.Equal(t, "project-id", target.Project.ID)
}

func Test_ResolveTarget_uses_direct_team_lookup_when_available(t *testing.T) {
	// Given: no "Teams" fixture entry, so a fallback scan would fail with a
	// missing-fixture error and fail the test.
	graphqlClient := fakeGraphQLClient{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"team": `{
			"team": {
				"id": "team-id",
				"key": "LIT",
				"name": "linctl-it",
				"description": null,
				"archivedAt": null,
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
	}

	// When
	target, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	// Then
	require.NoError(t, err)
	require.True(t, target.Confirmed)
	require.Equal(t, "team-id", target.Team.ID)
	require.Equal(t, "LIT", target.Team.Key)
}

func Test_ResolveTarget_falls_back_to_team_scan_when_direct_lookup_fails(t *testing.T) {
	// Given: no "team" fixture entry, so the direct lookup errors and
	// resolution must fall back to the paged "Teams" scan.
	graphqlClient := fakeGraphQLClient{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"Teams": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}`,
	}

	// When
	target, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "team-id", target.Team.ID)
}

func Test_ResolveTarget_refuses_project_next_page_without_cursor(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"Teams": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}`,
		"TargetProject": `{
			"project": {
				"id": "project-id",
				"name": "fixture",
				"teams": {
					"nodes": [],
					"pageInfo": {"hasNextPage": true, "endCursor": null}
				}
			}
		}`,
	}

	_, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:     "org-id",
		TeamKey:   "LIT",
		TeamID:    "team-id",
		ProjectID: "project-id",
	})

	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
	require.ErrorContains(t, err, "without endCursor")
}

func Test_ResolveTarget_refuses_when_project_next_page_fetch_fails(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"Teams": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}`,
		"TargetProject": `{
			"project": {
				"id": "project-id",
				"name": "fixture",
				"teams": {
					"nodes": [],
					"pageInfo": {"hasNextPage": true, "endCursor": "project-cursor-1"}
				}
			}
		}`,
		"TargetProject:project-cursor-1": "",
	}

	_, err := ResolveTarget(context.Background(), graphqlClient, config.Target{
		OrgID:     "org-id",
		TeamKey:   "LIT",
		TeamID:    "team-id",
		ProjectID: "project-id",
	})

	require.ErrorContains(t, err, "resolve project")
}

func Test_ResolveViewer_reports_only_the_viewer_and_organization(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"Viewer": `{
			"viewer": {
				"id": "app-user-id",
				"name": "linctl",
				"displayName": "linctl bot",
				"email": "app@example.com",
				"app": true,
				"isAssignable": true,
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
	}

	identity, err := ResolveViewer(context.Background(), graphqlClient)

	require.NoError(t, err)
	require.Equal(t, ViewerIdentity{
		Viewer: TargetViewer{
			ID:          "app-user-id",
			Name:        "linctl",
			DisplayName: "linctl bot",
			Email:       "app@example.com",
			App:         true,
			Assignable:  true,
		},
		Org: TargetOrg{ID: "org-id", Name: "Kyanite", URLKey: "kyanite"},
	}, identity)
}

func Test_ResolveViewer_keeps_the_auth_failure_of_a_failed_viewer_query(t *testing.T) {
	_, err := ResolveViewer(context.Background(), errorGraphQLClient{err: fmt.Errorf("graphql http status 401: %w", ErrAuthFailed)})

	require.ErrorIs(t, err, ErrAuthFailed)
	require.ErrorContains(t, err, "resolve viewer")
}

func Test_ResolveViewer_refuses_a_response_without_viewer_or_organization(t *testing.T) {
	tests := map[string]string{
		"null viewer":           `{"viewer": null}`,
		"viewer without org id": `{"viewer": {"id": "user-id", "organization": {"id": ""}}}`,
	}
	for name, payload := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := ResolveViewer(context.Background(), fakeGraphQLClient{"Viewer": payload})

			require.ErrorContains(t, err, "Linear returned no viewer or organization")
		})
	}
}

type fakeGraphQLClient map[string]string

func (client fakeGraphQLClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	payload, ok := client[requestKey(request)]
	if !ok {
		payload, ok = client[request.OpName]
	}
	if !ok {
		return errors.New("missing fake response for " + requestKey(request))
	}

	wrapped := []byte(`{"data":` + payload + `}`)
	return json.Unmarshal(wrapped, response)
}

func requestKey(request *graphql.Request) string {
	var variables map[string]any
	if request.Variables == nil {
		return request.OpName
	}
	data, err := json.Marshal(request.Variables)
	if err != nil {
		return request.OpName
	}
	if err := json.Unmarshal(data, &variables); err != nil {
		return request.OpName
	}
	after, ok := variables["after"].(string)
	if !ok || after == "" {
		return request.OpName
	}

	return request.OpName + ":" + after
}
