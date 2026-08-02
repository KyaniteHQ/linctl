package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newGroupCommand builds a parent command that only routes to subcommands.
// A bare invocation prints help; an unknown subcommand token is an error so
// agent pipelines fail loudly instead of parsing help output. A group command
// changes nothing in Linear, so it carries the read safety class explicitly
// and never depends on the wording of its Short text.
func newGroupCommand(use string, short string) *cobra.Command {
	command := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) == 0 {
				return command.Help()
			}

			return fmt.Errorf("unknown command %q for %q", args[0], command.CommandPath())
		},
	}
	annotateCommand(command, commandSafetyAnnotation, string(CommandSafetyRead))

	return command
}
