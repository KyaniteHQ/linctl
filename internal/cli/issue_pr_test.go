package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func runIssuePRWithHostileTitle(t *testing.T, args []string) string {
	t.Helper()
	restore := useCommandRuntime(t, graphqlPayloadOverride{
		inner: commandFlowFakeClient{},
		payloads: map[string]string{
			"issue": `{"issue":` + commandIssueJSON(
				"LIT-1", "fix `touch pwned`, $(id) and it's", "todo-state", "Todo", "unstarted",
			) + `}`,
		},
	})
	defer restore()
	output := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs(args)

	require.NoError(t, command.ExecuteContext(context.Background()))

	return output.String()
}

func Test_IssuePR_single_quotes_command_substitution_and_embedded_quotes(t *testing.T) {
	output := runIssuePRWithHostileTitle(t, []string{"issue", "pr", "LIT-1"})

	require.Equal(
		t,
		"gh pr create --title 'LIT-1 fix `touch pwned`, $(id) and it'\\''s' "+
			"--body 'https://linear.app/kyanite/issue/LIT-1'\n",
		output,
	)
}

func Test_IssuePR_json_keeps_argv_form(t *testing.T) {
	output := runIssuePRWithHostileTitle(t, []string{"--json", "--compact", "issue", "pr", "LIT-1"})

	require.JSONEq(t, `{
		"title": "LIT-1 fix `+"`touch pwned`"+`, $(id) and it's",
		"body": "https://linear.app/kyanite/issue/LIT-1",
		"command": [
			"gh", "pr", "create",
			"--title", "LIT-1 fix `+"`touch pwned`"+`, $(id) and it's",
			"--body", "https://linear.app/kyanite/issue/LIT-1"
		]
	}`, output)
}
