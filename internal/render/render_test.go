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
	err := WriteJSON(&buffer, value, false)

	// Then
	require.NoError(t, err)
	require.Equal(t, "{\n  \"org_id\": \"org-id\"\n}\n", buffer.String()) //nolint:testifylint // indentation is the asserted behavior.
}
