package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeCyclePort struct {
	created    client.CycleSummary
	createReq  client.CycleCreateRequest
	createErr  error
	updated    client.CycleSummary
	updateReq  client.CycleUpdateRequest
	updateErr  error
	archived   client.CycleSummary
	archiveID  string
	archiveErr error
}

func (port *fakeCyclePort) CreateCycle(
	_ context.Context,
	request client.CycleCreateRequest,
) (client.CycleSummary, error) {
	port.createReq = request

	return port.created, port.createErr
}

func (port *fakeCyclePort) UpdateCycle(
	_ context.Context,
	request client.CycleUpdateRequest,
) (client.CycleSummary, error) {
	port.updateReq = request

	return port.updated, port.updateErr
}

func (port *fakeCyclePort) ArchiveCycle(_ context.Context, cycleID string) (client.CycleSummary, error) {
	port.archiveID = cycleID

	return port.archived, port.archiveErr
}

func Test_runCycleCreate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeCyclePort{
		created: client.CycleSummary{ID: "cycle-id", Name: "Created cycle", Status: "active"},
	}

	err := runCycleCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.CycleCreateRequest{
			Name:        "Created cycle",
			Description: "description",
			StartsAt:    "2026-06-01T00:00:00Z",
			EndsAt:      "2026-06-15T00:00:00Z",
		},
	)

	require.NoError(t, err)
	require.Equal(t, "Created cycle", port.createReq.Name)
	require.Equal(t, "description", port.createReq.Description)
	require.Equal(t, "2026-06-01T00:00:00Z", port.createReq.StartsAt)
	require.Equal(t, "2026-06-15T00:00:00Z", port.createReq.EndsAt)
	require.Contains(t, stdout.String(), "cycle-id Created cycle [active]")
}

func Test_runCycleCreate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeCyclePort{createErr: errors.New("create failed")}

	err := runCycleCreate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.CycleCreateRequest{Name: "Created cycle"},
	)

	require.ErrorContains(t, err, "create failed")
}

func Test_runCycleUpdate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeCyclePort{
		updated: client.CycleSummary{ID: "cycle-id", Name: "Updated cycle", Status: "active"},
	}

	err := runCycleUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.CycleUpdateRequest{ID: "cycle-id", Name: "Updated cycle"},
	)

	require.NoError(t, err)
	require.Equal(t, "cycle-id", port.updateReq.ID)
	require.Equal(t, "Updated cycle", port.updateReq.Name)
	require.Contains(t, stdout.String(), "cycle-id Updated cycle [active]")
}

func Test_runCycleUpdate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeCyclePort{updateErr: errors.New("update failed")}

	err := runCycleUpdate(
		context.Background(),
		command,
		&rootOptions{},
		port,
		client.CycleUpdateRequest{ID: "cycle-id", Name: "Updated cycle"},
	)

	require.ErrorContains(t, err, "update failed")
}

func Test_runCycleArchive_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeCyclePort{
		archived: client.CycleSummary{ID: "cycle-id", Name: "Archived cycle", Status: "active"},
	}

	err := runCycleArchive(context.Background(), command, &rootOptions{}, port, "cycle-id")

	require.NoError(t, err)
	require.Equal(t, "cycle-id", port.archiveID)
	require.Contains(t, stdout.String(), "cycle-id Archived cycle [active]")
}

func Test_runCycleArchive_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeCyclePort{archiveErr: errors.New("archive failed")}

	err := runCycleArchive(context.Background(), command, &rootOptions{}, port, "cycle-id")

	require.ErrorContains(t, err, "archive failed")
}

func Test_cycleClientAdapter_forwards_to_client(t *testing.T) {
	adapter := commandAdapterFor(testCommandRuntime(cycleCommandFlowFakeClient{}))
	ctx := context.Background()

	created, err := adapter.CreateCycle(ctx, client.CycleCreateRequest{
		Name:     "Created cycle",
		StartsAt: "2026-06-01T00:00:00Z",
		EndsAt:   "2026-06-15T00:00:00Z",
	})
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	updated, err := adapter.UpdateCycle(ctx, client.CycleUpdateRequest{ID: "cycle-id", Name: "Updated cycle"})
	require.NoError(t, err)
	require.NotEmpty(t, updated.ID)

	archived, err := adapter.ArchiveCycle(ctx, "cycle-id")
	require.NoError(t, err)
	require.NotEmpty(t, archived.ID)
}
