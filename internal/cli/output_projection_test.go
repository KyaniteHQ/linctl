package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_projectJSONFields_projects_list_envelopes guards that --json --fields
// projects per item of the collection named by a command's annotated
// collection key, across representative envelope keys.
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

			projected, err := projectJSONFieldsWithCollectionKey(envelope, "id,name", testCase.key)

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

// Test_projectJSONFields_leaves_detail_arrays_whole guards that a command
// without an annotated collection key always projects the whole object: a
// detail object's embedded arrays (a time schedule's "entries", a project's
// "teams", a dependency graph's several arrays) are never mistaken for a
// projection collection, no matter what the array is named.
func Test_projectJSONFields_leaves_detail_arrays_whole(t *testing.T) {
	t.Run("detail with incidental array", func(t *testing.T) {
		detail := map[string]any{
			"id":      "schedule-1",
			"name":    "On call",
			"entries": []any{map[string]any{"id": "entry-1"}},
		}

		projected, err := projectJSONFieldsWithCollectionKey(detail, "id,name", "")

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

		projected, err := projectJSONFieldsWithCollectionKey(graph, "id", "")

		require.NoError(t, err)
		require.Equal(t, map[string]any{"id": "issue-1"}, projected)
	})

	t.Run("detail with collection-named array", func(t *testing.T) {
		detail := map[string]any{
			"name":  "the-project",
			"teams": []any{map[string]any{"name": "team-one"}, map[string]any{"name": "team-two"}},
		}

		projected, err := projectJSONFieldsWithCollectionKey(detail, "name", "")

		require.NoError(t, err)
		require.Equal(t, map[string]any{"name": "the-project"}, projected)
	})
}
