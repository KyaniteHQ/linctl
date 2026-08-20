package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

type authReadinessRequest struct {
	AccessToken    string
	TokenActor     string
	TokenScopes    []string
	ExpectedTarget config.Target
	ExpectedActor  string
	RequiredScopes []string
	Timeout        time.Duration
}

type authReadinessReport struct {
	Actor  string                `json:"actor"`
	Target client.ResolvedTarget `json:"target"`
}

func requireAuthReadiness(ctx context.Context, request authReadinessRequest) (authReadinessReport, error) {
	readiness, err := checkAuthReadiness(ctx, request)
	if err != nil {
		return authReadinessReport{}, mapAuthReadinessError(err)
	}
	if err := requireScopes(request.TokenScopes, request.RequiredScopes); err != nil {
		return authReadinessReport{}, err
	}
	if request.ExpectedActor != "" && readiness.Actor != request.ExpectedActor {
		return authReadinessReport{}, auth.NewError(
			auth.ErrorCodeActorMismatch,
			fmt.Sprintf("expected actor %q but resolved %q", request.ExpectedActor, readiness.Actor),
		)
	}

	return readiness, nil
}

func requireLoggedAuthReadiness(
	ctx context.Context,
	logger *slog.Logger,
	request authReadinessRequest,
) (authReadinessReport, error) {
	if logger == nil {
		logger = discardLogger
	}
	logger.DebugContext(
		ctx,
		"auth readiness check started",
		"expected_actor", request.ExpectedActor,
		"required_scopes", strings.Join(request.RequiredScopes, ","),
		"org", request.ExpectedTarget.OrgID,
		"team_key", request.ExpectedTarget.TeamKey,
		"team_id", request.ExpectedTarget.TeamID,
		"project", request.ExpectedTarget.ProjectID,
	)

	readiness, err := requireAuthReadiness(ctx, request)
	if err != nil {
		logger.DebugContext(
			ctx,
			"auth readiness check failed",
			"expected_actor", request.ExpectedActor,
			"error_code", errorCode(err),
		)

		return authReadinessReport{}, err
	}

	projectID := ""
	if readiness.Target.Project != nil {
		projectID = readiness.Target.Project.ID
	}
	logger.DebugContext(
		ctx,
		"auth readiness check succeeded",
		"actor", readiness.Actor,
		"org", readiness.Target.Org.ID,
		"team_key", readiness.Target.Team.Key,
		"team_id", readiness.Target.Team.ID,
		"project", projectID,
		"confirmed", readiness.Target.Confirmed,
	)

	return readiness, nil
}

func defaultCheckAuthReadiness(ctx context.Context, request authReadinessRequest) (authReadinessReport, error) {
	if request.AccessToken == "" {
		return authReadinessReport{}, auth.NewError(auth.ErrorCodeNotConfigured, "missing OAuth access token")
	}
	graphqlClient := newAuthReadinessGraphQLClient(request.AccessToken, request.Timeout)
	target, err := client.ResolveTarget(ctx, graphqlClient, request.ExpectedTarget)
	if err != nil {
		return authReadinessReport{}, err
	}

	return authReadinessReport{Actor: actorFromViewer(target.Viewer), Target: target}, nil
}

func actorFromViewer(viewer client.TargetViewer) string {
	if viewer.App {
		return appActor
	}

	return userActor
}

func mapAuthReadinessError(err error) error {
	var authErr *auth.AuthError
	var tokenErr *auth.TokenEndpointError
	switch {
	case errors.As(err, &authErr):
		return err
	case errors.As(err, &tokenErr):
		return err
	case errors.Is(err, client.ErrTargetNotConfigured):
		return auth.WrapError(
			auth.ErrorCodeTargetNotConfigured,
			"no pinned target is configured: set org_id, team_key, and team_id in .linctl.toml",
			err,
		)
	case errors.Is(err, client.ErrTargetMismatch):
		return auth.WrapError(
			auth.ErrorCodeTargetMismatch,
			"OAuth authorization does not match the pinned target",
			err,
		)
	default:
		return err
	}
}

func requireScopes(actual []string, required []string) error {
	missing := missingScopes(actual, required)
	if len(missing) == 0 {
		return nil
	}

	return auth.NewError(
		auth.ErrorCodeMissingScope,
		"missing OAuth scopes: "+strings.Join(missing, ",")+
			"; run linctl auth configure --scopes "+strings.Join(required, ",")+
			" then linctl auth app or linctl auth login",
	)
}

func missingScopes(actual []string, required []string) []string {
	actualSet := map[string]bool{}
	for _, scope := range actual {
		actualSet[scope] = true
	}
	missing := []string{}
	for _, scope := range required {
		if !actualSet[scope] {
			missing = append(missing, scope)
		}
	}

	return missing
}

func requiredScopes(app auth.AppConfig) []string {
	if len(app.Scopes) > 0 {
		return slices.Clone(app.Scopes)
	}

	return slices.Clone(defaultOAuthScopes)
}

func normalizedScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return nil
	}
	normalized := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		for _, part := range auth.SplitScopes(scope) {
			part = strings.TrimSpace(part)
			if part != "" && !slices.Contains(normalized, part) {
				normalized = append(normalized, part)
			}
		}
	}

	return normalized
}
