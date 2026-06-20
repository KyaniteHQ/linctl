//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCustomViewCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.CustomViewList, client.CustomViewSummary]{
		Use:           "custom-view",
		Short:         "Read Linear custom views",
		ListShort:     "List visible custom views",
		LimitHelp:     "maximum custom views to return",
		GetUse:        "get CUSTOM_VIEW_ID",
		GetShort:      "Get one custom view by id or slug",
		LoadList:      loadCustomViewList,
		PageWithItems: customViewPageWithItems,
		LoadGet:       loadCustomView,
		WriteItem:     writeCustomView,
	})
}

func writeCustomView(
	command *cobra.Command,
	options *rootOptions,
	view client.CustomViewSummary,
) error {
	if wrote, err := writeIDOnly(command, options, view.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, view)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", view.ID, view.Name, view.ModelName)
}

func loadCustomViewList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.CustomViewList, []client.CustomViewSummary, error) {
	views, err := client.ListCustomViews(ctx, runtime.graphqlClient, limit)
	return views, views.CustomViews, err
}

func loadCustomView(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.CustomViewSummary, error) {
	return client.GetCustomViewByID(ctx, runtime.graphqlClient, id)
}

func customViewPageWithItems(
	page client.CustomViewList,
	views []client.CustomViewSummary,
) client.CustomViewList {
	page.CustomViews = views
	return page
}
