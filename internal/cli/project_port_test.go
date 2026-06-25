package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeProjectPort struct {
	created    client.ProjectSummary
	createReq  client.ProjectCreateRequest
	createErr  error
	updated    client.ProjectSummary
	updateReq  client.ProjectUpdateRequest
	updateErr  error
	archived   client.ProjectSummary
	archiveID  string
	archiveErr error
}

func (port *fakeProjectPort) CreateProject(
	_ context.Context,
	request client.ProjectCreateRequest,
) (client.ProjectSummary, error) {
	port.createReq = request

	return port.created, port.createErr
}

func (port *fakeProjectPort) UpdateProject(
	_ context.Context,
	request client.ProjectUpdateRequest,
) (client.ProjectSummary, error) {
	port.updateReq = request

	return port.updated, port.updateErr
}

func (port *fakeProjectPort) ArchiveProject(_ context.Context, projectID string) (client.ProjectSummary, error) {
	port.archiveID = projectID

	return port.archived, port.archiveErr
}

func Test_runProjectCreate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectPort{
		created: client.ProjectSummary{ID: "project-id", Name: "Created project", Status: client.ProjectStatus{Name: "Backlog"}},
	}

	err := runProjectCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectCreateRequest{Name: "Created project", Description: "description"},
	)

	require.NoError(t, err)
	require.Equal(t, "Created project", port.createReq.Name)
	require.Equal(t, "description", port.createReq.Description)
	require.Contains(t, stdout.String(), "project-id Created project [Backlog]")
}

func Test_runProjectCreate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectPort{createErr: errors.New("create failed")}

	err := runProjectCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectCreateRequest{Name: "Created project"},
	)

	require.ErrorContains(t, err, "create failed")
}

func Test_runProjectUpdate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectPort{
		updated: client.ProjectSummary{ID: "project-id", Name: "Updated project", Status: client.ProjectStatus{Name: "Started"}},
	}

	err := runProjectUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectUpdateRequest{ID: "project-id", Name: "Updated project"},
	)

	require.NoError(t, err)
	require.Equal(t, "project-id", port.updateReq.ID)
	require.Equal(t, "Updated project", port.updateReq.Name)
	require.Contains(t, stdout.String(), "project-id Updated project [Started]")
}

func Test_runProjectUpdate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectPort{updateErr: errors.New("update failed")}

	err := runProjectUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectUpdateRequest{ID: "project-id", Name: "Updated project"},
	)

	require.ErrorContains(t, err, "update failed")
}

func Test_runProjectArchive_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectPort{
		archived: client.ProjectSummary{ID: "project-id", Name: "Archived project", Status: client.ProjectStatus{Name: "Canceled"}},
	}

	err := runProjectArchive(context.Background(), command, &rootOptions{}, port, "project-id")

	require.NoError(t, err)
	require.Equal(t, "project-id", port.archiveID)
	require.Contains(t, stdout.String(), "project-id Archived project [Canceled]")
}

func Test_runProjectArchive_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectPort{archiveErr: errors.New("archive failed")}

	err := runProjectArchive(context.Background(), command, &rootOptions{}, port, "project-id")

	require.ErrorContains(t, err, "archive failed")
}

func Test_projectClientAdapter_forwards_to_client(t *testing.T) {
	adapter := commandAdapterFor(testCommandRuntime(commandFlowFakeClient{}))
	ctx := context.Background()

	created, err := adapter.CreateProject(ctx, client.ProjectCreateRequest{Name: "Created project"})
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	updated, err := adapter.UpdateProject(ctx, client.ProjectUpdateRequest{ID: "project-id", Name: "Updated project"})
	require.NoError(t, err)
	require.NotEmpty(t, updated.ID)

	archived, err := adapter.ArchiveProject(ctx, "project-id")
	require.NoError(t, err)
	require.NotEmpty(t, archived.ID)
}
