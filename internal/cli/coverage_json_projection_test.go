package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CliOutputHelpers_cover_json_projection_and_sort_edges(t *testing.T) {
	projected, err := projectJSONFieldsWithCollectionKey(
		map[string]any{"issues": []any{map[string]any{"identifier": "LIT-1", "state": map[string]any{"name": "Todo"}}}},
		"identifier,state.name",
		"issues",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"issues": []any{map[string]any{"identifier": "LIT-1", "state": map[string]any{"name": "Todo"}}},
	}, projected)

	projected, err = projectJSONFieldsWithCollectionKey(map[string]any{"identifier": "LIT-1"}, "identifier", "")
	require.NoError(t, err)
	require.Equal(t, map[string]any{"identifier": "LIT-1"}, projected)

	projected, err = projectJSONFieldsWithCollectionKey(map[string]any{"identifier": "LIT-1"}, "identifier,, ", "")
	require.NoError(t, err)
	require.Equal(t, map[string]any{"identifier": "LIT-1"}, projected)

	// An annotated collection key that is absent from the payload falls back to
	// whole-object projection instead of failing.
	projected, err = projectJSONFieldsWithCollectionKey(map[string]any{"identifier": "LIT-1"}, "identifier", "issues")
	require.NoError(t, err)
	require.Equal(t, map[string]any{"identifier": "LIT-1"}, projected)

	_, err = projectJSONFieldsWithCollectionKey(map[string]any{"bad": func() {}}, "bad", "")
	require.ErrorContains(t, err, "marshal output")

	_, err = projectJSONFieldsWithCollectionKey([]string{"not-an-object"}, "id", "")
	require.ErrorContains(t, err, "decode output")

	_, err = projectJSONFieldsWithCollectionKey(map[string]any{"issues": []any{"bad-item"}}, "identifier", "issues")
	require.ErrorContains(t, err, "item is not an object")

	_, err = projectJSONFieldsWithCollectionKey(
		map[string]any{"issues": []any{map[string]any{"title": "Missing id"}}}, "identifier", "issues",
	)
	require.ErrorContains(t, err, "field \"identifier\" is not present")

	_, err = projectJSONFieldsWithCollectionKey(map[string]any{"identifier": "LIT-1"}, "missing", "")
	require.ErrorContains(t, err, "field \"missing\" is not present")

	_, err = projectJSONFieldsWithCollectionKey(map[string]any{"state": "Todo"}, "state.name", "")
	require.ErrorContains(t, err, "field \"state\" is not an object")

	items := []client.IssueSummary{
		{Identifier: "LIT-2", Title: "Zebra"},
		{Identifier: "LIT-1", Title: "Alpha"},
	}
	sortedItems, err := sortByJSONField(items, "", "asc")
	require.NoError(t, err)
	require.Equal(t, items, sortedItems)

	sortedItems, err = sortByJSONField(items, "title", "asc")
	require.NoError(t, err)
	require.Equal(t, "Alpha", sortedItems[0].Title)

	_, err = sortByJSONField(items, "title", "sideways")
	require.ErrorContains(t, err, "invalid sort order")

	_, err = sortByJSONField(items, "missing", "asc")
	require.ErrorContains(t, err, "sort field \"missing\" is not present")

	_, err = sortByJSONField([]map[string]any{{"state": "Todo"}}, "state.name", "asc")
	require.ErrorContains(t, err, "not an object path")

	var nilItems []client.IssueSummary
	sortedItems, err = sortByJSONField(nilItems, "title", "asc")
	require.NoError(t, err)
	require.Nil(t, sortedItems)

	type numbered struct {
		Number float64 `json:"number"`
	}
	numberedItems := []numbered{{Number: 2}, {Number: 10}, {Number: 1}}
	sortedNumbers, err := sortByJSONField(numberedItems, "number", "asc")
	require.NoError(t, err)
	require.Equal(t, []numbered{{Number: 1}, {Number: 2}, {Number: 10}}, sortedNumbers)

	sortedNumbers, err = sortByJSONField(numberedItems, "number", "desc")
	require.NoError(t, err)
	require.Equal(t, []numbered{{Number: 10}, {Number: 2}, {Number: 1}}, sortedNumbers)

	type counted struct {
		Count int `json:"count"`
	}
	countedItems := []counted{{Count: 2}, {Count: 10}, {Count: 1}}
	sortedCounts, err := sortByJSONField(countedItems, "count", "asc")
	require.NoError(t, err)
	require.Equal(t, []counted{{Count: 1}, {Count: 2}, {Count: 10}}, sortedCounts)

	type mixed struct {
		Value any `json:"value"`
	}
	mixedItems := []mixed{{Value: 10.0}, {Value: "abc"}, {Value: 2.0}}
	sortedMixed, err := sortByJSONField(mixedItems, "value", "asc")
	require.NoError(t, err)
	require.Equal(t, []mixed{{Value: 2.0}, {Value: 10.0}, {Value: "abc"}}, sortedMixed)

	sortedMixed, err = sortByJSONField(mixedItems, "value", "desc")
	require.NoError(t, err)
	require.Equal(t, []mixed{{Value: "abc"}, {Value: 10.0}, {Value: 2.0}}, sortedMixed)

	_, err = sortByJSONField([]map[string]any{{"bad": func() {}}}, "bad", "asc")
	require.ErrorContains(t, err, "marshal output")

	destination := map[string]any{}
	require.NoError(t, copyJSONPath(map[string]any{"id": "issue-id"}, destination, nil))
	require.Empty(t, destination)
}
