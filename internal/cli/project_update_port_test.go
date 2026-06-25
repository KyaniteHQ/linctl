package cli

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// fakeProjectUpdatePort is an in-memory Command Port: it captures the request
// and returns a domain summary, so the create command's body-resolution and
// health-normalization logic is tested without canned GraphQL JSON.
type fakeProjectUpdatePort struct {
	created     client.ProjectUpdateSummary
	createReq   client.ProjectUpdateCreateRequest
	createCalls int
	createErr   error
}

func (port *fakeProjectUpdatePort) CreateProjectUpdate(
	_ context.Context,
	request client.ProjectUpdateCreateRequest,
) (client.ProjectUpdateSummary, error) {
	port.createCalls++
	port.createReq = request

	return port.created, port.createErr
}

func Test_runProjectUpdateCreate_normalizes_health_through_the_port(t *testing.T) {
	command, stdout, stderr := bufferedCommand()
	port := &fakeProjectUpdatePort{
		created: client.ProjectUpdateSummary{ID: "pu-1", Health: "atRisk", DisplayName: "Omer", Body: "posted"},
	}

	err := runProjectUpdateCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectUpdateCreateRequest{ProjectID: "proj-1", Body: "raw body"},
		"at-risk", // alias -> normalized to canonical before reaching the port
		"",
	)

	require.NoError(t, err)
	require.Equal(t, 1, port.createCalls)
	require.Equal(t, "proj-1", port.createReq.ProjectID)
	require.Equal(t, "raw body", port.createReq.Body)
	require.Equal(t, "atRisk", port.createReq.Health)
	require.Contains(t, stderr.String(), "normalized")
	require.Contains(t, stdout.String(), "pu-1")
}

func Test_runProjectUpdateCreate_rejects_invalid_health_before_calling_the_port(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectUpdatePort{}

	err := runProjectUpdateCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectUpdateCreateRequest{ProjectID: "proj-1", Body: "body"},
		"bogus",
		"",
	)

	require.ErrorContains(t, err, "unknown health")
	require.Equal(t, 0, port.createCalls)
}

func Test_runProjectUpdateCreate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectUpdatePort{createErr: errors.New("create failed")}

	err := runProjectUpdateCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectUpdateCreateRequest{ProjectID: "proj-1", Body: "body"},
		"",
		"",
	)

	require.ErrorContains(t, err, "create failed")
	require.Equal(t, 1, port.createCalls)
}

func Test_runProjectUpdateCreate_reads_body_from_stdin(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	command.SetIn(strings.NewReader("body from stdin"))
	port := &fakeProjectUpdatePort{created: client.ProjectUpdateSummary{ID: "pu-1"}}

	err := runProjectUpdateCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectUpdateCreateRequest{ProjectID: "proj-1", Body: "-"},
		"",
		"",
	)

	require.NoError(t, err)
	require.Equal(t, 1, port.createCalls)
	require.Equal(t, "body from stdin", port.createReq.Body)
	require.Contains(t, stdout.String(), "pu-1")
}

func Test_runProjectUpdateCreate_reads_body_from_file(t *testing.T) {
	command, _, _ := bufferedCommand()
	path := writeTempTextFile(t, "body from file")
	port := &fakeProjectUpdatePort{created: client.ProjectUpdateSummary{ID: "pu-1"}}

	err := runProjectUpdateCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectUpdateCreateRequest{ProjectID: "proj-1"},
		"",
		path,
	)

	require.NoError(t, err)
	require.Equal(t, 1, port.createCalls)
	require.Equal(t, "body from file", port.createReq.Body)
}

func Test_runProjectUpdateCreate_rejects_body_and_body_file_before_reading_stdin(t *testing.T) {
	command, _, _ := bufferedCommand()
	// stdin errors if read; the --body/--body-file conflict check must fire
	// first, so the failing reader is never touched.
	command.SetIn(commandFailingReader{})
	path := writeTempTextFile(t, "body from file")
	port := &fakeProjectUpdatePort{}

	err := runProjectUpdateCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectUpdateCreateRequest{ProjectID: "proj-1", Body: "-"},
		"",
		path,
	)

	require.ErrorContains(t, err, "mutually exclusive")
	require.Equal(t, 0, port.createCalls)
}
