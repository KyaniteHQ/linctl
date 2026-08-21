package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
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
	return writeItemLine(
		command, options, template, template.ID,
		"%s %s [%s] %s",
		template.ID, template.Name, template.Type, templateScopeLabel(template.TeamKey),
	)
}

func templateScopeLabel(teamKey string) string {
	scope := "organization"
	if teamKey != "" {
		scope = "team " + teamKey
	}

	return scope
}
