package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeCommentPort struct {
	updated   client.CommentSummary
	updateReq client.CommentUpdateRequest
	updateErr error
	deletedID string
	deleteID  string
	deleteErr error
}

func (port *fakeCommentPort) UpdateComment(
	_ context.Context,
	request client.CommentUpdateRequest,
) (client.CommentSummary, error) {
	port.updateReq = request

	return port.updated, port.updateErr
}

func (port *fakeCommentPort) DeleteComment(_ context.Context, commentID string) (string, error) {
	port.deleteID = commentID

	return port.deletedID, port.deleteErr
}

func Test_runCommentUpdate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeCommentPort{
		updated: client.CommentSummary{ID: "comment-id", DisplayName: "Omer", Body: "updated body"},
	}

	err := runCommentUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.CommentUpdateRequest{ID: "comment-id", Body: "updated body"},
		"",
	)

	require.NoError(t, err)
	require.Equal(t, "comment-id", port.updateReq.ID)
	require.Equal(t, "updated body", port.updateReq.Body)
	require.Contains(t, stdout.String(), "comment-id Omer updated body")
}

func Test_runCommentUpdate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeCommentPort{updateErr: errors.New("update failed")}

	err := runCommentUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.CommentUpdateRequest{ID: "comment-id", Body: "updated body"},
		"",
	)

	require.ErrorContains(t, err, "update failed")
}

func Test_runCommentDelete_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeCommentPort{deletedID: "comment-id"}

	err := runCommentDelete(context.Background(), command, &rootOptions{}, port, "comment-id")

	require.NoError(t, err)
	require.Equal(t, "comment-id", port.deleteID)
	require.Contains(t, stdout.String(), "comment-id deleted")
}

func Test_runCommentDelete_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeCommentPort{deleteErr: errors.New("delete failed")}

	err := runCommentDelete(context.Background(), command, &rootOptions{}, port, "comment-id")

	require.ErrorContains(t, err, "delete failed")
}

func Test_commentClientAdapter_forwards_to_client(t *testing.T) {
	adapter := commentAdapterFor(testCommandRuntime(commandFlowFakeClient{}))
	ctx := context.Background()

	comment, err := adapter.UpdateComment(ctx, client.CommentUpdateRequest{ID: "comment-id", Body: "updated body"})
	require.NoError(t, err)
	require.NotEmpty(t, comment.ID)

	deletedID, err := adapter.DeleteComment(ctx, "comment-id")
	require.NoError(t, err)
	require.NotEmpty(t, deletedID)
}
