package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeProjectMilestonePort struct {
	created   client.ProjectMilestoneSummary
	createReq client.ProjectMilestoneCreateRequest
	createErr error
	updated   client.ProjectMilestoneSummary
	updateReq client.ProjectMilestoneUpdateRequest
	updateErr error
}

func (port *fakeProjectMilestonePort) CreateProjectMilestone(
	_ context.Context,
	request client.ProjectMilestoneCreateRequest,
) (client.ProjectMilestoneSummary, error) {
	port.createReq = request

	return port.created, port.createErr
}

func (port *fakeProjectMilestonePort) UpdateProjectMilestone(
	_ context.Context,
	request client.ProjectMilestoneUpdateRequest,
) (client.ProjectMilestoneSummary, error) {
	port.updateReq = request

	return port.updated, port.updateErr
}

func Test_runProjectMilestoneCreate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectMilestonePort{
		created: client.ProjectMilestoneSummary{ID: "project-milestone-id", Name: "Created milestone", Status: "next"},
	}

	err := runProjectMilestoneCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectMilestoneCreateRequest{
			ProjectID:   "project-id",
			Name:        "Created milestone",
			Description: "description",
			TargetDate:  "2026-07-01",
		},
	)

	require.NoError(t, err)
	require.Equal(t, "project-id", port.createReq.ProjectID)
	require.Equal(t, "Created milestone", port.createReq.Name)
	require.Equal(t, "description", port.createReq.Description)
	require.Equal(t, "2026-07-01", port.createReq.TargetDate)
	require.Contains(t, stdout.String(), "project-milestone-id Created milestone [next]")
}

func Test_runProjectMilestoneCreate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectMilestonePort{createErr: errors.New("create failed")}

	err := runProjectMilestoneCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "Created milestone"},
	)

	require.ErrorContains(t, err, "create failed")
}

func Test_runProjectMilestoneUpdate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectMilestonePort{
		updated: client.ProjectMilestoneSummary{ID: "project-milestone-id", Name: "Updated milestone", Status: "done"},
	}

	err := runProjectMilestoneUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectMilestoneUpdateRequest{
			ID:          "project-milestone-id",
			Name:        "Updated milestone",
			Description: "description",
			TargetDate:  "2026-07-15",
		},
	)

	require.NoError(t, err)
	require.Equal(t, "project-milestone-id", port.updateReq.ID)
	require.Equal(t, "Updated milestone", port.updateReq.Name)
	require.Equal(t, "description", port.updateReq.Description)
	require.Equal(t, "2026-07-15", port.updateReq.TargetDate)
	require.Contains(t, stdout.String(), "project-milestone-id Updated milestone [done]")
}

func Test_runProjectMilestoneUpdate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectMilestonePort{updateErr: errors.New("update failed")}

	err := runProjectMilestoneUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.ProjectMilestoneUpdateRequest{ID: "project-milestone-id", Name: "Updated milestone"},
	)

	require.ErrorContains(t, err, "update failed")
}

func Test_projectMilestoneClientAdapter_forwards_to_client(t *testing.T) {
	adapter := commandAdapterFor(testCommandRuntime(commandFlowFakeClient{}))
	ctx := context.Background()

	created, err := adapter.CreateProjectMilestone(ctx, client.ProjectMilestoneCreateRequest{
		ProjectID: "project-id",
		Name:      "Created milestone",
	})
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	updated, err := adapter.UpdateProjectMilestone(ctx, client.ProjectMilestoneUpdateRequest{
		ID:   "project-milestone-id",
		Name: "Updated milestone",
	})
	require.NoError(t, err)
	require.NotEmpty(t, updated.ID)
}
