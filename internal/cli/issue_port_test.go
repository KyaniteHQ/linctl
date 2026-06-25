package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// fakeIssuePort is an in-memory Command Port: it returns domain summaries
// directly, so issue command logic is tested without canned GraphQL JSON.
type fakeIssuePort struct {
	created     client.IssueSummary
	createReq   client.IssueCreateRequest
	createCalls int
	createErr   error
	closed      client.IssueSummary
	closeID     string
	closeCalls  int
	template    client.IssueTemplateContent
	templateErr error
	updated     client.IssueSummary
	updateReq   client.IssueUpdateRequest
	updateCalls int
	updateErr   error
	commented   client.IssueCommentResult
	commentReq  client.IssueCommentRequest
	commentErr  error
}

func (port *fakeIssuePort) UpdateIssue(
	_ context.Context,
	request client.IssueUpdateRequest,
) (client.IssueSummary, error) {
	port.updateCalls++
	port.updateReq = request

	return port.updated, port.updateErr
}

func (port *fakeIssuePort) CommentOnIssue(
	_ context.Context,
	request client.IssueCommentRequest,
) (client.IssueCommentResult, error) {
	port.commentReq = request

	return port.commented, port.commentErr
}

func (port *fakeIssuePort) CreateIssue(
	_ context.Context,
	request client.IssueCreateRequest,
) (client.IssueSummary, error) {
	port.createCalls++
	port.createReq = request

	return port.created, port.createErr
}

func (port *fakeIssuePort) CloseIssue(_ context.Context, issueID string) (client.IssueSummary, error) {
	port.closeCalls++
	port.closeID = issueID

	return port.closed, nil
}

func (port *fakeIssuePort) GetIssueTemplateContent(
	_ context.Context,
	_ string,
) (client.IssueTemplateContent, error) {
	return port.template, port.templateErr
}

func bufferedCommand() (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	command := &cobra.Command{}
	var stdout, stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)

	return command, &stdout, &stderr
}

func Test_runIssueCreate_assembles_request_through_the_port(t *testing.T) {
	command, stdout, stderr := bufferedCommand()
	port := &fakeIssuePort{
		created:  client.IssueSummary{ID: "iss-id", Identifier: "LIT-9", Title: "Created", State: "Todo"},
		template: client.IssueTemplateContent{Title: "Template title", Description: "Template body"},
	}
	estimate := 5

	err := runIssueCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.IssueCreateRequest{}, // empty title/description -> template fills them
		issueCreateFlags{templateID: "tmpl-1", state: "in progress", priority: "high"},
		&estimate,
	)

	require.NoError(t, err)
	require.Equal(t, 1, port.createCalls)
	// template defaults filled the empty fields
	require.Equal(t, "Template title", port.createReq.Title)
	require.Equal(t, "Template body", port.createReq.Description)
	// alias normalization reached the port as canonical values
	require.Equal(t, "started", port.createReq.StateType)
	require.Equal(t, "2", port.createReq.Priority)
	// estimate gating forwarded the resolved pointer
	require.NotNil(t, port.createReq.Estimate)
	require.Equal(t, 5, *port.createReq.Estimate)
	// normalization emitted a stderr note, and the issue rendered to stdout
	require.Contains(t, stderr.String(), "normalized")
	require.Contains(t, stdout.String(), "LIT-9")
}

func Test_runIssueCreate_dry_run_never_calls_the_port(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeIssuePort{}

	err := runIssueCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.IssueCreateRequest{Title: "Draft only"},
		issueCreateFlags{dryRun: true},
		nil,
	)

	require.NoError(t, err)
	require.Equal(t, 0, port.createCalls)
	require.Contains(t, stdout.String(), "Draft only")
}

func Test_runIssueCreate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeIssuePort{createErr: errors.New("create failed")}

	err := runIssueCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.IssueCreateRequest{Title: "X"},
		issueCreateFlags{},
		nil,
	)

	require.ErrorContains(t, err, "create failed")
	require.Equal(t, 1, port.createCalls)
}

func Test_runIssueCreate_propagates_template_error_before_creating(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeIssuePort{templateErr: errors.New("no such template")}

	err := runIssueCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.IssueCreateRequest{Title: "X"},
		issueCreateFlags{templateID: "missing"},
		nil,
	)

	require.ErrorContains(t, err, "no such template")
	require.Equal(t, 0, port.createCalls)
}

func Test_runIssueClose_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeIssuePort{
		closed: client.IssueSummary{ID: "c-id", Identifier: "LIT-3", Title: "Done", State: "Done"},
	}

	err := runIssueClose(context.Background(), command, &rootOptions{}, port, "LIT-3")

	require.NoError(t, err)
	require.Equal(t, 1, port.closeCalls)
	require.Equal(t, "LIT-3", port.closeID)
	require.Contains(t, stdout.String(), "LIT-3")
}

func Test_runIssueUpdate_assembles_request_through_the_port(t *testing.T) {
	command, stdout, stderr := bufferedCommand()
	port := &fakeIssuePort{
		updated: client.IssueSummary{ID: "u-id", Identifier: "LIT-1", Title: "Renamed", State: "Done"},
	}
	estimate := 8

	err := runIssueUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.IssueUpdateRequest{ID: "LIT-1", Title: "Renamed"},
		issueUpdateFlags{state: "done", priority: "urgent"},
		&estimate,
	)

	require.NoError(t, err)
	require.Equal(t, 1, port.updateCalls)
	require.Equal(t, "LIT-1", port.updateReq.ID)
	require.Equal(t, "completed", port.updateReq.StateType)
	require.Equal(t, "1", port.updateReq.Priority)
	require.NotNil(t, port.updateReq.Estimate)
	require.Equal(t, 8, *port.updateReq.Estimate)
	require.Contains(t, stderr.String(), "normalized")
	require.Contains(t, stdout.String(), "LIT-1")
}

func Test_runIssueUpdate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeIssuePort{updateErr: errors.New("update failed")}

	err := runIssueUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.IssueUpdateRequest{ID: "LIT-1"},
		issueUpdateFlags{},
		nil,
	)

	require.ErrorContains(t, err, "update failed")
	require.Equal(t, 1, port.updateCalls)
}

func Test_runIssueBodyWriteCommand_comments_through_the_port(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeIssuePort{
		commented: client.IssueCommentResult{ID: "cmt-1", Issue: client.IssueSummary{Identifier: "LIT-1"}},
	}

	err := runIssueBodyWriteCommand(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.IssueCommentRequest{ID: "LIT-1", Body: "looks good"},
		"",
	)

	require.NoError(t, err)
	require.Equal(t, "looks good", port.commentReq.Body)
	require.Equal(t, "LIT-1", port.commentReq.ID)
	require.Contains(t, stdout.String(), "comment cmt-1 on LIT-1")
}

func Test_runIssueBodyWriteCommand_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeIssuePort{commentErr: errors.New("comment failed")}

	err := runIssueBodyWriteCommand(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.IssueCommentRequest{ID: "LIT-1", Body: "x"},
		"",
	)

	require.ErrorContains(t, err, "comment failed")
}

// Test_issueClientAdapter_forwards_to_client proves the production adapter wires
// each port method to the right client free function. The canned GraphQL JSON is
// confined here, at the adapter seam, rather than smeared across command tests.
func Test_issueClientAdapter_forwards_to_client(t *testing.T) {
	adapter := issueAdapterFor(testCommandRuntime(commandFlowFakeClient{}))
	ctx := context.Background()

	created, err := adapter.CreateIssue(ctx, client.IssueCreateRequest{Title: "Created issue"})
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	closed, err := adapter.CloseIssue(ctx, "LIT-1")
	require.NoError(t, err)
	require.NotEmpty(t, closed.ID)

	_, err = adapter.GetIssueTemplateContent(ctx, "tmpl-1")
	require.NoError(t, err)

	updated, err := adapter.UpdateIssue(ctx, client.IssueUpdateRequest{ID: "LIT-1", Title: "Renamed"})
	require.NoError(t, err)
	require.NotEmpty(t, updated.ID)

	comment, err := adapter.CommentOnIssue(ctx, client.IssueCommentRequest{ID: "LIT-1", Body: "note"})
	require.NoError(t, err)
	require.NotEmpty(t, comment.ID)
}
