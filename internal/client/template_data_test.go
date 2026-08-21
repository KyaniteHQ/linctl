package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CanonicalTemplateData_ignores_whitespace_and_key_order(t *testing.T) {
	left, err := CanonicalTemplateData([]byte(`{ "b": 1, "a": { "z": true, "y": "ok" } }`))
	require.NoError(t, err)
	right, err := CanonicalTemplateData([]byte(`{"a":{"y":"ok","z":true},"b":1}`))
	require.NoError(t, err)
	require.Equal(t, string(left), string(right))
	require.JSONEq(t, `{"a":{"y":"ok","z":true},"b":1}`, string(left))
}

func Test_CanonicalTemplateData_preserves_number_lexemes(t *testing.T) {
	one, err := CanonicalTemplateData([]byte(`{"n":1}`))
	require.NoError(t, err)
	onePointZero, err := CanonicalTemplateData([]byte(`{"n":1.0}`))
	require.NoError(t, err)
	require.Equal(t, `{"n":1}`, string(one))            //nolint:testifylint // 1 and 1.0 must stay distinct
	require.Equal(t, `{"n":1.0}`, string(onePointZero)) //nolint:testifylint // 1 and 1.0 must stay distinct
	require.NotEqual(t, string(one), string(onePointZero))
}

func Test_CanonicalTemplateData_unwraps_one_encoded_object_layer(t *testing.T) {
	canonical, err := CanonicalTemplateData([]byte(`"{\"title\":\"Bug\",\"n\":1.0}"`))
	require.NoError(t, err)
	require.Equal(t, `{"n":1.0,"title":"Bug"}`, string(canonical)) //nolint:testifylint // lexeme must stay 1.0
}

func Test_CanonicalTemplateData_rejects_trailing_JSON(t *testing.T) {
	_, err := CanonicalTemplateData([]byte(`{"a":1}{"b":2}`))
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CanonicalTemplateData_rejects_trailing_delimiters(t *testing.T) {
	for _, raw := range []string{`{"a":1}]`, `{"a":1}}`, `{"a":1}]garbage`} {
		_, err := CanonicalTemplateData([]byte(raw))
		require.ErrorIs(t, err, ErrWriteInvalid, raw)
	}
}

func Test_CanonicalTemplateData_allows_trailing_whitespace(t *testing.T) {
	canonical, err := CanonicalTemplateData([]byte("{\"a\":1} \n\t"))
	require.NoError(t, err)
	require.JSONEq(t, `{"a":1}`, string(canonical))
}

func Test_mustJSON_panics_when_value_cannot_marshal(t *testing.T) {
	require.Panics(t, func() {
		mustJSON(make(chan int))
	})
}

func Test_CanonicalTemplateData_rejects_non_objects(t *testing.T) {
	for _, raw := range []string{``, `null`, `[]`, `123`, `"hello"`, `"[]"`, `"null"`, `true`} {
		_, err := CanonicalTemplateData([]byte(raw))
		require.ErrorIs(t, err, ErrWriteInvalid, raw)
	}
}
