package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_WritePathInputGuards_report_INVALID_WRITE(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "archive empty project id",
			args: []string{"project", "archive", ""},
			want: "project id is required",
		},
		{
			name: "label color",
			args: []string{"label", "create", "--name", "Bug", "--color", "notacolor"},
			want: "color must be #RRGGBB",
		},
		{
			name: "project label color",
			args: []string{"project-label", "create", "--name", "Bug", "--color", "notacolor", "--org-wide"},
			want: "color must be #RRGGBB",
		},
		{
			name: "self relation two reference forms",
			args: []string{"issue", "relate", "LIT-42", "issue-uuid-42", "--type", "blocks"},
			want: "cannot relate to itself",
		},
		{
			name: "template data-file stdin",
			args: []string{
				"template", "create",
				"--id", "7c9e6679-7425-40de-944b-e07fc1f90ae7",
				"--name", "Bug report",
				"--type", "issue",
				"--data-file", "-",
			},
			want: "not stdin",
		},
		{
			name: "workflow-state type",
			args: []string{
				"workflow-state", "create",
				"--id", "550e8400-e29b-41d4-a716-446655440000",
				"--name", "Ready",
				"--type", "triage",
				"--color", "#f2c94c",
			},
			want: "type must be backlog, unstarted, started, completed, or canceled",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			var stdout, stderr bytes.Buffer
			err := execute(
				context.Background(),
				BuildInfo{},
				strings.NewReader(""),
				&stdout,
				&stderr,
				tt.args,
			)

			require.ErrorIs(t, err, client.ErrWriteInvalid)
			require.Equal(t, "INVALID_WRITE", errorCode(err))
			require.ErrorContains(t, err, tt.want)
			require.Contains(t, stderr.String(), `"error_code":"INVALID_WRITE"`)
		})
	}
}
