package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addAuditEntryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "audit-entry",
		Short: "Read Linear audit entry catalogs",
	}
	command.AddCommand(&cobra.Command{
		Use:   "types",
		Short: "List Linear audit entry types",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			types, err := client.ListAuditEntryTypes(ctx, runtime.graphqlClient)
			if err != nil {
				return err
			}

			return writeAuditEntryTypes(command, options, types)
		},
	})
	root.AddCommand(command)
}

func writeAuditEntryTypes(
	command *cobra.Command,
	options *rootOptions,
	types client.AuditEntryTypeList,
) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, types)
	}

	for _, entryType := range types.AuditEntryTypes {
		if err := writeAuditEntryType(command, options, entryType); err != nil {
			return err
		}
	}
	return nil
}

func writeAuditEntryType(
	command *cobra.Command,
	options *rootOptions,
	entryType client.AuditEntryTypeSummary,
) error {
	if wrote, err := writeIDOnly(command, options, entryType.Type); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, entryType)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s", entryType.Type, entryType.Description)
}
