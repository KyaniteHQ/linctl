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
	require.ErrorIs(t, err, ErrTargetNotConfigured)

	_, err = ResolveTarget(context.Background(), fakeGraphQLClient{
		"Viewer": `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":  "",
	}, matchingTarget())
	require.ErrorContains(t, err, "resolve teams")

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
	require.ErrorContains(t, err, "resolve project")

	graphqlClient := fakeGraphQLClient{
		"Viewer":        `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":         `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"TargetProject": `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"other-team","key":"ABC","name":"other","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`,
	}

	_, err = ResolveTarget(context.Background(), graphqlClient, matchingTarget())
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = listTeamWorkflowStates(context.Background(), fakeGraphQLClient{
		"WorkflowStatesByTeam": workflowStatesByTeamJSON(""),
	}, "team-id")
	require.NoError(t, err)
	_, err = selectWorkflowStateID(nil, "completed")
	require.ErrorIs(t, err, ErrWriteInvalid)

	err = requireTargetMatch(config.Target{OrgID: "org-id", TeamID: "team-id", TeamKey: "LIT"}, config.Target{
		OrgID:   "other-org",
		TeamID:  "team-id",
		TeamKey: "LIT",
	})
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_WriteGuardScenarios_refuse_mismatched_resources(t *testing.T) {
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
	guard := guardedClient{
		graphqlClient: graphqlClient,
		target: ResolvedTarget{
			Team: TargetTeam{ID: "team-id", Key: "LIT"},
		},
	}

	_, err := guard.requireIssue(context.Background(), "ABC-1")
	require.ErrorIs(t, err, ErrTargetMismatch)

	err = guard.requireProject(context.Background(), "project-id")
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = newGuardedClient(context.Background(), errorGraphQLClient{err: errors.New("resolve failed")}, matchingTarget())
	require.ErrorContains(t, err, "resolve failed")

	guard.graphqlClient = errorGraphQLClient{err: errors.New("read issue failed")}
	_, err = guard.requireIssue(context.Background(), "LIT-1")
	require.ErrorContains(t, err, "read issue failed")

	guard.graphqlClient = errorGraphQLClient{err: errors.New("read project failed")}
	err = guard.requireProject(context.Background(), "project-id")
	require.ErrorContains(t, err, "read project failed")
}

func Test_GuardedWrites_return_target_resolution_errors(t *testing.T) {
	graphqlClient := errorGraphQLClient{err: errors.New("resolve failed")}
	testCases := []struct {
		name string
		run  func() error
	}{
		{
			name: "remove issue label",
			run: func() error {
				_, err := RemoveIssueLabel(context.Background(), graphqlClient, matchingTarget(), IssueLabelAssociationRequest{
					IssueID: "LIT-1",
					LabelID: "label-id",
				})
				return err
			},
		},
		{
			name: "update label",
			run: func() error {
				_, err := UpdateLabel(context.Background(), graphqlClient, matchingTarget(), LabelUpdateRequest{
					ID: "label-id", Name: "updated",
				})
				return err
			},
		},
		{
			name: "retire label",
			run: func() error {
				_, err := RetireLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", false)
				return err
			},
		},
		{
			name: "restore label",
			run: func() error {
				_, err := RestoreLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", false)
				return err
			},
		},
		{
			name: "remove project label",
			run: func() error {
				_, err := RemoveProjectLabel(
					context.Background(), graphqlClient, matchingTarget(),
					ProjectLabelAssociationRequest{ProjectID: "project-id", LabelID: "label-id"},
				)
				return err
			},
		},
		{
			name: "update project label",
			run: func() error {
				_, err := UpdateProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelUpdateRequest{
					ID: "label-id", Name: "updated", OrgWide: true,
				})
				return err
			},
		},
		{
			name: "retire project label",
			run: func() error {
				_, err := RetireProjectLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", true)
				return err
			},
		},
		{
			name: "restore project label",
			run: func() error {
				_, err := RestoreProjectLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", true)
				return err
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			require.ErrorContains(t, testCase.run(), "resolve failed")
		})
	}
}

func Test_FakeGraphQLClient_respects_context_and_missing_operations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := fakeGraphQLClient{}.MakeRequest(ctx, &graphql.Request{OpName: "Viewer"}, &graphql.Response{})
	require.Error(t, err)

	err = fakeGraphQLClient{}.MakeRequest(context.Background(), &graphql.Request{OpName: "Viewer"}, &graphql.Response{})
	require.ErrorContains(t, err, "missing fake response")
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

	guard, err := newGuardedClient(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
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
