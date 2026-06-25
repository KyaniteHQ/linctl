//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addTemplateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.TemplateList, client.TemplateSummary]{
		Use:           "template",
		Short:         "Read Linear templates",
		ListShort:     "List visible Linear templates",
		LimitHelp:     "maximum templates to print",
		GetUse:        "get TEMPLATE_ID",
		GetShort:      "Get one template by id",
		LoadList:      loadTemplateList,
		PageWithItems: templatePageWithItems,
		LoadGet:       loadTemplate,
		WriteItem:     writeTemplate,
	})
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

func loadTemplateList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TemplateList, []client.TemplateSummary, error) {
	templates, err := client.ListTemplates(ctx, runtime.graphqlClient, limit)
	return templates, templates.Templates, err
}

func loadTemplate(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.TemplateSummary, error) {
	return client.GetTemplateByID(ctx, runtime.graphqlClient, id)
}

func templatePageWithItems(page client.TemplateList, templates []client.TemplateSummary) client.TemplateList {
	page.Templates = templates
	return page
}
