package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ClientWriteScenarios_guard_writes_and_report_results(t *testing.T) {
	// Given
	t.Run("invalid requests fail before network", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{})

		_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{ID: "LIT-1"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
			ID:          "LIT-1",
			Description: "description",
			Append:      "append",
		})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{Title: "missing id"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{ID: "LIT-1"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{Body: "body"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = StartIssue(context.Background(), graphqlClient, matchingTarget(), "")
		require.Error(t, err)

		_, err = CloseIssue(context.Background(), graphqlClient, matchingTarget(), "")
		require.Error(t, err)

		_, err = CreateProject(context.Background(), graphqlClient, matchingTarget(), ProjectCreateRequest{})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{ID: "project-id"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CreateProjectMilestone(context.Background(), graphqlClient, matchingTarget(), ProjectMilestoneCreateRequest{
			Name: "name",
		})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CreateProjectMilestone(context.Background(), graphqlClient, matchingTarget(), ProjectMilestoneCreateRequest{
			ProjectID: "project-id",
		})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateProjectMilestone(
			context.Background(),
			graphqlClient,
			matchingTarget(),
			ProjectMilestoneUpdateRequest{ID: "project-milestone-id"},
		)
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CreateCycle(context.Background(), graphqlClient, matchingTarget(), CycleCreateRequest{
			EndsAt: "2026-07-15T00:00:00Z",
		})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CreateCycle(context.Background(), graphqlClient, matchingTarget(), CycleCreateRequest{
			StartsAt: "2026-07-01T00:00:00Z",
		})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateCycle(context.Background(), graphqlClient, matchingTarget(), CycleUpdateRequest{ID: "cycle-id"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = ArchiveCycle(context.Background(), graphqlClient, matchingTarget(), "")
		require.ErrorIs(t, err, ErrWriteInvalid)
	})

	t.Run("issue comment succeeds", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-12",
				Title:      "comment target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueCommentCreate": `{"commentCreate":{"success":true,"comment":{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","issue":` + issueJSON(issueFixture{
				Identifier: "LIT-12",
				Title:      "comment target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}}}`,
		})

		comment, err := CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{
			ID:   "LIT-12",
			Body: "hello",
		})

		require.NoError(t, err)
		require.Equal(t, "comment-id", comment.ID)
		require.Equal(t, "LIT-12", comment.Issue.Identifier)
	})

	t.Run("issue update succeeds", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-21",
				Title:      "update target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-21",
				Title:      "updated",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}}`,
		})

		issue, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
			ID:    "LIT-21",
			Title: "updated",
		})

		require.NoError(t, err)
		require.Equal(t, "updated", issue.Title)
	})

	t.Run("issue update appends to description", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSONWithDescription(issueFixture{
				Identifier: "LIT-22",
				Title:      "append target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}, "Existing description") + `}`,
			"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-22",
				Title:      "append target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}}`,
		})

		issue, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
			ID:     "LIT-22",
			Append: "Progress note",
		})

		require.NoError(t, err)
		require.Equal(t, "append target", issue.Title)
		require.Equal(t, "Progress note", appendIssueDescription("", "Progress note"))
		require.Equal(t, "Existing description\n\nProgress note", appendIssueDescription("Existing description\n", "Progress note"))
	})

	t.Run("project update and archive succeed", func(t *testing.T) {
		graphqlClient := projectWriteFakeClient(map[string]string{
			"project": `{"project":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "fixture",
				Status: "Backlog",
			}) + `}`,
			"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "updated",
				Status: "Started",
			}) + `}}`,
			"ProjectArchive": `{"projectArchive":{"success":true,"entity":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "updated",
				Status: "Canceled",
			}) + `}}`,
		})

		project, err := UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{
			ID:   "project-id",
			Name: "updated",
		})
		require.NoError(t, err)
		require.Equal(t, "updated", project.Name)

		project, err = ArchiveProject(context.Background(), graphqlClient, matchingTarget(), "project-id")
		require.NoError(t, err)
		require.Equal(t, "Canceled", project.Status.Name)
	})
}
