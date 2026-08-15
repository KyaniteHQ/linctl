package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_ProjectAddTeamCommandFlow_reports_runtime_errors(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"project", "add-team", "project-id", "--to-team-id", "ops-team-id"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "runtime failed")
}

func Test_ProjectAddTeam_requires_destination(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"project", "add-team", "project-id"})

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "to-team")
}

func Test_ProjectAddTeam_writes_the_updated_project(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"project", "add-team", "project-id", "--to-team-id", "ops-team-id"})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "project-id")
	require.Empty(t, stderr.String())
}

func Test_IssueMoveTeamCommandFlow_reports_runtime_errors(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"issue", "move-team", "LIT-1", "--to-team-id", "ops-team-id"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "runtime failed")
}

func Test_IssueMoveTeam_requires_destination(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"issue", "move-team", "LIT-1"})

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "to-team")
}

func Test_IssueMoveTeam_writes_the_moved_issue(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"issue", "move-team", "LIT-1", "--to-team-id", "ops-team-id"})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "LIT-1")
	require.Empty(t, stderr.String())
}

func Test_IssueMoveProjectCommandFlow_reports_runtime_errors(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"issue", "move-project", "LIT-1", "--to-project-id", "eoir-project-id"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "runtime failed")
}

func Test_IssueMoveProject_requires_destination(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"issue", "move-project", "LIT-1"})

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "to-project-id")
}

func Test_IssueMoveProject_writes_the_moved_issue(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"issue", "move-project", "LIT-1", "--to-project-id", "eoir-project-id", "--format", "full"})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "LIT-1")
	require.Contains(t, stdout.String(), "project=EOIR Case Scraper")
	require.Empty(t, stderr.String())
}
