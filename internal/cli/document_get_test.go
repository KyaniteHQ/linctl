package cli

import (
	"bytes"
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
		DocumentSummary: client.DocumentSummary{ID: "document-id", Title: "Team note", ParentType: "team"},
	}

	err := writeDocumentDetail(command, &rootOptions{}, document)

	require.NoError(t, err)
	require.Equal(t, "document-id Team note [team]\n", out.String())
}

func Test_writeDocumentDetail_surfaces_stdout_write_errors(t *testing.T) {
	command := &cobra.Command{}
	command.SetOut(failingWriter{err: errors.New("stdout closed")})
	document := client.DocumentDetail{
		DocumentSummary: client.DocumentSummary{ID: "document-id", Title: "Team note"},
		Content:         "body",
	}

	err := writeDocumentDetail(command, &rootOptions{}, document)

	require.ErrorContains(t, err, "stdout closed")
}
