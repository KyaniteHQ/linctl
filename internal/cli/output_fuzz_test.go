package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// FuzzProjectJSONFields asserts the field projection parser never panics on
// arbitrary --fields strings and always returns a non-nil value on success.
func FuzzProjectJSONFields(f *testing.F) {
	seeds := []string{
		"",
		"id",
		"id,name",
		"nested.deep.leaf",
		" , , ",
		"id,,name",
		"...",
		"a.",
		".b",
		"issues",
		"missing.path.here",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	value := map[string]any{
		"id":     "x",
		"name":   "y",
		"nested": map[string]any{"deep": map[string]any{"leaf": 1}},
		"issues": []any{map[string]any{"id": "i1", "name": "n1"}},
	}

	f.Fuzz(func(t *testing.T, fields string) {
		result, err := projectJSONFields(value, fields)
		if err != nil {
			return
		}

		require.NotNil(t, result)
	})
}

// FuzzSortByJSONField asserts the list sorter never panics on arbitrary field
// and order strings and preserves the element count whenever it succeeds.
func FuzzSortByJSONField(f *testing.F) {
	seeds := []struct {
		field string
		order string
	}{
		{field: "", order: "asc"},
		{field: "id", order: "asc"},
		{field: "id", order: "desc"},
		{field: "name", order: "sideways"},
		{field: "missing", order: "asc"},
		{field: "nested.leaf", order: "desc"},
		{field: ".", order: ""},
	}
	for _, seed := range seeds {
		f.Add(seed.field, seed.order)
	}

	items := []map[string]any{
		{"id": "2", "name": "b", "nested": map[string]any{"leaf": "y"}},
		{"id": "1", "name": "a", "nested": map[string]any{"leaf": "x"}},
		{"id": "3", "name": "c", "nested": map[string]any{"leaf": "z"}},
	}

	f.Fuzz(func(t *testing.T, field string, order string) {
		sorted, err := sortByJSONField(items, field, order)
		if err != nil {
			return
		}

		require.Len(t, sorted, len(items))
	})
}
