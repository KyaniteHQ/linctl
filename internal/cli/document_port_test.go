package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeDocumentPort struct {
	created   client.DocumentSummary
	createReq client.DocumentCreateRequest
	createErr error
	updated   client.DocumentSummary
	updateReq client.DocumentUpdateRequest
	updateErr error
}

func (port *fakeDocumentPort) CreateDocument(
	_ context.Context,
	request client.DocumentCreateRequest,
) (client.DocumentSummary, error) {
	port.createReq = request

	return port.created, port.createErr
}

func (port *fakeDocumentPort) UpdateDocument(
	_ context.Context,
	request client.DocumentUpdateRequest,
) (client.DocumentSummary, error) {
	port.updateReq = request

	return port.updated, port.updateErr
}

func Test_runDocumentCreate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeDocumentPort{
		created: client.DocumentSummary{ID: "document-id", Title: "Created doc", ParentType: "team"},
	}

	err := runDocumentCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.DocumentCreateRequest{Title: "Created doc", Content: "body"},
		"",
	)

	require.NoError(t, err)
	require.Equal(t, "Created doc", port.createReq.Title)
	require.Equal(t, "body", port.createReq.Content)
	require.Contains(t, stdout.String(), "document-id Created doc [team]")
}

func Test_runDocumentCreate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeDocumentPort{createErr: errors.New("create failed")}

	err := runDocumentCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.DocumentCreateRequest{Title: "Created doc"},
		"",
	)

	require.ErrorContains(t, err, "create failed")
}

func Test_runDocumentUpdate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeDocumentPort{
		updated: client.DocumentSummary{ID: "document-id", Title: "Updated doc", ParentType: "team"},
	}

	err := runDocumentUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.DocumentUpdateRequest{ID: "document-id", Title: "Updated doc", Content: "body"},
		"",
	)

	require.NoError(t, err)
	require.Equal(t, "document-id", port.updateReq.ID)
	require.Equal(t, "Updated doc", port.updateReq.Title)
	require.Equal(t, "body", port.updateReq.Content)
	require.Contains(t, stdout.String(), "document-id Updated doc [team]")
}

func Test_runDocumentUpdate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeDocumentPort{updateErr: errors.New("update failed")}

	err := runDocumentUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.DocumentUpdateRequest{ID: "document-id", Title: "Updated doc"},
		"",
	)

	require.ErrorContains(t, err, "update failed")
}

func Test_documentClientAdapter_forwards_to_client(t *testing.T) {
	adapter := documentAdapterFor(testCommandRuntime(commandFlowFakeClient{}))
	ctx := context.Background()

	created, err := adapter.CreateDocument(ctx, client.DocumentCreateRequest{Title: "Created doc"})
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	updated, err := adapter.UpdateDocument(ctx, client.DocumentUpdateRequest{ID: "document-id", Title: "Updated doc"})
	require.NoError(t, err)
	require.NotEmpty(t, updated.ID)
}
