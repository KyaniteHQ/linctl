package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HasMore_reports_whether_further_pages_exist(t *testing.T) {
	require.True(t, Page{HasNextPage: true}.HasMore())
	require.False(t, Page{}.HasMore())

	// TemplateList and SemanticSearchList keep their own shape instead of
	// embedding Page, so each carries its own HasMore.
	require.True(t, TemplateList{HasNextPage: true}.HasMore())
	require.False(t, TemplateList{}.HasMore())
	require.False(t, SemanticSearchList{}.HasMore())
}
