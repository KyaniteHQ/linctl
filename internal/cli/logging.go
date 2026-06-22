package cli

import (
	"io"
	"log/slog"
	"strings"
)

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
