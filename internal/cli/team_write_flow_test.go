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

// Test_TeamCreateCommandFlow_reports_runtime_errors covers the
// buildCommandRuntime error branch inside the command's RunE closure.
func Test_TeamCreateCommandFlow_reports_runtime_errors(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"team", "create", "--name", "Operations", "--org-wide"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "runtime failed")
}

// Test_TeamCreateCommandFlow_reports_writer_errors covers the render step: the
// guarded write succeeds but stdout fails, so the error must propagate.
func Test_TeamCreateCommandFlow_reports_writer_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"team", "create", "--name", "Operations", "--org-wide"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "write failed")
}

// Test_TeamCreate_requires_org_wide proves the end-to-end refusal: a Team is
// what a pin names, so a create cannot land inside the Pinned Target's team and
// is a hard stop when --org-wide is omitted, with the flag named on stderr and
// nothing written to stdout.
func Test_TeamCreate_requires_org_wide(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"team", "create", "--name", "Operations"})

	require.ErrorIs(t, err, client.ErrTargetMismatch)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "org-wide")
}

// Test_TeamCreate_writes_the_created_team is the allow arm the refusal above
// needs: without it a command that refused everything would pass just as well.
func Test_TeamCreate_writes_the_created_team(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"team", "create", "--name", "Operations", "--key", "OPS", "--org-wide"})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "team-id")
	require.Empty(t, stderr.String())
}

func Test_TeamDeleteCommandFlow_reports_runtime_errors(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"team", "delete", "ops-team-id", "--org-wide"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "runtime failed")
}

func Test_TeamDeleteCommandFlow_reports_writer_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"team", "delete", "ops-team-id", "--org-wide"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "write failed")
}

// Test_TeamDelete_requires_org_wide proves the end-to-end refusal with the
// flag named on stderr and nothing written to stdout.
func Test_TeamDelete_requires_org_wide(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"team", "delete", "ops-team-id"})

	require.ErrorIs(t, err, client.ErrTargetMismatch)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "org-wide")
}

// Test_TeamDelete_refuses_the_pinned_team proves the pin's own subject is a
// hard stop end to end: the fixture pin is team-id.
func Test_TeamDelete_refuses_the_pinned_team(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"team", "delete", "team-id", "--org-wide"})

	require.ErrorIs(t, err, client.ErrTargetMismatch)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "pinned team")
}

// Test_TeamDelete_writes_the_deletion is the allow arm: another team in the
// organization is archived and the irreversibility warning is printed.
func Test_TeamDelete_writes_the_deletion(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"team", "delete", "ops-team-id", "--org-wide"})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "archived team ops-team-id")
	require.Contains(t, stdout.String(), "cannot be undone via linctl")
	require.Empty(t, stderr.String())
}
