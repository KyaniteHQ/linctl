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
		{name: "custom_views", key: "custom_views"},
		{name: "project_labels", key: "project_labels"},
		{name: "project_statuses", key: "project_statuses"},
		{name: "spans", key: "spans"},
		{name: "git_automation_states", key: "git_automation_states"},
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

// Test_projectJSONFields_leaves_detail_arrays_whole guards the curation rule in
// projectCollection: the projected-collection key set is an allowlist, not
// generic top-level []any detection. A detail object carries scalar fields plus
// an incidental array that is NOT a collection (a time schedule's "entries"),
// and a dependency graph carries several arrays at once. Both must project as a
// single object — projecting per-element would return the wrong entity. If a
// future change replaces the allowlist with "the single top-level array", these
// cases fail.
func Test_projectJSONFields_leaves_detail_arrays_whole(t *testing.T) {
	t.Run("detail with incidental array", func(t *testing.T) {
		detail := map[string]any{
			"id":      "schedule-1",
			"name":    "On call",
			"entries": []any{map[string]any{"id": "entry-1"}},
		}

		projected, err := projectJSONFields(detail, "id,name")

		require.NoError(t, err)
		require.Equal(t, map[string]any{"id": "schedule-1", "name": "On call"}, projected)
	})

	t.Run("multiple top-level arrays", func(t *testing.T) {
		graph := map[string]any{
			"id":         "issue-1",
			"children":   []any{map[string]any{"id": "child-1"}},
			"blocks":     []any{},
			"blocked_by": []any{},
		}

		projected, err := projectJSONFields(graph, "id")

		require.NoError(t, err)
		require.Equal(t, map[string]any{"id": "issue-1"}, projected)
	})
}
