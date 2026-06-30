package oauth

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClientUsesDefaultConfigValues(t *testing.T) {
	client := NewClient(ClientConfig{})

	require.Equal(t, defaultTokenEndpoint, client.endpoint)
	require.Equal(t, defaultRevocationEndpoint, client.revocationEndpoint)
	require.Same(t, http.DefaultClient, client.httpClient)
	require.NotNil(t, client.now)
	require.WithinDuration(t, time.Now(), client.now(), time.Second)
}

func TestScopeListUnmarshalJSONAcceptsNull(t *testing.T) {
	var response tokenEndpointResponse

	err := json.Unmarshal([]byte(`{"scope":null}`), &response)

	require.NoError(t, err)
	require.Nil(t, []string(response.Scopes))
}

func TestScopeListUnmarshalJSONRejectsUnsupportedValue(t *testing.T) {
	var response tokenEndpointResponse

	err := json.Unmarshal([]byte(`{"scope":42}`), &response)

	require.ErrorContains(t, err, "scope must be a string or array")
}
