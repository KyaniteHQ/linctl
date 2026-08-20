package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_pageWithItems_preserves_page_shape_and_collection_semantics(t *testing.T) {
	type item struct {
		ID string `json:"id"`
	}
	type page struct {
		Items         []item `json:"items"`
		HasNextPage   bool   `json:"has_next_page"`
		EndCursor     string `json:"end_cursor"`
		ContextMarker string `json:"context_marker"`
	}

	original := page{
		Items:         []item{{ID: "original"}},
		HasNextPage:   true,
		EndCursor:     "cursor",
		ContextMarker: "context",
	}
	items := []item{{ID: "b"}, {ID: "a"}}
	sorted, err := sortByJSONField(items, "id", "asc")
	require.NoError(t, err)

	result, err := pageWithItems(original, sorted)
	require.NoError(t, err)
	require.Equal(t, []item{{ID: "original"}}, original.Items)
	require.True(t, bytes.Equal(
		[]byte(`{"items":[{"id":"a"},{"id":"b"}],"has_next_page":true,"end_cursor":"cursor","context_marker":"context"}`),
		compactJSON(t, result),
	))

	nilResult, err := pageWithItems(original, []item(nil))
	require.NoError(t, err)
	require.True(t, bytes.Equal(
		[]byte(`{"items":null,"has_next_page":true,"end_cursor":"cursor","context_marker":"context"}`),
		compactJSON(t, nilResult),
	))

	emptyResult, err := pageWithItems(original, []item{})
	require.NoError(t, err)
	require.True(t, bytes.Equal(
		[]byte(`{"items":[],"has_next_page":true,"end_cursor":"cursor","context_marker":"context"}`),
		compactJSON(t, emptyResult),
	))

	pointerOriginal := &original
	pointerResult, err := pageWithItems(pointerOriginal, sorted)
	require.NoError(t, err)
	require.NotSame(t, pointerOriginal, pointerResult)
	require.Equal(t, []item{{ID: "original"}}, pointerOriginal.Items)
	require.Equal(t, sorted, pointerResult.Items)
}

func Test_pageWithItems_rejects_malformed_shapes(t *testing.T) {
	type noCollection struct {
		Value string `json:"value"`
	}
	type ignoredCollection struct {
		Items []string `json:"-"`
	}
	type ambiguousCollection struct {
		Items  []string `json:"items"`
		Labels []string `json:"labels"`
	}
	type wrongCollection struct {
		Items []int `json:"items"`
	}
	type pointerChain **struct {
		Items []string `json:"items"`
	}

	_, err := pageWithItems(any(nil), []string{"item"})
	require.ErrorContains(t, err, "value is invalid")
	_, err = pageWithItems("not a page", []string{"item"})
	require.ErrorContains(t, err, "must be a struct or pointer to a struct")
	_, err = pageWithItems(noCollection{}, []string{"item"})
	require.ErrorContains(t, err, "no exported JSON slice field")
	hiddenCollection := reflect.StructOf([]reflect.StructField{{
		Name:    "items",
		PkgPath: "github.com/KyaniteHQ/linctl/internal/cli",
		Type:    reflect.TypeOf([]string(nil)),
		Tag:     `json:"items"`,
	}})
	_, err = pageCollectionField(hiddenCollection)
	require.ErrorContains(t, err, "no exported JSON slice field")
	_, err = pageWithItems(ignoredCollection{}, []string{"item"})
	require.ErrorContains(t, err, "no exported JSON slice field")
	_, err = pageWithItems(ambiguousCollection{}, []string{"item"})
	require.ErrorContains(t, err, "multiple exported JSON slice fields")
	_, err = pageWithItems(wrongCollection{}, []string{"item"})
	require.ErrorContains(t, err, `collection field "Items" has type []int, not []string`)
	var nilPage *struct {
		Items []string `json:"items"`
	}
	_, err = pageWithItems(nilPage, []string{"item"})
	require.ErrorContains(t, err, "must not be nil")
	var chainedPage pointerChain
	_, err = pageWithItems(chainedPage, []string{"item"})
	require.ErrorContains(t, err, "must point directly to a struct")
	nested := &struct {
		Items []string `json:"items"`
	}{}
	chainedPage = &nested
	_, err = pageWithItems(chainedPage, []string{"item"})
	require.ErrorContains(t, err, "must point directly to a struct")
}

func Test_writePageJSON_rejects_malformed_page(t *testing.T) {
	command := &cobra.Command{}
	command.SetOut(&bytes.Buffer{})

	err := writePageJSON(command, &rootOptions{json: true, compact: true}, "not a page", []string{"item"})

	require.ErrorContains(t, err, "must be a struct or pointer to a struct")
}

func compactJSON(t *testing.T, value any) []byte {
	t.Helper()
	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	return encoded
}

func Test_runReadListCommand_notes_truncation_on_stderr(t *testing.T) {
	type item struct {
		ID string `json:"id"`
	}
	type page struct {
		Items []item `json:"items"`
		client.Page
	}
	original := buildCommandRuntime
	buildCommandRuntime = func(context.Context, *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, nil
	}
	t.Cleanup(func() {
		buildCommandRuntime = original
	})

	tests := []struct {
		name       string
		options    rootOptions
		truncated  bool
		failStderr bool
		wantStdout string
		wantStderr string
		wantErr    string
	}{
		{
			name:       "text",
			truncated:  true,
			wantStdout: "item-id\n",
			wantStderr: "note: listing capped at 50 items; more pages exist\n",
		},
		{
			name:       "id-only",
			options:    rootOptions{idOnly: true},
			truncated:  true,
			wantStdout: "item-id\n",
			wantStderr: "note: listing capped at 50 items; more pages exist\n",
		},
		{
			name:       "quiet",
			options:    rootOptions{quiet: true},
			truncated:  true,
			wantStdout: "",
			wantStderr: "",
		},
		{
			name:       "json",
			options:    rootOptions{json: true, compact: true},
			truncated:  true,
			wantStdout: `{"items":[{"id":"item-id"}],"has_next_page":true,"end_cursor":"next"}` + "\n",
			wantStderr: "note: listing capped at 50 items; more pages exist\n",
		},
		{
			name:       "complete page",
			wantStdout: "item-id\n",
			wantStderr: "",
		},
		{
			name:       "stderr write error",
			truncated:  true,
			failStderr: true,
			wantErr:    "write failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command, stdout, stderr := bufferedCommand()
			if test.failStderr {
				command.SetErr(commandFailingWriter{})
			}
			loaded := page{
				Items: []item{{ID: "item-id"}},
				Page:  client.Page{HasNextPage: test.truncated},
			}
			if test.truncated {
				cursor := "next"
				loaded.EndCursor = &cursor
			}

			err := runReadListCommand(
				context.Background(),
				command,
				nil,
				&test.options,
				50,
				func(context.Context, commandRuntime, []string, int) (page, error) {
					return loaded, nil
				},
				func(command *cobra.Command, options *rootOptions, row item) error {
					return writeItemLine(command, options, row, row.ID, "%s", row.ID)
				},
			)

			if test.wantErr != "" {
				require.ErrorContains(t, err, test.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.wantStdout, stdout.String())
			require.Equal(t, test.wantStderr, stderr.String())
		})
	}
}
