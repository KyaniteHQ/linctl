package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeLabelPort struct {
	created    client.LabelSummary
	createReq  client.LabelCreateRequest
	createErr  error
	updated    client.LabelSummary
	updateReq  client.LabelUpdateRequest
	updateErr  error
	retired    client.LabelSummary
	retireID   string
	retireOrg  bool
	retireErr  error
	restored   client.LabelSummary
	restoreID  string
	restoreOrg bool
	restoreErr error
}

func (port *fakeLabelPort) CreateLabel(_ context.Context, request client.LabelCreateRequest) (client.LabelSummary, error) {
	port.createReq = request

	return port.created, port.createErr
}

func (port *fakeLabelPort) UpdateLabel(_ context.Context, request client.LabelUpdateRequest) (client.LabelSummary, error) {
	port.updateReq = request

	return port.updated, port.updateErr
}

func (port *fakeLabelPort) RetireLabel(_ context.Context, id string, orgWide bool) (client.LabelSummary, error) {
	port.retireID = id
	port.retireOrg = orgWide

	return port.retired, port.retireErr
}

func (port *fakeLabelPort) RestoreLabel(_ context.Context, id string, orgWide bool) (client.LabelSummary, error) {
	port.restoreID = id
	port.restoreOrg = orgWide

	return port.restored, port.restoreErr
}

func Test_runLabelCreate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeLabelPort{created: client.LabelSummary{ID: "label-id", Name: "Created label", Color: "#ff0000"}}

	err := runLabelCreate(context.Background(), command, &rootOptions{}, port, client.LabelCreateRequest{Name: "Created label"})

	require.NoError(t, err)
	require.Equal(t, "Created label", port.createReq.Name)
	require.Contains(t, stdout.String(), "label-id Created label #ff0000")
}

func Test_runLabelCreate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeLabelPort{createErr: errors.New("create failed")}

	err := runLabelCreate(context.Background(), command, &rootOptions{}, port, client.LabelCreateRequest{Name: "Created label"})

	require.ErrorContains(t, err, "create failed")
}

func Test_runLabelUpdate_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeLabelPort{updated: client.LabelSummary{ID: "label-id", Name: "Updated label", Color: "#ff0000"}}

	err := runLabelUpdate(context.Background(), command, &rootOptions{}, port, client.LabelUpdateRequest{ID: "label-id", Name: "Updated label"})

	require.NoError(t, err)
	require.Equal(t, "label-id", port.updateReq.ID)
	require.Contains(t, stdout.String(), "label-id Updated label #ff0000")
}

func Test_runLabelUpdate_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeLabelPort{updateErr: errors.New("update failed")}

	err := runLabelUpdate(context.Background(), command, &rootOptions{}, port, client.LabelUpdateRequest{ID: "label-id"})

	require.ErrorContains(t, err, "update failed")
}

func Test_runLabelRetire_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeLabelPort{retired: client.LabelSummary{ID: "label-id", Name: "Retired label", Color: "#ff0000"}}

	err := runLabelRetire(context.Background(), command, &rootOptions{}, port, "label-id", true)

	require.NoError(t, err)
	require.Equal(t, "label-id", port.retireID)
	require.True(t, port.retireOrg)
	require.Contains(t, stdout.String(), "label-id Retired label #ff0000")
}

func Test_runLabelRetire_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeLabelPort{retireErr: errors.New("retire failed")}

	err := runLabelRetire(context.Background(), command, &rootOptions{}, port, "label-id", false)

	require.ErrorContains(t, err, "retire failed")
}

func Test_runLabelRestore_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeLabelPort{restored: client.LabelSummary{ID: "label-id", Name: "Restored label", Color: "#ff0000"}}

	err := runLabelRestore(context.Background(), command, &rootOptions{}, port, "label-id", true)

	require.NoError(t, err)
	require.Equal(t, "label-id", port.restoreID)
	require.True(t, port.restoreOrg)
	require.Contains(t, stdout.String(), "label-id Restored label #ff0000")
}

func Test_runLabelRestore_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeLabelPort{restoreErr: errors.New("restore failed")}

	err := runLabelRestore(context.Background(), command, &rootOptions{}, port, "label-id", false)

	require.ErrorContains(t, err, "restore failed")
}
