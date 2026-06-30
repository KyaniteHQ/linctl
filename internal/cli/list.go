package cli

import (
	"context"
	"errors"

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

type readGetLoader[Item any] func(context.Context, commandRuntime, string) (Item, error)

type readListGetSpec[Page any, Item any] struct {
	Use           string
	Short         string
	ListShort     string
	LimitHelp     string
	GetUse        string
	GetShort      string
	LoadList      readListLoader[Page, Item]
	PageWithItems readListPage[Page, Item]
	LoadGet       readGetLoader[Item]
	WriteItem     readListItemWriter[Item]
}

func addReadListGetCommand[Page any, Item any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	spec readListGetSpec[Page, Item],
) *cobra.Command {
	limit := 50
	parentCommand := &cobra.Command{
		Use:   spec.Use,
		Short: spec.Short,
	}

	listCommand := &cobra.Command{
		Use:   "list",
		Short: spec.ListShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				spec.LoadList,
				spec.PageWithItems,
				spec.WriteItem,
			)
		},
	}
	annotateReadCollectionCommand(listCommand, collectionKeyForPage[Page]())
	listCommand.Flags().IntVar(&limit, "limit", limit, spec.LimitHelp)

	getCommand := &cobra.Command{
		Use:   spec.GetUse,
		Short: spec.GetShort,
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			item, err := spec.LoadGet(ctx, runtime, args[0])
			if err != nil {
				return err
			}

			return spec.WriteItem(command, options, item)
		},
	}
	parentCommand.AddCommand(listCommand, getCommand)
	root.AddCommand(parentCommand)

	return parentCommand
}

func runReadListCommand[Page any, Item any](
	ctx context.Context,
	command *cobra.Command,
	args []string,
	options *rootOptions,
	limit int,
	loader readListLoader[Page, Item],
	pageWithItems readListPage[Page, Item],
	writeOne readListItemWriter[Item],
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
		annotateReadCollectionCommand(command, collectionKeyForPage[Page]())
		return writeJSONValue(command, options, pageWithItems(page, items))
	}
	for _, item := range items {
		if err := writeOne(command, options, item); err != nil {
			return err
		}
	}

	return nil
}

func addChildListCommand[List any, Item any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	use string,
	short string,
	limitHelp string,
	fetch func(commandRuntime, string, int) (List, error),
	count func(List) int,
	sortList func(List) (List, error),
	writeItem any,
	items func(List) []Item,
) {
	limit := 50
	command := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			list, err := fetch(runtime, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, count(list)); err != nil {
				return err
			}
			list, err = sortList(list)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, list)
			}
			for _, item := range items(list) {
				if err := writeChildListItem(command, options, writeItem, item); err != nil {
					return err
				}
			}

			return nil
		},
	}
	annotateReadCollectionCommand(command, collectionKeyForPage[List]())
	command.Flags().IntVar(&limit, "limit", limit, "maximum "+limitHelp+" to return")
	root.AddCommand(command)
}

func writeChildListItem[Item any](
	command *cobra.Command,
	options *rootOptions,
	writeItem any,
	item Item,
) error {
	switch writer := writeItem.(type) {
	case readListItemWriter[Item]:
		return writer(command, options, item)
	case func(*cobra.Command, *rootOptions, Item) error:
		return writer(command, options, item)
	case func(*cobra.Command, Item) error:
		return writer(command, item)
	default:
		return errors.New("child list writer has unsupported signature")
	}
}
