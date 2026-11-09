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

func Test_AgentSessionCreateCommandFlow_reports_runtime_errors(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"agent-session", "create", "--issue", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "runtime failed")
}

func Test_AgentSessionCreateCommandFlow_reports_writer_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"agent-session", "create", "--issue", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "write failed")
}

func Test_AgentSessionCreate_requires_an_issue(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"agent-session", "create"})

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "issue is required")
}

func Test_AgentSessionCreate_writes_the_created_session(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"agent-session", "create", "--issue", "LIT-1"})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "agent-session-id")
	require.Empty(t, stderr.String())
}
