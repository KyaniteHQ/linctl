package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_writeDocumentDetail_prints_only_the_line_without_content(t *testing.T) {
	command := &cobra.Command{}
	out := &bytes.Buffer{}
	command.SetOut(out)
	document := client.DocumentDetail{
		DocumentSummary: client.DocumentSummary{ID: "document-id", Title: "Project note", ParentType: "project"},
	}

	err := writeDocumentDetail(command, &rootOptions{}, document)

	require.NoError(t, err)
	require.Equal(t, "document-id Project note [project]\n", out.String())
}

func Test_DocumentGet_surfaces_stdout_write_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(failingWriter{err: errors.New("stdout closed")})
	command.SetArgs([]string{"document", "get", "document-id"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "stdout closed")
}
