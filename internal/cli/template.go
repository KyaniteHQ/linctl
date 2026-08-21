package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addTemplateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	spec := readListGetSpec[client.TemplateList, client.TemplateSummary]{
		Use:       "template",
		Short:     "Read and write Linear templates",
		ListShort: "List visible Linear templates",
		LimitHelp: "maximum templates to print",
		GetUse:    "get TEMPLATE_ID",
		GetShort:  "Get one template by id",
		LoadList:  clientList(client.ListTemplates),
		LoadGet:   clientGet(client.GetTemplateByID),
		WriteItem: writeTemplate,
	}
	templateCommand := addReadListGetCommand(ctx, root, options, spec)
	addTemplateContentCommand(ctx, templateCommand, options)
	addTemplateCreateCommand(ctx, templateCommand, options)
	addTemplateUpdateCommand(ctx, templateCommand, options)
}

func writeTemplate(command *cobra.Command, options *rootOptions, template client.TemplateSummary) error {
	return writeItem(command, options, template, template.ID,
		func(command *cobra.Command, _ *rootOptions, template client.TemplateSummary) error {
			scope := "organization"
			if template.TeamKey != "" {
				scope = "team " + template.TeamKey
			}

			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s [%s] %s",
				template.ID,
				template.Name,
				template.Type,
				scope,
			)
		})
}
