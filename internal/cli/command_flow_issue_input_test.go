package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CommandFlows_read_issue_comment_body_from_stdin(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{expectedCommentBody: "stdin body\nsecond line"})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(strings.NewReader("stdin body\nsecond line"))
	command.SetOut(&output)
	command.SetArgs([]string{"issue", "comment", "LIT-1", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "comment comment-id on LIT-1")
}

func Test_CommandFlows_read_issue_text_from_files(t *testing.T) {
	descriptionFile := writeTempTextFile(t, "description from file")
	appendFile := writeTempTextFile(t, "append from file")
	commentFile := writeTempTextFile(t, "comment from file")
	replyFile := writeTempTextFile(t, "reply from file")

	tests := []struct {
		name string
		args []string
		fake commandFlowFakeClient
	}{
		{
			name: "create description",
			args: []string{"issue", "create", "--title", "Created issue", "--description-file", descriptionFile},
			fake: commandFlowFakeClient{expectedCreateDescription: "description from file"},
		},
		{
			name: "update append",
			args: []string{"issue", "update", "LIT-1", "--append-file", appendFile},
			fake: commandFlowFakeClient{expectedUpdateDescription: "Existing description\n\nappend from file"},
		},
		{
			name: "comment body",
			args: []string{"issue", "comment", "LIT-1", "--body-file", commentFile},
			fake: commandFlowFakeClient{expectedCommentBody: "comment from file"},
		},
		{
			name: "reply body",
			args: []string{"issue", "reply", "LIT-1", "comment-id", "--body-file", replyFile},
			fake: commandFlowFakeClient{expectedCommentBody: "reply from file", expectedCommentParentID: "comment-id"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.NotEmpty(t, output.String())
		})
	}
}

func Test_CommandFlows_report_issue_text_file_errors(t *testing.T) {
	textFile := writeTempTextFile(t, "from file")
	missingFile := filepath.Join(t.TempDir(), "missing.md")
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "create description conflict",
			args:     []string{"issue", "create", "--title", "Created issue", "--description", "inline", "--description-file", textFile},
			contains: "description and description-file are mutually exclusive",
		},
		{
			name:     "update description conflict",
			args:     []string{"issue", "update", "LIT-1", "--description", "inline", "--description-file", textFile},
			contains: "description and description-file are mutually exclusive",
		},
		{
			name:     "update append conflict",
			args:     []string{"issue", "update", "LIT-1", "--append", "inline", "--append-file", textFile},
			contains: "append and append-file are mutually exclusive",
		},
		{
			name:     "comment body conflict",
			args:     []string{"issue", "comment", "LIT-1", "--body", "inline", "--body-file", textFile},
			contains: "body and body-file are mutually exclusive",
		},
		{
			name:     "reply body conflict",
			args:     []string{"issue", "reply", "LIT-1", "comment-id", "--body", "inline", "--body-file", textFile},
			contains: "body and body-file are mutually exclusive",
		},
		{
			name:     "missing file",
			args:     []string{"issue", "comment", "LIT-1", "--body-file", missingFile},
			contains: "read body from file",
		},
		{
			name:     "create unknown state alias",
			args:     []string{"issue", "create", "--title", "T", "--state", "sprinting"},
			contains: "unknown state type",
		},
		{
			name:     "create unknown priority alias",
			args:     []string{"issue", "create", "--title", "T", "--priority", "blocker"},
			contains: "unknown priority",
		},
		{
			name:     "update unknown state alias",
			args:     []string{"issue", "update", "LIT-1", "--state", "sprinting"},
			contains: "unknown state type",
		},
		{
			name:     "update unknown priority alias",
			args:     []string{"issue", "update", "LIT-1", "--priority", "blocker"},
			contains: "unknown priority",
		},
		{
			name:     "list unknown status alias",
			args:     []string{"issue", "list", "--status", "sprinting"},
			contains: "unknown state type",
		},
		{
			name:     "project update body conflict",
			args:     []string{"project-update", "create", "project-id", "--body", "inline", "--body-file", textFile},
			contains: "body and body-file are mutually exclusive",
		},
		{
			name:     "comment update body conflict",
			args:     []string{"comment", "update", "comment-id", "--body", "inline", "--body-file", textFile},
			contains: "body and body-file are mutually exclusive",
		},
		{
			name:     "project update unknown health alias",
			args:     []string{"project-update", "create", "project-id", "--health", "sideways"},
			contains: "unknown health",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.contains)
		})
	}
}
