package cli

import (
	"context"

	"github.com/spf13/cobra"
)

type readListLoader[Page any, Item any] func(
	context.Context,
	commandRuntime,
	[]string,
	int,
) (Page, []Item, error)

type readListPage[Page any, Item any] func(Page, []Item) Page

type readListItemWriter[Item any] func(*cobra.Command, *rootOptions, Item) error

func runReadListCommand[Page any, Item any](
	ctx context.Context,
	command *cobra.Command,
	args []string,
	options *rootOptions,
	limit int,
	loader readListLoader[Page, Item],
	pageWithItems readListPage[Page, Item],
	writeItem readListItemWriter[Item],
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	page, items, err := loader(ctx, runtime, args, limit)
	if err != nil {
		return err
	}
	if err := ensureNonEmpty(options, len(items)); err != nil {
		return err
	}
	items, err = sortByJSONField(items, options.sortField, options.sortOrder)
	if err != nil {
		return err
	}
	if options.json {
		return writeJSONValue(command, options, pageWithItems(page, items))
	}
	for _, item := range items {
		if err := writeItem(command, options, item); err != nil {
			return err
		}
	}

	return nil
}
