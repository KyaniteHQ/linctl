package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func documentJSON(title string, projectID string) string {
	project := `null`
	if projectID != "" {
		project = `{"id":"` + projectID + `","name":"fixture"}`
	}

	return `{
		"id":"document-id",
		"title":"` + title + `",
		"slugId":"doc-slug",
		"archivedAt":null,
		"project":` + project + `,
		"issue":null
	}`
}

func Test_CreateDocument_returns_created_document_when_target_matches(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"DocumentCreate": `{"documentCreate":{"success":true,"document":` +
			documentJSON("created", "project-id") + `}}`,
	})}

	document, err := CreateDocument(context.Background(), recorder, matchingTarget(), DocumentCreateRequest{
		Title:   "created",
		Content: "body",
	})

	require.NoError(t, err)
	require.Equal(t, "document-id", document.ID)
	require.Equal(t, "created", document.Title)
	require.JSONEq(t, `{
		"input": {
			"title": "created",
			"content": "body",
			"projectId": "project-id"
		}
	}`, string(recorder.variablesFor(t, "DocumentCreate")))
}

func Test_CreateDocument_refuses_when_team_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	_, err := CreateDocument(context.Background(), graphqlClient, config.Target{
		OrgID:     "org-id",
		TeamKey:   "WRONG",
		TeamID:    "wrong-id",
		ProjectID: "project-id",
	}, DocumentCreateRequest{Title: "blocked"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_DocumentWrites_require_a_pinned_project_before_linear_access(t *testing.T) {
	writes := map[string]func(graphqlClient *recordingGraphQLClient) error{
		"create": func(graphqlClient *recordingGraphQLClient) error {
			_, err := CreateDocument(context.Background(), graphqlClient, teamOnlyTarget(), DocumentCreateRequest{Title: "x"})
			return err
		},
		"update": func(graphqlClient *recordingGraphQLClient) error {
			_, err := UpdateDocument(
				context.Background(), graphqlClient, teamOnlyTarget(), DocumentUpdateRequest{ID: "document-id", Title: "x"},
			)
			return err
		},
	}
	for name, write := range writes {
		t.Run(name, func(t *testing.T) {
			recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{})}

			err := write(recorder)

			require.ErrorIs(t, err, ErrWriteInvalid)
			require.ErrorContains(t, err, "project_id")
			require.Empty(t, recorder.requests)
		})
	}
}

func Test_CreateDocument_requires_title(t *testing.T) {
	_, err := CreateDocument(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), DocumentCreateRequest{},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateDocument_wraps_mutation_error(t *testing.T) {
	sentinel := errors.New("linear unavailable")
	graphqlClient := issueWriteFakeClient(map[string]string{}).withError(sentinel)

	_, err := CreateDocument(
		context.Background(), graphqlClient, matchingTarget(), DocumentCreateRequest{Title: "x"},
	)

	require.ErrorIs(t, err, sentinel)
	require.ErrorContains(t, err, "create document")
}

func Test_CreateDocument_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"DocumentCreate": `{"documentCreate":{"success":false,"document":` +
			documentJSON("x", "project-id") + `}}`,
	})

	_, err := CreateDocument(
		context.Background(), graphqlClient, matchingTarget(), DocumentCreateRequest{Title: "x"},
	)

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_UpdateDocument_checks_and_updates_the_requested_document(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"document":       `{"document":` + documentJSON("current", "project-id") + `}`,
		"DocumentUpdate": `{"documentUpdate":{"success":true,"document":` + documentJSON("updated", "project-id") + `}}`,
	})}

	document, err := UpdateDocument(context.Background(), recorder, matchingTarget(), DocumentUpdateRequest{
		ID:    "requested-document-id",
		Title: "updated",
	})

	require.NoError(t, err)
	require.Equal(t, "updated", document.Title)
	require.JSONEq(t, `{"id":"requested-document-id"}`, string(recorder.variablesFor(t, "document")))
	require.JSONEq(t, `{"id":"requested-document-id","input":{"title":"updated"}}`,
		string(recorder.variablesFor(t, "DocumentUpdate")))
}

func Test_UpdateDocument_refuses_a_document_outside_any_project(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"document": `{"document":` + documentJSON("current", "") + `}`,
	})}

	_, err := UpdateDocument(context.Background(), recorder, matchingTarget(), DocumentUpdateRequest{
		ID:    "document-id",
		Title: "x",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.Zero(t, recorder.countOf("DocumentUpdate"))
}

func Test_UpdateDocument_refuses_when_project_differs(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"document": `{"document":` + documentJSON("current", "other-project") + `}`,
	})}

	_, err := UpdateDocument(context.Background(), recorder, matchingTarget(), DocumentUpdateRequest{
		ID:      "document-id",
		Content: "x",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.Zero(t, recorder.countOf("DocumentUpdate"))
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
		OrgID:     "org-id",
		TeamKey:   "WRONG",
		TeamID:    "wrong-id",
		ProjectID: "project-id",
	}, DocumentUpdateRequest{ID: "document-id", Title: "x"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateDocument_wraps_document_read_error(t *testing.T) {
	sentinel := errors.New("linear unavailable")
	graphqlClient := issueWriteFakeClient(map[string]string{}).withError(sentinel)

	_, err := UpdateDocument(
		context.Background(), graphqlClient, matchingTarget(),
		DocumentUpdateRequest{ID: "requested-document-id", Title: "x"},
	)

	require.ErrorIs(t, err, sentinel)
	require.ErrorContains(t, err, "get document requested-document-id")
}

func Test_UpdateDocument_wraps_mutation_error(t *testing.T) {
	sentinel := errors.New("linear unavailable")
	graphqlClient := issueWriteFakeClient(map[string]string{
		"document": `{"document":` + documentJSON("current", "project-id") + `}`,
	}).withError(sentinel)

	_, err := UpdateDocument(
		context.Background(), graphqlClient, matchingTarget(),
		DocumentUpdateRequest{ID: "requested-document-id", Title: "x"},
	)

	require.ErrorIs(t, err, sentinel)
	require.ErrorContains(t, err, "update document requested-document-id")
}

func Test_UpdateDocument_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"document":       `{"document":` + documentJSON("current", "project-id") + `}`,
		"DocumentUpdate": `{"documentUpdate":{"success":false,"document":` + documentJSON("x", "project-id") + `}}`,
	})

	_, err := UpdateDocument(
		context.Background(), graphqlClient, matchingTarget(), DocumentUpdateRequest{ID: "document-id", Title: "x"},
	)

	require.ErrorIs(t, err, ErrMutationFailed)
}
