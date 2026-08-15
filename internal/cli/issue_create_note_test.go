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

func Test_IssueCreate_notes_the_project_it_landed_in(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{expectedCreateTitle: "Created issue"})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"issue", "create", "--title", "Created issue"})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "LIT-2 Created issue")
	require.Equal(t, "note: created in project \"Pinned project\"\n", stderr.String())
}

func Test_IssueCreate_stays_quiet_when_a_parent_proves_the_project(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{expectedCreateTitle: "Created issue"})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr,
		[]string{"issue", "create", "--title", "Created issue", "--parent", "LIT-1"})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "LIT-2 Created issue")
	require.Empty(t, stderr.String())
}

func Test_IssueCreate_reports_a_closed_stderr(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{expectedCreateTitle: "Created issue"})
	defer restore()

	var stdout bytes.Buffer
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout,
		failingWriter{err: errors.New("stderr closed")},
		[]string{"issue", "create", "--title", "Created issue"})

	require.ErrorContains(t, err, "stderr closed")
}

func Test_NoteCreatedProject_names_the_team_when_the_issue_has_no_project(t *testing.T) {
	var stderr bytes.Buffer
	command := (cobraCommandWithIO{err: &stderr}).command()

	err := noteCreatedProject(command, "", client.IssueSummary{Team: "LIT"})

	require.NoError(t, err)
	require.Equal(t, "note: created in team LIT with no project\n", stderr.String())
}
