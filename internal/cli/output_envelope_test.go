package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_errorCode_maps_sentinels_and_fallbacks(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code string
	}{
		{name: "target mismatch", err: fmt.Errorf("%w: x", client.ErrTargetMismatch), code: "TARGET_MISMATCH"},
		{name: "rate limited", err: fmt.Errorf("%w: x", client.ErrRateLimited), code: "RATE_LIMITED"},
		{name: "mutation failed", err: fmt.Errorf("%w: x", client.ErrMutationFailed), code: "MUTATION_FAILED"},
		{name: "invalid write", err: fmt.Errorf("%w: x", client.ErrWriteInvalid), code: "INVALID_WRITE"},
		{name: "graphql", err: fmt.Errorf("%w: x", client.ErrGraphQL), code: "GRAPHQL_ERROR"},
		{name: "not found", err: errors.New("get issue LIT-1: not found"), code: "NOT_FOUND"},
		{name: "fallback", err: errors.New("something unexpected"), code: "INTERNAL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.code, errorCode(tt.err))
		})
	}
}

func Test_isNotFoundError_matches_suffix(t *testing.T) {
	require.True(t, isNotFoundError(errors.New("get issue by vcs branch x: not found")))
	require.False(t, isNotFoundError(errors.New("some other failure")))
}

func Test_writeErrorEnvelope_emits_code_and_message(t *testing.T) {
	var buf bytes.Buffer

	err := writeErrorEnvelope(&buf, fmt.Errorf("%w: team differs", client.ErrTargetMismatch))

	require.NoError(t, err)
	var envelope errorEnvelope
	require.NoError(t, json.Unmarshal(buf.Bytes(), &envelope))
	require.Equal(t, "TARGET_MISMATCH", envelope.ErrorCode)
	require.Contains(t, envelope.Message, "team differs")
}

func Test_writeErrorEnvelope_returns_write_error(t *testing.T) {
	err := writeErrorEnvelope(commandFailingWriter{}, errors.New("boom"))

	require.Error(t, err)
}
