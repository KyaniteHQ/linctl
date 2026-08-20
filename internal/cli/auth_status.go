package cli

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/render"
)

type authConfigReport struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURI  string   `json:"redirect_uri,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

type authTokenReport struct {
	Status    string     `json:"status"`
	Type      string     `json:"type,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Scopes    []string   `json:"scopes,omitempty"`
}

type authTargetStatusReport struct {
	Status   string            `json:"status"`
	Expected map[string]string `json:"expected,omitempty"`
	Resolved map[string]string `json:"resolved,omitempty"`
}

type authStatusReport struct {
	App       authConfigReport       `json:"app"`
	Token     authTokenReport        `json:"token"`
	Actor     string                 `json:"actor,omitempty"`
	Scopes    []string               `json:"scopes,omitempty"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
	TokenType string                 `json:"token_type,omitempty"`
	Target    authTargetStatusReport `json:"target"`
}

type authLogoutReport struct {
	Token            string   `json:"token"`
	App              string   `json:"app"`
	Revoked          []string `json:"revoked,omitempty"`
	RevocationFailed bool     `json:"revocation_failed,omitempty"`
}

var afterAuthStatusTokenAcquired func()

func writeCurrentOrRefreshedAuthStatus(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	authContext authCommandContext,
	app auth.AppConfig,
	token auth.TokenState,
) error {
	token, readiness, err := currentOrRefreshedAuthToken(ctx, authContext, app, token, options.timeout)
	if err != nil {
		return err
	}
	if afterAuthStatusTokenAcquired != nil {
		afterAuthStatusTokenAcquired()
	}

	return writeAuthStatus(command, options, newAuthStatusReport(app, token, readiness))
}

func newAuthStatusReport(
	app auth.AppConfig,
	token auth.TokenState,
	readiness authReadinessReport,
) authStatusReport {
	return authStatusReport{
		App: redactedAppConfigReport(app),
		Token: authTokenReport{
			Status:    presence(token.AccessToken),
			Type:      token.TokenType,
			ExpiresAt: token.ExpiresAt,
			Scopes:    slices.Clone(token.Scopes),
		},
		Actor:     cmp.Or(readiness.Actor, token.Actor),
		Scopes:    slices.Clone(token.Scopes),
		ExpiresAt: token.ExpiresAt,
		TokenType: token.TokenType,
		Target: authTargetStatusReport{
			Status:   "ready",
			Expected: targetMap(readiness.Target.Expected),
			Resolved: targetMap(readiness.Target.Resolved),
		},
	}
}

func redactedAppConfigReport(app auth.AppConfig) authConfigReport {
	return authConfigReport{
		ClientID:     presence(app.ClientID),
		ClientSecret: presence(app.ClientSecret),
		RedirectURI:  app.RedirectURI,
		Scopes:       slices.Clone(app.Scopes),
	}
}

func presence(value string) string {
	if value == "" {
		return "missing"
	}

	return "set"
}

func writeAuthStatus(command *cobra.Command, options *rootOptions, status authStatusReport) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, status)
	}

	return writeAuthStatusHuman(command, status)
}

func writeAuthStatusHuman(command *cobra.Command, status authStatusReport) error {
	return render.WriteLine(
		command.OutOrStdout(),
		"auth %s actor %s scopes %s target %s",
		status.Token.Status,
		defaultString(status.Actor, "unknown"),
		strings.Join(status.Scopes, ","),
		status.Target.Status,
	)
}
