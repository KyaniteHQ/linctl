package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeProjectLabelAssociationPort struct {
	added     client.ProjectSummary
	addReq    client.ProjectLabelAssociationRequest
	addErr    error
	removed   client.ProjectSummary
	removeReq client.ProjectLabelAssociationRequest
	removeErr error
}

func (port *fakeProjectLabelAssociationPort) AddProjectLabel(
	_ context.Context,
	request client.ProjectLabelAssociationRequest,
) (client.ProjectSummary, error) {
	port.addReq = request

	return port.added, port.addErr
}

func (port *fakeProjectLabelAssociationPort) RemoveProjectLabel(
	_ context.Context,
	request client.ProjectLabelAssociationRequest,
) (client.ProjectSummary, error) {
	port.removeReq = request

	return port.removed, port.removeErr
}

func Test_runProjectAddLabel_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectLabelAssociationPort{
		added: client.ProjectSummary{ID: "project-id", Name: "Labeled project", Status: client.ProjectStatus{Name: "Backlog"}},
	}

	err := runProjectAddLabel(context.Background(), command, &rootOptions{}, port, client.ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.NoError(t, err)
	require.Equal(t, "project-id", port.addReq.ProjectID)
	require.Equal(t, "label-id", port.addReq.LabelID)
	require.Contains(t, stdout.String(), "project-id Labeled project [Backlog]")
}

func Test_runProjectAddLabel_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectLabelAssociationPort{addErr: errors.New("add label failed")}

	err := runProjectAddLabel(context.Background(), command, &rootOptions{}, port, client.ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.ErrorContains(t, err, "add label failed")
}

func Test_runProjectRemoveLabel_calls_the_port_and_renders(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeProjectLabelAssociationPort{
		removed: client.ProjectSummary{ID: "project-id", Name: "Unlabeled project", Status: client.ProjectStatus{Name: "Backlog"}},
	}

	err := runProjectRemoveLabel(context.Background(), command, &rootOptions{}, port, client.ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.NoError(t, err)
	require.Equal(t, "project-id", port.removeReq.ProjectID)
	require.Equal(t, "label-id", port.removeReq.LabelID)
	require.Contains(t, stdout.String(), "project-id Unlabeled project [Backlog]")
}

func Test_runProjectRemoveLabel_propagates_port_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	port := &fakeProjectLabelAssociationPort{removeErr: errors.New("remove label failed")}

	err := runProjectRemoveLabel(context.Background(), command, &rootOptions{}, port, client.ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.ErrorContains(t, err, "remove label failed")
}
