package client

import (
	"context"
	"encoding/json"
	"errors"
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

type fakeGraphQLClient map[string]string

func (client fakeGraphQLClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	payload, ok := client[request.OpName]
	if !ok {
		return errors.New("missing fake response for " + request.OpName)
	}

	wrapped := []byte(`{"data":` + payload + `}`)
	return json.Unmarshal(wrapped, response)
}
