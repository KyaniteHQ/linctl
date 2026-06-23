package client

import (
	"context"
	"testing"

	"github.com/KyaniteHQ/linctl/internal/config"
	"github.com/stretchr/testify/require"
)

func Test_ClientWriteFailureScenarios_refuse_unpinned_targets(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})
	emptyTarget := config.Target{}

	_, err := CreateIssue(context.Background(), graphqlClient, emptyTarget, IssueCreateRequest{Title: "title"})
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = UpdateIssue(context.Background(), graphqlClient, emptyTarget, IssueUpdateRequest{ID: "LIT-1", Title: "title"})
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = CommentOnIssue(context.Background(), graphqlClient, emptyTarget, IssueCommentRequest{ID: "LIT-1", Body: "body"})
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = StartIssue(context.Background(), graphqlClient, emptyTarget, "LIT-1")
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = CloseIssue(context.Background(), graphqlClient, emptyTarget, "LIT-1")
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = CreateProject(context.Background(), graphqlClient, emptyTarget, ProjectCreateRequest{Name: "name"})
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = UpdateProject(context.Background(), graphqlClient, emptyTarget, ProjectUpdateRequest{ID: "project-id", Name: "name"})
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = ArchiveProject(context.Background(), graphqlClient, emptyTarget, "project-id")
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = CreateProjectMilestone(
		context.Background(),
		graphqlClient,
		emptyTarget,
		ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"},
	)
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = UpdateProjectMilestone(
		context.Background(),
		graphqlClient,
		emptyTarget,
		ProjectMilestoneUpdateRequest{ID: "project-milestone-id", Name: "name"},
	)
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = CreateCycle(
		context.Background(),
		graphqlClient,
		emptyTarget,
		CycleCreateRequest{StartsAt: "2026-07-01T00:00:00Z", EndsAt: "2026-07-15T00:00:00Z"},
	)
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = UpdateCycle(
		context.Background(),
		graphqlClient,
		emptyTarget,
		CycleUpdateRequest{ID: "cycle-id", Name: "name"},
	)
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = ArchiveCycle(context.Background(), graphqlClient, emptyTarget, "cycle-id")
	require.ErrorIs(t, err, ErrTargetMismatch)
}
