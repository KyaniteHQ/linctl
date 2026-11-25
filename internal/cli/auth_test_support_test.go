package cli

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
	"github.com/KyaniteHQ/linctl/internal/oauth"
)

func useAuthPaths(t *testing.T, paths auth.Paths) func() {
	t.Helper()
	original := authDefaultPaths
	authDefaultPaths = func(auth.Env) (auth.Paths, error) {
		return paths, nil
	}
	return func() {
		authDefaultPaths = original
	}
}

func useAuthPathsError(t *testing.T, err error) func() {
	t.Helper()
	original := authDefaultPaths
	authDefaultPaths = func(auth.Env) (auth.Paths, error) {
		return auth.Paths{}, err
	}
	return func() {
		authDefaultPaths = original
	}
}

func useAuthCommandHooks(
	t *testing.T,
	paths auth.Paths,
	oauthClient *fakeOAuthTokenClient,
	readiness *fakeAuthReadinessChecker,
) func() {
	t.Helper()
	restorePaths := useAuthPaths(t, paths)
	originalOAuthClient := newAuthOAuthClient
	originalReadiness := checkAuthReadiness
	newAuthOAuthClient = func(time.Duration) authOAuthClient {
		return oauthClient
	}
	checkAuthReadiness = readiness.check
	return func() {
		checkAuthReadiness = originalReadiness
		newAuthOAuthClient = originalOAuthClient
		restorePaths()
	}
}

type fakeOAuthTokenClient struct {
	grant                    auth.TokenState
	err                      error
	clientCredentialsCalls   int
	authorizationCodeCalls   int
	authorizationCodeRequest oauth.AuthorizationCodeRequest
	refreshTokenCalls        int
	refreshTokenRequest      oauth.RefreshTokenRequest
	beforeRefresh            func()
	revokeTokenCalls         int
	revokeTokenRequests      []oauth.RevocationRequest
	revokeTokenErr           error
	beforeRevoke             func()
}

func (client *fakeOAuthTokenClient) ClientCredentials(
	_ context.Context,
	_ oauth.ClientCredentialsRequest,
) (auth.TokenState, error) {
	client.clientCredentialsCalls++
	if client.err != nil {
		return auth.TokenState{}, client.err
	}

	return client.grant, nil
}

func (client *fakeOAuthTokenClient) RefreshToken(
	_ context.Context,
	request oauth.RefreshTokenRequest,
) (auth.TokenState, error) {
	client.refreshTokenCalls++
	client.refreshTokenRequest = request
	if client.beforeRefresh != nil {
		client.beforeRefresh()
	}
	if client.err != nil {
		return auth.TokenState{}, client.err
	}

	return client.grant, nil
}

func (client *fakeOAuthTokenClient) ExchangeAuthorizationCode(
	_ context.Context,
	request oauth.AuthorizationCodeRequest,
) (auth.TokenState, error) {
	client.authorizationCodeCalls++
	client.authorizationCodeRequest = request
	if client.err != nil {
		return auth.TokenState{}, client.err
	}

	return client.grant, nil
}

func (client *fakeOAuthTokenClient) RevokeToken(_ context.Context, request oauth.RevocationRequest) error {
	client.revokeTokenCalls++
	client.revokeTokenRequests = append(client.revokeTokenRequests, request)
	if client.beforeRevoke != nil {
		client.beforeRevoke()
		client.beforeRevoke = nil
	}
	if client.revokeTokenErr != nil {
		return client.revokeTokenErr
	}

	return nil
}

type fakeAuthReadinessChecker struct {
	report         authReadinessReport
	err            error
	accessToken    string
	expectedActor  string
	requiredScopes []string
	rejectedTokens []string
	beforeReturn   func()
}

func (checker *fakeAuthReadinessChecker) check(
	_ context.Context,
	request authReadinessRequest,
) (authReadinessReport, error) {
	checker.accessToken = request.AccessToken
	checker.expectedActor = request.ExpectedActor
	checker.requiredScopes = request.RequiredScopes
	if checker.beforeReturn != nil {
		checker.beforeReturn()
	}
	if checker.err != nil {
		return authReadinessReport{}, checker.err
	}
	if slices.Contains(checker.rejectedTokens, request.AccessToken) {
		return authReadinessReport{}, fmt.Errorf("resolve viewer: %w: graphql http status 401", client.ErrAuthFailed)
	}

	return checker.report, nil
}

func readyAuthReport(actor string) authReadinessReport {
	return authReadinessReport{
		Actor: actor,
		Target: client.ResolvedTarget{
			Org:       client.TargetOrg{ID: "org-id"},
			Team:      client.TargetTeam{ID: "team-id", Key: "LIT"},
			Expected:  config.Target{OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id"},
			Resolved:  config.Target{OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id"},
			Confirmed: true,
		},
	}
}
