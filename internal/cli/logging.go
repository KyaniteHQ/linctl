package cli

import (
	"errors"
	"io"
	"log/slog"
	"strings"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// discardLogger is the fallback for runtimes built without a logger (test
// fakes), so logging call-sites never need a nil check.
var discardLogger = slog.New(slog.DiscardHandler)

// log returns the runtime logger, or a discarding logger when one was not set.
func (runtime commandRuntime) log() *slog.Logger {
	if runtime.logger == nil {
		return discardLogger
	}

	return runtime.logger
}

// logTargetResolution records the outcome of a target resolution at debug
// level. The mismatch attribute flags a Target Mismatch refusal versus an
// unrelated resolve error; the underlying error is already surfaced to the
// user by the command, so this stays at debug rather than warn.
func logTargetResolution(logger *slog.Logger, target client.ResolvedTarget, err error) {
	if err != nil {
		logger.Debug(
			"target unresolved",
			"mismatch", errors.Is(err, client.ErrTargetMismatch),
			"error", err.Error(),
		)

		return
	}

	logger.Debug(
		"target resolved",
		"org", target.Org.ID,
		"team_id", target.Team.ID,
		"team_key", target.Team.Key,
		"confirmed", target.Confirmed,
	)
}

// newDiagnosticLogger builds the stderr diagnostic logger. Without --debug it
// emits warnings and above (for example an insecure-config notice); with
// --debug it emits the full debug stream. LINCTL_DEBUG_JSON=1 selects JSON
// output over the default text format.
func newDiagnosticLogger(debug bool, jsonOutput bool, destination io.Writer) *slog.Logger {
	level := slog.LevelWarn
	if debug {
		level = slog.LevelDebug
	}

	handlerOptions := &slog.HandlerOptions{Level: level}
	if jsonOutput {
		return slog.New(slog.NewJSONHandler(destination, handlerOptions))
	}

	return slog.New(slog.NewTextHandler(destination, handlerOptions))
}

// transportDiagnosticWriter forwards the transport's line diagnostics into the
// slog stream at debug level so every diagnostic shares one sink and format.
type transportDiagnosticWriter struct {
	logger *slog.Logger
}

func (writer transportDiagnosticWriter) Write(payload []byte) (int, error) {
	writer.logger.Debug("transport", "detail", strings.TrimRight(string(payload), "\n"))

	return len(payload), nil
}

// newTransportDiagnosticWriter wires transport diagnostics into the logger only
// under --debug; otherwise it returns nil so the transport stays silent.
func newTransportDiagnosticWriter(logger *slog.Logger, debug bool) io.Writer {
	if !debug {
		return nil
	}

	return transportDiagnosticWriter{logger: logger}
}
