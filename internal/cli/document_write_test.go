package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_DocumentWrites_refuse_without_a_pinned_project_before_reading_content(t *testing.T) {
	for _, args := range [][]string{
		{"document", "create", "--title", "x", "--content", "-"},
		{"document", "update", "document-id", "--content", "-"},
	} {
		t.Run(args[1], func(t *testing.T) {
			recorder := &recordingGraphQLClient{}
			original := buildCommandRuntime
			buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				runtime := testCommandRuntime(recorder)
				runtime.config.Target.ProjectID = ""
				return runtime, nil
			}
			defer func() { buildCommandRuntime = original }()

			var stdout, stderr bytes.Buffer
			err := execute(context.Background(), BuildInfo{}, initiativeBodyFailingReader{}, &stdout, &stderr, args)

			require.ErrorIs(t, err, client.ErrWriteInvalid)
			require.ErrorContains(t, err, "project_id")
			require.Empty(t, recorder.calls)
		})
	}
}
