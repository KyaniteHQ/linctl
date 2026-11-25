package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_projectJSONFields_projects_list_envelopes guards that --json --fields
// projects per item of the collection named by a command's annotated
// collection key.
func Test_projectJSONFields_projects_list_envelopes(t *testing.T) {
	envelope := map[string]any{
		"favorites": []any{
			map[string]any{"id": "id-1", "name": "first", "extra": "drop"},
		},
	}

	projected, err := projectJSONFieldsWithCollectionKey(envelope, "id,name", "favorites")

	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"favorites": []any{map[string]any{"id": "id-1", "name": "first"}},
	}, projected)
}

// Test_projectJSONFields_leaves_detail_arrays_whole guards that a command
// without an annotated collection key projects the whole object: a detail
// object's embedded array is never mistaken for a projection collection.
func Test_projectJSONFields_leaves_detail_arrays_whole(t *testing.T) {
	detail := map[string]any{
		"id":      "schedule-1",
		"name":    "On call",
		"entries": []any{map[string]any{"id": "entry-1"}},
	}

	projected, err := projectJSONFieldsWithCollectionKey(detail, "id,name", "")

	require.NoError(t, err)
	require.Equal(t, map[string]any{"id": "schedule-1", "name": "On call"}, projected)
}
