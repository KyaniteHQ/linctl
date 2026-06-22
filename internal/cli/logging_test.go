package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_newDiagnosticLogger_debug_emits_text_diagnostics(t *testing.T) {
	var buffer bytes.Buffer
	logger := newDiagnosticLogger(true, false, &buffer)

	logger.Debug("runtime ready", "team_key", "LIT")

	output := buffer.String()
	require.Contains(t, output, "level=DEBUG")
	require.Contains(t, output, "msg=\"runtime ready\"")
	require.Contains(t, output, "team_key=LIT")
}

func Test_newDiagnosticLogger_without_debug_drops_debug_but_keeps_warnings(t *testing.T) {
	var buffer bytes.Buffer
	logger := newDiagnosticLogger(false, false, &buffer)

	logger.Debug("debug line")
	require.Empty(t, buffer.String())

	logger.Warn("config world-readable")
	require.Contains(t, buffer.String(), "level=WARN")
	require.Contains(t, buffer.String(), "config world-readable")
}

func Test_newDiagnosticLogger_json_output_is_structured(t *testing.T) {
	var buffer bytes.Buffer
	logger := newDiagnosticLogger(true, true, &buffer)

	logger.Debug("graphql_retry", "attempt", 2)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buffer.Bytes(), &entry))
	require.Equal(t, "graphql_retry", entry["msg"])
	require.Equal(t, "DEBUG", entry["level"])
	require.InEpsilon(t, float64(2), entry["attempt"], 0)
}

func Test_newTransportDiagnosticWriter_forwards_trimmed_lines_under_debug(t *testing.T) {
	var buffer bytes.Buffer
	writer := newTransportDiagnosticWriter(newDiagnosticLogger(true, false, &buffer), true)
	require.NotNil(t, writer)

	line := []byte("graphql_response attempt=1 status=200\n")
	count, err := writer.Write(line)

	require.NoError(t, err)
	require.Equal(t, len(line), count)
	require.Contains(t, buffer.String(), `detail="graphql_response attempt=1 status=200"`)
}

func Test_newTransportDiagnosticWriter_is_nil_without_debug(t *testing.T) {
	require.Nil(t, newTransportDiagnosticWriter(newDiagnosticLogger(false, false, io.Discard), false))
}

func Test_commandRuntime_log_falls_back_to_discard_logger(t *testing.T) {
	require.Same(t, discardLogger, commandRuntime{}.log())

	logger := newDiagnosticLogger(true, false, io.Discard)
	require.Same(t, logger, commandRuntime{logger: logger}.log())
}

func Test_logTargetResolution_marks_target_mismatch(t *testing.T) {
	var buffer bytes.Buffer
	logger := newDiagnosticLogger(true, false, &buffer)

	logTargetResolution(logger, client.ResolvedTarget{}, fmt.Errorf("%w: org", client.ErrTargetMismatch))

	require.Contains(t, buffer.String(), "msg=\"target unresolved\"")
	require.Contains(t, buffer.String(), "mismatch=true")
}

func Test_logTargetResolution_marks_non_mismatch_errors(t *testing.T) {
	var buffer bytes.Buffer
	logger := newDiagnosticLogger(true, false, &buffer)

	logTargetResolution(logger, client.ResolvedTarget{}, errors.New("resolve viewer failed"))

	require.Contains(t, buffer.String(), "msg=\"target unresolved\"")
	require.Contains(t, buffer.String(), "mismatch=false")
}

func Test_logTargetResolution_records_resolved_target(t *testing.T) {
	var buffer bytes.Buffer
	logger := newDiagnosticLogger(true, false, &buffer)

	logTargetResolution(logger, client.ResolvedTarget{
		Org:       client.TargetOrg{ID: "org-id"},
		Team:      client.TargetTeam{ID: "team-id", Key: "LIT"},
		Confirmed: true,
	}, nil)

	output := buffer.String()
	require.Contains(t, output, "msg=\"target resolved\"")
	require.Contains(t, output, "team_key=LIT")
	require.Contains(t, output, "confirmed=true")
}
