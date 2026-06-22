package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func documentJSON(title string, teamID string, teamKey string, projectID string) string {
	return `{
		"id":"document-id",
		"title":"` + title + `",
		"slugId":"doc-slug",
		"archivedAt":null,
		"project":{"id":"` + projectID + `","name":"fixture"},
		"team":{"id":"` + teamID + `","key":"` + teamKey + `","name":"linctl"},
		"issue":null,
		"cycle":null
	}`
}

func Test_CreateDocument_returns_created_document_when_target_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"DocumentCreate": `{"documentCreate":{"success":true,"document":` +
			documentJSON("created", "team-id", "LIT", "project-id") + `}}`,
	})

	document, err := CreateDocument(context.Background(), graphqlClient, matchingTarget(), DocumentCreateRequest{
		Title:   "created",
		Content: "body",
	})

	require.NoError(t, err)
	require.Equal(t, "document-id", document.ID)
	require.Equal(t, "created", document.Title)
}

func Test_CreateDocument_refuses_when_team_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	_, err := CreateDocument(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, DocumentCreateRequest{Title: "blocked"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateDocument_requires_title(t *testing.T) {
	_, err := CreateDocument(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), DocumentCreateRequest{},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateDocument_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	_, err := CreateDocument(
		context.Background(), graphqlClient, matchingTarget(), DocumentCreateRequest{Title: "x"},
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateDocument_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"DocumentCreate": `{"documentCreate":{"success":false,"document":` +
			documentJSON("x", "team-id", "LIT", "project-id") + `}}`,
	})

	_, err := CreateDocument(
		context.Background(), graphqlClient, matchingTarget(), DocumentCreateRequest{Title: "x"},
	)

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_UpdateDocument_returns_updated_document_when_target_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"document":       `{"document":` + documentJSON("current", "team-id", "LIT", "project-id") + `}`,
		"DocumentUpdate": `{"documentUpdate":{"success":true,"document":` + documentJSON("updated", "team-id", "LIT", "project-id") + `}}`,
	})

	document, err := UpdateDocument(context.Background(), graphqlClient, matchingTarget(), DocumentUpdateRequest{
		ID:    "document-id",
		Title: "updated",
	})

	require.NoError(t, err)
	require.Equal(t, "updated", document.Title)
}

func Test_UpdateDocument_refuses_when_team_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"document": `{"document":` + documentJSON("current", "other-team", "OTHER", "project-id") + `}`,
	})

	_, err := UpdateDocument(context.Background(), graphqlClient, matchingTarget(), DocumentUpdateRequest{
		ID:    "document-id",
		Title: "x",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateDocument_refuses_when_project_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"document": `{"document":` + documentJSON("current", "team-id", "LIT", "other-project") + `}`,
	})

	_, err := UpdateDocument(context.Background(), graphqlClient, matchingTarget(), DocumentUpdateRequest{
		ID:      "document-id",
		Content: "x",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateDocument_requires_id(t *testing.T) {
	_, err := UpdateDocument(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		DocumentUpdateRequest{Title: "x"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateDocument_requires_a_field(t *testing.T) {
	_, err := UpdateDocument(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		DocumentUpdateRequest{ID: "document-id"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateDocument_refuses_when_target_unresolved(t *testing.T) {
	_, err := UpdateDocument(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, DocumentUpdateRequest{ID: "document-id", Title: "x"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateDocument_wraps_document_read_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	_, err := UpdateDocument(
		context.Background(), graphqlClient, matchingTarget(), DocumentUpdateRequest{ID: "document-id", Title: "x"},
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateDocument_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"document": `{"document":` + documentJSON("current", "team-id", "LIT", "project-id") + `}`,
	})

	_, err := UpdateDocument(
		context.Background(), graphqlClient, matchingTarget(), DocumentUpdateRequest{ID: "document-id", Title: "x"},
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateDocument_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"document":       `{"document":` + documentJSON("current", "team-id", "LIT", "project-id") + `}`,
		"DocumentUpdate": `{"documentUpdate":{"success":false,"document":` + documentJSON("x", "team-id", "LIT", "project-id") + `}}`,
	})

	_, err := UpdateDocument(
		context.Background(), graphqlClient, matchingTarget(), DocumentUpdateRequest{ID: "document-id", Title: "x"},
	)

	require.ErrorIs(t, err, ErrMutationFailed)
}
