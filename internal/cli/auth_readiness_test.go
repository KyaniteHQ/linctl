package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
)

type linearGraphQLRoundTripper struct {
	payloads       map[string]string
	urls           []string
	authorizations []string
}

func (linear *linearGraphQLRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	var body struct {
		OperationName string `json:"operationName"`
	}
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		return nil, err
	}
	linear.urls = append(linear.urls, request.URL.String())
	linear.authorizations = append(linear.authorizations, request.Header.Get("Authorization"))
	response := `{"errors":[{"message":"viewer unavailable"}]}`
	if payload, ok := linear.payloads[body.OperationName]; ok {
		response = `{"data":` + payload + `}`
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(response)),
		Request:    request,
	}, nil
}

func useLinearGraphQL(t *testing.T, viewerIsApp bool, viewerAvailable bool) *linearGraphQLRoundTripper {
	t.Helper()
	viewer, _ := commandFlowBasePayload("Viewer")
	if viewerIsApp {
		viewer = strings.Replace(viewer, `"app":false`, `"app":true`, 1)
	}
	teams, _ := commandFlowBasePayload("Teams")
	project, _ := commandFlowBasePayload("TargetProject")
	linear := &linearGraphQLRoundTripper{payloads: map[string]string{"Teams": teams, "TargetProject": project}}
	if viewerAvailable {
		linear.payloads["Viewer"] = viewer
	}
	original := http.DefaultTransport
	http.DefaultTransport = linear
	t.Cleanup(func() {
		http.DefaultTransport = original
	})

	return linear
}

func Test_AuthStatus_checks_readiness_against_linear(t *testing.T) {
	pinned := []string{"--json", "--org", "org-id", "--team", "LIT", "--team-id", "team-id"}
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	appToken := auth.TokenState{
		AccessToken: "saved-access-token",
		TokenType:   "Bearer",
		Scopes:      []string{"read"},
		ExpiresAt:   &expiresAt,
		Actor:       "app",
		GrantType:   authGrantClientCredentials,
	}
	userToken := appToken
	userToken.Actor = "user"
	userToken.RefreshToken = "saved-refresh-token"
	userToken.GrantType = authGrantAuthorizationCode
	unscopedToken := appToken
	unscopedToken.Scopes = nil
	tests := []struct {
		name          string
		injectedToken string
		token         auth.TokenState
		viewerIsApp   bool
		viewerDown    bool
		args          []string
		wantActor     string
		wantError     string
		wantErrorCode string
	}{
		{name: "app token resolved as the app actor", token: appToken, viewerIsApp: true, args: pinned, wantActor: "app"},
		{name: "user token resolved as the user actor", token: userToken, args: pinned, wantActor: "user"},
		{
			name:          "injected token reports the actor linear resolves",
			injectedToken: "injected-access-token",
			args:          pinned,
			wantActor:     "user",
		},
		{
			name:          "user token resolved as an app actor",
			token:         userToken,
			viewerIsApp:   true,
			args:          pinned,
			wantErrorCode: "AUTH_ACTOR_MISMATCH",
		},
		{
			name:          "token missing a configured scope",
			token:         unscopedToken,
			viewerIsApp:   true,
			args:          pinned,
			wantErrorCode: "MISSING_SCOPE",
		},
		{
			name:          "no pinned target",
			token:         appToken,
			viewerIsApp:   true,
			args:          []string{"--json"},
			wantErrorCode: "AUTH_TARGET_NOT_CONFIGURED",
		},
		{
			name:          "saved state without an access token",
			token:         auth.TokenState{RefreshToken: "saved-refresh-token", Actor: "user"},
			args:          pinned,
			wantErrorCode: "AUTH_NOT_CONFIGURED",
		},
		{
			name:       "linear answers the viewer query with an error",
			token:      appToken,
			viewerDown: true,
			args:       pinned,
			wantError:  "viewer unavailable",
		},
		{name: "quiet", token: appToken, viewerIsApp: true, args: slices.Concat([]string{"--quiet"}, pinned[1:])},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			if test.injectedToken != "" {
				t.Setenv("LINCTL_OAUTH_ACCESS_TOKEN", test.injectedToken)
			}
			paths := cliAuthTestPaths(t)
			require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
				App:   auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret", Scopes: []string{"read"}},
				Token: test.token,
			}))
			restorePaths := useAuthPaths(t, paths)
			defer restorePaths()
			useLinearGraphQL(t, test.viewerIsApp, !test.viewerDown)
			var stdout bytes.Buffer

			err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{},
				slices.Concat(test.args, []string{"auth", "status"}))

			if test.wantError != "" {
				require.ErrorContains(t, err, test.wantError)
				return
			}
			if test.wantErrorCode != "" {
				require.Error(t, err)
				require.Equal(t, test.wantErrorCode, errorCode(err))
				return
			}
			require.NoError(t, err)
			if test.wantActor == "" {
				require.Empty(t, stdout.String())
				return
			}
			var report authStatusReport
			require.NoError(t, json.Unmarshal(stdout.Bytes(), &report))
			require.Equal(t, test.wantActor, report.Actor)
			require.Equal(t, "ready", report.Target.Status)
			require.Equal(t, "team-id", report.Target.Resolved["team_id"])
		})
	}
}

func Test_AuthStatus_sends_the_saved_token_to_linear_and_logs_the_resolved_project(t *testing.T) {
	t.Chdir(t.TempDir())
	paths := cliAuthTestPaths(t)
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret", Scopes: []string{"read"}},
		Token: auth.TokenState{
			AccessToken: "saved-access-token",
			Scopes:      []string{"read"},
			ExpiresAt:   &expiresAt,
			Actor:       "app",
			GrantType:   authGrantClientCredentials,
		},
	}))
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	linear := useLinearGraphQL(t, true, true)
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
		"--debug",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"--project", "project-id",
		"auth",
		"status",
	})

	require.NoError(t, err)
	require.NotEmpty(t, linear.urls)
	for index, requestURL := range linear.urls {
		require.Equal(t, "https://api.linear.app/graphql", requestURL)
		require.Equal(t, "Bearer saved-access-token", linear.authorizations[index])
	}
	succeeded := ""
	for line := range strings.Lines(stderr.String()) {
		if strings.Contains(line, `msg="auth readiness check succeeded"`) {
			succeeded = line
		}
	}
	require.Contains(t, succeeded, "project=project-id")
	require.NotContains(t, stderr.String(), "saved-access-token")
	require.NotContains(t, stderr.String(), "client-secret")
}
