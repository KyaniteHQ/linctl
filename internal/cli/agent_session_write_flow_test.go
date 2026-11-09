package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
