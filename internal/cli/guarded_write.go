package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// runGuardedWrite builds the command runtime, runs the guarded write through it,
// and renders the resulting summary. It is the shared flow for the simple
// guarded-write commands whose only per-entity differences are the write call
// and the renderer; richer write surfaces (issue, project-update) depend on a
// Command Port instead so their request-assembly logic is testable in isolation.
func runGuardedWrite[T any](
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	write func(commandRuntime) (T, error),
	render func(*cobra.Command, *rootOptions, T) error,
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	result, err := write(runtime)
	if err != nil {
		return err
	}

	return render(command, options, result)
}
