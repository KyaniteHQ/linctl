package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TokenState_formatting_never_contains_secret_values(t *testing.T) {
	t.Parallel()
	token := TokenState{
		AccessToken:  "raw-access",
		RefreshToken: "raw-refresh",
		Actor:        "raw-actor",
	}
	app := AppConfig{
		ClientID:     "raw-client-id",
		ClientSecret: "raw-client-secret",
	}

	outputs := make([]string, 0, 10)
	outputs = append(
		outputs,
		fmt.Sprintf("%v", token),
		fmt.Sprintf("%+v", token),
		fmt.Sprintf("%#v", token),
		fmt.Sprintf("%s", token), //nolint:staticcheck // exercising the Stringer verb explicitly.
		fmt.Sprintf("%v", app),
		fmt.Sprintf("%+v", app),
		fmt.Sprintf("%#v", app),
		fmt.Sprintf("%s", app), //nolint:staticcheck // exercising the Stringer verb explicitly.
	)

	tokenJSON, err := json.Marshal(token)
	require.NoError(t, err)
	appJSON, err := json.Marshal(app)
	require.NoError(t, err)
	outputs = append(outputs, string(tokenJSON), string(appJSON))

	for _, output := range outputs {
		require.NotContains(t, output, "raw-access")
		require.NotContains(t, output, "raw-refresh")
		require.NotContains(t, output, "raw-actor")
		require.NotContains(t, output, "raw-client-id")
		require.NotContains(t, output, "raw-client-secret")
	}

	require.Contains(t, fmt.Sprintf("%+v", token), "set")
	require.Contains(t, string(tokenJSON), `"access_token":"set"`)
	require.Contains(t, string(appJSON), `"client_secret":"set"`)
}

func Test_Store_persistence_bytes_unchanged_by_redaction(t *testing.T) {
	t.Parallel()
	paths := testPaths(t)
	store := NewStore(paths)
	want := State{
		App: AppConfig{
			ClientID:     "client-id",
			ClientSecret: "raw-client-secret",
		},
		Token: TokenState{
			AccessToken:  "raw-access",
			RefreshToken: "raw-refresh",
		},
	}

	require.NoError(t, store.Save(context.Background(), want))

	appBytes, err := os.ReadFile(paths.AppConfigPath)
	require.NoError(t, err)
	tokenBytes, err := os.ReadFile(paths.TokenPath)
	require.NoError(t, err)

	// require.Equal (not JSONEq) is intentional: this is a byte-for-byte tripwire
	// for indentation and key order, not a semantic JSON comparison.
	appWant := "{\n" +
		"  \"app\": {\n" +
		"    \"client_id\": \"client-id\",\n" +
		"    \"client_secret\": \"raw-client-secret\"\n" +
		"  }\n" +
		"}\n"
	tokenWant := "{\n" +
		"  \"token\": {\n" +
		"    \"access_token\": \"raw-access\",\n" +
		"    \"refresh_token\": \"raw-refresh\"\n" +
		"  }\n" +
		"}\n"
	require.Equal(t, appWant, string(appBytes))     //nolint:testifylint // exact byte comparison intentional, see comment above.
	require.Equal(t, tokenWant, string(tokenBytes)) //nolint:testifylint // exact byte comparison intentional, see comment above.

	got, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, want, got)
}
