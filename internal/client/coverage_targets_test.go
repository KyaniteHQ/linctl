package client

import (
	"context"
	"errors"
	"testing"

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
