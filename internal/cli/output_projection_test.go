package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_projectJSONFields_projects_list_envelopes guards that --json --fields
// projects per item for every collection envelope the read commands emit.
// The favorite/emoji/attachment envelopes were previously absent from the
// projectCollection dispatch table, so projection fell through to the
// single-object path and errored on the (missing) leaf field.
func Test_projectJSONFields_projects_list_envelopes(t *testing.T) {
	cases := []struct {
		name string
		key  string
	}{
		{name: "favorites", key: "favorites"},
		{name: "emojis", key: "emojis"},
		{name: "attachments", key: "attachments"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			envelope := map[string]any{
				testCase.key: []any{
					map[string]any{"id": "id-1", "name": "first", "extra": "drop"},
				},
			}

			projected, err := projectJSONFields(envelope, "id,name")

			require.NoError(t, err)
			result, ok := projected.(map[string]any)
			require.True(t, ok)
			items, ok := result[testCase.key].([]any)
			require.True(t, ok)
			require.Len(t, items, 1)
			item, ok := items[0].(map[string]any)
			require.True(t, ok)
			require.Equal(t, "id-1", item["id"])
			require.Equal(t, "first", item["name"])
			require.NotContains(t, item, "extra")
		})
	}
}
