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
