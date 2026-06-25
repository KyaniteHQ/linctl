package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addRateLimitCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "rate-limit",
		Short: "Read Linear rate-limit status",
	}
	command.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Print the authenticated Linear rate-limit status",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			status, err := client.GetRateLimitStatus(ctx, runtime.graphqlClient)
			if err != nil {
				return err
			}

			return writeRateLimitStatus(command, options, status)
		},
	})
	root.AddCommand(command)
}

func writeRateLimitStatus(command *cobra.Command, options *rootOptions, status client.RateLimitStatus) error {
	return writeItemNoID(command, options, status,
		func(command *cobra.Command, _ *rootOptions, status client.RateLimitStatus) error {
			if err := render.WriteLine(
				command.OutOrStdout(),
				"%s %s",
				status.Kind,
				emptyDash(status.Identifier),
			); err != nil {
				return err
			}
			for _, limit := range status.Limits {
				if err := render.WriteLine(
					command.OutOrStdout(),
					"%s remaining %.0f/%.0f reset %.0f",
					limit.Type,
					limit.RemainingAmount,
					limit.AllowedAmount,
					limit.Reset,
				); err != nil {
					return err
				}
			}

			return nil
		})
}
