package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeIssueLabelPort struct {
	added     client.IssueSummary
	addReq    client.IssueLabelAssociationRequest
	addErr    error
	removed   client.IssueSummary
	removeReq client.IssueLabelAssociationRequest
	removeErr error
}

func (port *fakeIssueLabelPort) AddIssueLabel(
	_ context.Context,
	request client.IssueLabelAssociationRequest,
) (client.IssueSummary, error) {
	port.addReq = request

	return port.added, port.addErr
}

func (port *fakeIssueLabelPort) RemoveIssueLabel(
	_ context.Context,
	request client.IssueLabelAssociationRequest,
) (client.IssueSummary, error) {
	port.removeReq = request

	return port.removed, port.removeErr
}

func Test_runIssueAddLabel_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeIssueLabelPort{added: client.IssueSummary{Identifier: "LIT-1", Title: "Labeled issue", StateType: "unstarted", State: "Todo"}}

	err := runIssueAddLabel(context.Background(), command, &rootOptions{}, port, client.IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-1", port.addReq.IssueID)
	require.Equal(t, "label-id", port.addReq.LabelID)
	require.Contains(t, stdout.String(), "LIT-1 Labeled issue [Todo]")
}

func Test_runIssueAddLabel_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeIssueLabelPort{addErr: errors.New("add label failed")}

	err := runIssueAddLabel(context.Background(), command, &rootOptions{}, port, client.IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.ErrorContains(t, err, "add label failed")
}

func Test_runIssueRemoveLabel_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeIssueLabelPort{removed: client.IssueSummary{Identifier: "LIT-1", Title: "Unlabeled issue", StateType: "unstarted", State: "Todo"}}

	err := runIssueRemoveLabel(context.Background(), command, &rootOptions{}, port, client.IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-1", port.removeReq.IssueID)
	require.Equal(t, "label-id", port.removeReq.LabelID)
	require.Contains(t, stdout.String(), "LIT-1 Unlabeled issue [Todo]")
}

func Test_runIssueRemoveLabel_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeIssueLabelPort{removeErr: errors.New("remove label failed")}

	err := runIssueRemoveLabel(context.Background(), command, &rootOptions{}, port, client.IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.ErrorContains(t, err, "remove label failed")
}
