package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_shellQuote_uses_posix_single_quotes(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "empty", value: "", want: "''"},
		{name: "plain", value: "LIT-1 Detail issue", want: "'LIT-1 Detail issue'"},
		{name: "backtick substitution", value: "fix `touch pwned`", want: "'fix `touch pwned`'"},
		{name: "dollar substitution", value: "fix $(id)", want: "'fix $(id)'"},
		{name: "embedded single quote", value: "it's broken", want: `'it'\''s broken'`},
		{name: "double quotes stay literal", value: `say "hi"`, want: `'say "hi"'`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, shellQuote(test.value))
		})
	}
}

func Test_writePullRequestPlan_single_quotes_command_substitution(t *testing.T) {
	output := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	plan := pullRequestPlanFromIssue(client.IssueSummary{
		Identifier: "LIT-1",
		Title:      "fix `touch pwned` and $(id)",
		URL:        "https://linear.app/kyanite/issue/LIT-1",
	})

	err := writePullRequestPlan(command, &rootOptions{}, plan)

	require.NoError(t, err)
	require.Equal(
		t,
		"gh pr create --title 'LIT-1 fix `touch pwned` and $(id)' --body 'https://linear.app/kyanite/issue/LIT-1'\n",
		output.String(),
	)
}

func Test_writePullRequestPlan_json_keeps_argv_form(t *testing.T) {
	output := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	plan := pullRequestPlanFromIssue(client.IssueSummary{
		Identifier: "LIT-1",
		Title:      "fix `touch pwned`",
		URL:        "https://linear.app/kyanite/issue/LIT-1",
	})

	err := writePullRequestPlan(command, &rootOptions{json: true}, plan)

	require.NoError(t, err)
	require.Contains(t, output.String(), `"title": "LIT-1 fix `+"`touch pwned`"+`"`)
	require.Contains(t, output.String(), `"gh"`)
	require.Contains(t, output.String(), `"--title"`)
	require.NotContains(t, output.String(), "gh pr create --title")
}
