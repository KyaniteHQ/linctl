package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// guardedWriteSpec describes one guarded-write command in the shared write
// pipeline. Configure registers flags and completions, Run performs the
// guarded client call with the built runtime, and Write renders the result.
// Pre-call steps (stdin body resolution, normalization, template lookup)
// belong inside Run, which receives the command and raw args.
type guardedWriteSpec[T any] struct {
	Use       string
	Short     string
	Args      cobra.PositionalArgs
	Configure func(*cobra.Command)
	Run       func(context.Context, *cobra.Command, commandRuntime, []string) (T, error)
	Write     func(*cobra.Command, *rootOptions, T) error
}

// addGuardedWriteCommand registers a guarded-write command: the shared RunE
// builds the command runtime, delegates to spec.Run, and renders through
// spec.Write. Registration stamps the write safety class explicitly so the
// static command inventory never infers it from prose.
func addGuardedWriteCommand[T any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	spec guardedWriteSpec[T],
) {
	command := &cobra.Command{
		Use:   spec.Use,
		Short: spec.Short,
		Args:  spec.Args,
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			item, err := spec.Run(ctx, command, runtime, args)
			if err != nil {
				return err
			}

			return spec.Write(command, options, item)
		},
	}
	if spec.Configure != nil {
		spec.Configure(command)
	}
	addCommandWithSafety(root, CommandSafetyWrite, command)
}
