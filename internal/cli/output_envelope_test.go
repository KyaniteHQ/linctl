package cli

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_errorCode_maps_sentinels_and_fallbacks(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code string
	}{
		{name: "target mismatch", err: fmt.Errorf("%w: x", client.ErrTargetMismatch), code: "TARGET_MISMATCH"},
		{name: "cross-organization relation", err: fmt.Errorf("%w: x", client.ErrCrossOrganizationRelation), code: "CROSS_ORGANIZATION_RELATION"},
		{name: "state mismatch", err: fmt.Errorf("%w: x", client.ErrStateMismatch), code: "STATE_MISMATCH"},
		{name: "transition denied", err: fmt.Errorf("%w: x", client.ErrTransitionDenied), code: "TRANSITION_DENIED"},
		{name: "write conflict", err: fmt.Errorf("%w: x", client.ErrWriteConflict), code: "CONFLICT"},
		{name: "target not configured", err: fmt.Errorf("%w: x", client.ErrTargetNotConfigured), code: "TARGET_NOT_CONFIGURED"},
		{name: "rate limited", err: fmt.Errorf("%w: x", client.ErrRateLimited), code: "RATE_LIMITED"},
		{name: "mutation failed", err: fmt.Errorf("%w: x", client.ErrMutationFailed), code: "MUTATION_FAILED"},
		{name: "invalid write", err: fmt.Errorf("%w: x", client.ErrWriteInvalid), code: "INVALID_WRITE"},
		{name: "pin exists", err: fmt.Errorf("%w: x", config.ErrPinExists), code: "INVALID_WRITE"},
		{name: "graphql", err: fmt.Errorf("%w: x", client.ErrGraphQL), code: "GRAPHQL_ERROR"},
		{name: "not found", err: fmt.Errorf("get issue LIT-1: %w", client.ErrNotFound), code: "NOT_FOUND"},
		{name: "auth readiness", err: auth.NewError(auth.ErrorCodeTargetMismatch, "wrong target"), code: "AUTH_TARGET_MISMATCH"},
		{name: "token endpoint", err: auth.NewTokenEndpointError(auth.ErrorCodeMissingScope, 400, "invalid_scope"), code: "MISSING_SCOPE"},
		{name: "fallback", err: errors.New("something unexpected"), code: "INTERNAL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.code, errorCode(tt.err))
		})
	}
}
