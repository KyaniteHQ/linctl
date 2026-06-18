package render

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WriteJSON_writes_indented_json_when_value_is_structured(t *testing.T) {
	// Given
	buffer := bytes.Buffer{}
	value := map[string]string{"org_id": "org-id"}

	// When
	err := WriteJSON(&buffer, value)

	// Then
	require.NoError(t, err)
	require.JSONEq(t, `{"org_id":"org-id"}`, buffer.String())
}
