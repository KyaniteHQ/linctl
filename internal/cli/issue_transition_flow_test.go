package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

func useTransitionRuntime(t *testing.T, allowlist config.Transitions) func() {
	t.Helper()
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		runtime := testCommandRuntime(wrapCommandFlowAfterWrite(commandFlowFakeClient{}))
		runtime.config.Target.Transitions = allowlist

		return runtime, nil
	}

	return func() {
		buildCommandRuntime = original
	}
}

// Test_IssueUpdate_refuses_a_transition_outside_the_allowlist proves the
// allowlist is enforced end to end from `.linctl.toml`: the fixture issue sits
// in Todo, the allowlist lets Todo move only to Done, and the JSON envelope
// carries the stable code so an agent can stop instead of retrying.
func Test_IssueUpdate_refuses_a_transition_outside_the_allowlist(t *testing.T) {
	restore := useTransitionRuntime(t, config.Transitions{"Todo": {"Done"}})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"--json", "issue", "update", "LIT-1", "--state", "In Review"})

	require.ErrorIs(t, err, client.ErrTransitionDenied)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), `"error_code":"TRANSITION_DENIED"`)
}

// Test_IssueStart_moves_when_the_allowlist_lists_the_started_state is the allow
// arm: `issue start` resolves the started state and Todo -> Started is listed.
func Test_IssueStart_moves_when_the_allowlist_lists_the_started_state(t *testing.T) {
	restore := useTransitionRuntime(t, config.Transitions{"Todo": {"Started"}})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"issue", "start", "LIT-1"})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "LIT-1")
	require.Empty(t, stderr.String())
}
