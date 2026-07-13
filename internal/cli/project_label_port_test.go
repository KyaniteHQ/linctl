package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeProjectLabelPort struct {
	created    client.ProjectLabelSummary
	createReq  client.ProjectLabelCreateRequest
	createErr  error
	updated    client.ProjectLabelSummary
	updateReq  client.ProjectLabelUpdateRequest
	updateErr  error
	retired    client.ProjectLabelSummary
	retireID   string
	retireOrg  bool
	retireErr  error
	restored   client.ProjectLabelSummary
	restoreID  string
	restoreOrg bool
	restoreErr error
}

func (port *fakeProjectLabelPort) CreateProjectLabel(
	_ context.Context, request client.ProjectLabelCreateRequest,
) (client.ProjectLabelSummary, error) {
	port.createReq = request

	return port.created, port.createErr
}

func (port *fakeProjectLabelPort) UpdateProjectLabel(
	_ context.Context, request client.ProjectLabelUpdateRequest,
) (client.ProjectLabelSummary, error) {
	port.updateReq = request

	return port.updated, port.updateErr
}

func (port *fakeProjectLabelPort) RetireProjectLabel(
	_ context.Context, id string, orgWide bool,
) (client.ProjectLabelSummary, error) {
	port.retireID = id
	port.retireOrg = orgWide

	return port.retired, port.retireErr
}

func (port *fakeProjectLabelPort) RestoreProjectLabel(
	_ context.Context, id string, orgWide bool,
) (client.ProjectLabelSummary, error) {
	port.restoreID = id
	port.restoreOrg = orgWide

	return port.restored, port.restoreErr
}

func Test_runProjectLabelCreate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectLabelPort{created: client.ProjectLabelSummary{
		ID: "project-label-id", Name: "Created project label", Color: "#f2c94c",
	}}

	err := runProjectLabelCreate(
		context.Background(), command, &rootOptions{}, port,
		client.ProjectLabelCreateRequest{Name: "Created project label", OrgWide: true},
	)

	require.NoError(t, err)
	require.Equal(t, "Created project label", port.createReq.Name)
	require.True(t, port.createReq.OrgWide)
	require.Contains(t, stdout.String(), "project-label-id Created project label #f2c94c")
}

func Test_runProjectLabelCreate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectLabelPort{createErr: errors.New("create failed")}

	err := runProjectLabelCreate(
		context.Background(), command, &rootOptions{}, port,
		client.ProjectLabelCreateRequest{Name: "Created project label"},
	)

	require.ErrorContains(t, err, "create failed")
}

func Test_runProjectLabelUpdate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectLabelPort{updated: client.ProjectLabelSummary{
		ID: "project-label-id", Name: "Updated project label", Color: "#f2c94c",
	}}

	err := runProjectLabelUpdate(
		context.Background(), command, &rootOptions{}, port,
		client.ProjectLabelUpdateRequest{ID: "project-label-id", Name: "Updated project label", OrgWide: true},
	)

	require.NoError(t, err)
	require.Equal(t, "project-label-id", port.updateReq.ID)
	require.Contains(t, stdout.String(), "project-label-id Updated project label #f2c94c")
}

func Test_runProjectLabelUpdate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectLabelPort{updateErr: errors.New("update failed")}

	err := runProjectLabelUpdate(
		context.Background(), command, &rootOptions{}, port,
		client.ProjectLabelUpdateRequest{ID: "project-label-id"},
	)

	require.ErrorContains(t, err, "update failed")
}

func Test_runProjectLabelRetire_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectLabelPort{retired: client.ProjectLabelSummary{
		ID: "project-label-id", Name: "Retired project label", Color: "#f2c94c",
	}}

	err := runProjectLabelRetire(context.Background(), command, &rootOptions{}, port, "project-label-id", true)

	require.NoError(t, err)
	require.Equal(t, "project-label-id", port.retireID)
	require.True(t, port.retireOrg)
	require.Contains(t, stdout.String(), "project-label-id Retired project label #f2c94c")
}

func Test_runProjectLabelRetire_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectLabelPort{retireErr: errors.New("retire failed")}

	err := runProjectLabelRetire(context.Background(), command, &rootOptions{}, port, "project-label-id", true)

	require.ErrorContains(t, err, "retire failed")
}

func Test_runProjectLabelRestore_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectLabelPort{restored: client.ProjectLabelSummary{
		ID: "project-label-id", Name: "Restored project label", Color: "#f2c94c",
	}}

	err := runProjectLabelRestore(context.Background(), command, &rootOptions{}, port, "project-label-id", true)

	require.NoError(t, err)
	require.Equal(t, "project-label-id", port.restoreID)
	require.True(t, port.restoreOrg)
	require.Contains(t, stdout.String(), "project-label-id Restored project label #f2c94c")
}

func Test_runProjectLabelRestore_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectLabelPort{restoreErr: errors.New("restore failed")}

	err := runProjectLabelRestore(context.Background(), command, &rootOptions{}, port, "project-label-id", true)

	require.ErrorContains(t, err, "restore failed")
}
