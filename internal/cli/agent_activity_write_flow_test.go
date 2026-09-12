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

func Test_AgentActivityCreateCommandFlow_reports_runtime_errors(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"agent-activity", "create", "agent-session-id", "--type", "thought", "--body", "Reading"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "runtime failed")
}

func Test_AgentActivityCreateCommandFlow_reports_writer_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"agent-activity", "create", "agent-session-id", "--type", "thought", "--body", "Reading"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "write failed")
}

// Test_AgentActivityCreate_refuses_an_unknown_type proves the end-to-end
// refusal: the accepted types are named on stderr and nothing reaches stdout.
func Test_AgentActivityCreate_refuses_an_unknown_type(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"agent-activity", "create", "agent-session-id", "--type", "prompt", "--body", "Reading"})

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "thought, elicitation, response, error, action")
}

// Test_AgentActivityCreate_writes_the_created_activity is the allow arm: the
// fixture session sits on LIT-1, which is on the pinned team and project.
func Test_AgentActivityCreate_writes_the_created_activity(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{
			"agent-activity", "create", "agent-session-id", "--type", "action",
			"--action", "read_file", "--parameter", "README.md", "--signal", "continue",
		})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "agent-activity-id session agent-session-id")
	require.Empty(t, stderr.String())
}
