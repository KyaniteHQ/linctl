package client

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_TargetFailureScenarios_refuse_unpinned_or_mismatched_targets(t *testing.T) {
	_, err := ResolveTarget(context.Background(), fakeGraphQLClient{}, config.Target{})
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = ResolveTarget(context.Background(), fakeGraphQLClient{
		"Viewer": `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":  "",
	}, matchingTarget())
	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve teams")

	_, err = ResolveTarget(context.Background(), fakeGraphQLClient{
		"Viewer":        `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"other-org","name":"Other","urlKey":"other"}}}`,
		"Teams":         `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"TargetProject": `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`,
	}, matchingTarget())
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = ResolveTarget(context.Background(), fakeGraphQLClient{
		"Viewer":        `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":         `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"TargetProject": "",
	}, matchingTarget())
	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve project")

	graphqlClient := fakeGraphQLClient{
		"Viewer":        `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":         `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"TargetProject": `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"other-team","key":"ABC","name":"other","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`,
	}

	_, err = ResolveTarget(context.Background(), graphqlClient, matchingTarget())
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = firstCompletedStateID(context.Background(), fakeGraphQLClient{
		"CompletedWorkflowStates": `{"workflowStates":{"nodes":[]}}`,
	}, "team-id")
	require.ErrorIs(t, err, ErrWriteInvalid)

	_, err = firstStartedStateID(context.Background(), fakeGraphQLClient{
		"StartedWorkflowStates": `{"workflowStates":{"nodes":[]}}`,
	}, "team-id")
	require.ErrorIs(t, err, ErrWriteInvalid)

	err = requireTargetMatch(config.Target{OrgID: "org-id", TeamID: "team-id", TeamKey: "LIT"}, config.Target{
		OrgID:   "other-org",
		TeamID:  "team-id",
		TeamKey: "LIT",
	})
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_WriteGuardScenarios_refuse_mismatched_resources(t *testing.T) {
	guard := writeGuard{
		target: ResolvedTarget{
			Team: TargetTeam{ID: "team-id", Key: "LIT"},
		},
	}
	graphqlClient := fakeGraphQLClient{
		"issue": `{"issue":` + strings.ReplaceAll(issueJSON(issueFixture{
			Identifier: "ABC-1",
			Title:      "wrong team",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}), `"key":"LIT"`, `"key":"ABC"`) + `}`,
		"project": `{"project":` + strings.ReplaceAll(projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "wrong-team",
			Status: "Backlog",
		}), `"key":"LIT"`, `"key":"ABC"`) + `}`,
	}

	_, err := guard.requireIssue(context.Background(), graphqlClient, "ABC-1")
	require.ErrorIs(t, err, ErrTargetMismatch)

	err = guard.requireProject(context.Background(), graphqlClient, "project-id")
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = newWriteGuard(context.Background(), errorGraphQLClient{err: errors.New("resolve failed")}, matchingTarget())
	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve failed")

	_, err = guard.requireIssue(context.Background(), errorGraphQLClient{err: errors.New("read issue failed")}, "LIT-1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "read issue failed")

	err = guard.requireProject(context.Background(), errorGraphQLClient{err: errors.New("read project failed")}, "project-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "read project failed")
}

func Test_FakeGraphQLClient_respects_context_and_missing_operations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := fakeGraphQLClient{}.MakeRequest(ctx, &graphql.Request{OpName: "Viewer"}, &graphql.Response{})
	require.Error(t, err)

	err = fakeGraphQLClient{}.MakeRequest(context.Background(), &graphql.Request{OpName: "Viewer"}, &graphql.Response{})
	require.Error(t, err)
	require.True(t, errors.Is(err, errors.New("missing fake response for Viewer")) || strings.Contains(err.Error(), "missing fake response"))
}

func Test_TargetScenarios_allow_unpinned_project_and_matching_team(t *testing.T) {
	require.Nil(t, optionalString(""))
	require.Equal(t, "value", *optionalString("value"))
	require.Equal(t, "value", *stringPtr("value"))
	require.Equal(t, 7, *intPtr(7))
	require.True(t, *boolPtr(true))
	require.Nil(t, issueDependencyParent(nil))
	require.True(t, projectHasTeam(ProjectSummary{Teams: []ProjectTeam{{ID: "team-id", Key: "LIT"}}}, "team-id", "LIT"))
	require.False(t, projectHasTeam(ProjectSummary{Teams: []ProjectTeam{{ID: "team-id", Key: "ABC"}}}, "team-id", "LIT"))

	guard, err := newWriteGuard(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	require.NoError(t, err)
	require.Nil(t, guard.target.Project)

	err = validateProjectUpdateRequest(ProjectUpdateRequest{Name: "missing id"})
	require.ErrorIs(t, err, ErrWriteInvalid)

	err = validateProjectMilestoneUpdateRequest(ProjectMilestoneUpdateRequest{Name: "missing id"})
	require.ErrorIs(t, err, ErrWriteInvalid)

	err = validateCycleUpdateRequest(CycleUpdateRequest{Name: "missing id"})
	require.ErrorIs(t, err, ErrWriteInvalid)
}
