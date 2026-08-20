package cli

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"
)

// readListLoader produces one page of a list command. The items are read back
// out of the page by pageItems, so a loader never restates which field holds
// them.
type readListLoader[Page any] func(
	context.Context,
	commandRuntime,
	[]string,
	int,
) (Page, error)

type readListItemWriter[Item any] func(*cobra.Command, *rootOptions, Item) error

type readGetLoader[Item any] func(context.Context, commandRuntime, string) (Item, error)

type childListFetcher[Page any] func(
	context.Context,
	graphql.Client,
	string,
	int,
) (Page, error)

// clientList lifts a plain client list call into a readListLoader. Use it when
// the command reads nothing from the runtime but the GraphQL client.
func clientList[Page any](
	list func(context.Context, graphql.Client, int) (Page, error),
) readListLoader[Page] {
	return func(ctx context.Context, runtime commandRuntime, _ []string, limit int) (Page, error) {
		return list(ctx, runtime.graphqlClient, limit)
	}
}

// clientGet lifts a plain client get call into a readGetLoader.
func clientGet[Item any](
	get func(context.Context, graphql.Client, string) (Item, error),
) readGetLoader[Item] {
	return func(ctx context.Context, runtime commandRuntime, id string) (Item, error) {
		return get(ctx, runtime.graphqlClient, id)
	}
}

// addChildListCommand registers a one-argument child listing command: fetch a
// page for the parent entity id and write its items through the shared list
// pipeline.
func addChildListCommand[Page any, Item any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	use string,
	short string,
	limitHelp string,
	fetch childListFetcher[Page],
	writeItem readListItemWriter[Item],
) {
	addListCommand(ctx, root, options, listCommandSpec[Page, Item]{
		Use:       use,
		Short:     short,
		LimitHelp: limitHelp,
		Args:      cobra.ExactArgs(1),
		Load: func(ctx context.Context, runtime commandRuntime, args []string, limit int) (Page, error) {
			return fetch(ctx, runtime.graphqlClient, args[0], limit)
		},
		WriteItem: writeItem,
	})
}

type readListGetSpec[Page any, Item any] struct {
	Use       string
	Short     string
	ListShort string
	LimitHelp string
	GetUse    string
	GetShort  string
	LoadList  readListLoader[Page]
	LoadGet   readGetLoader[Item]
	WriteItem readListItemWriter[Item]
}

func addReadListGetCommand[Page any, Item any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	spec readListGetSpec[Page, Item],
) *cobra.Command {
	limit := defaultListLimit
	parentCommand := newGroupCommand(spec.Use, spec.Short)

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
				spec.WriteItem,
			)
		},
	}
	annotateReadCollectionCommand(listCommand, mustCollectionKeyForList[Page, Item]())
	listCommand.Flags().IntVar(&limit, "limit", limit, spec.LimitHelp)

	parentCommand.AddCommand(listCommand)
	addReadGetCommand(ctx, parentCommand, options, readGetSpec[Item]{
		Use:   spec.GetUse,
		Short: spec.GetShort,
		Load:  spec.LoadGet,
		Write: spec.WriteItem,
	})
	root.AddCommand(parentCommand)

	return parentCommand
}

// readGetSpec describes one single-item read command: Load fetches the item
// named by the command's one argument, and Write renders it. Configure
// registers flags and completions when the command has any.
type readGetSpec[T any] struct {
	Use       string
	Short     string
	Configure func(*cobra.Command)
	Load      readGetLoader[T]
	Write     readListItemWriter[T]
}

// addReadGetCommand registers a one-argument read command through the shared
// pipeline: build the runtime, load the item, write it. Registration stamps
// the read safety class explicitly. It returns the registered command so
// callers can attach subcommands.
func addReadGetCommand[T any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	spec readGetSpec[T],
) *cobra.Command {
	command := &cobra.Command{
		Use:   spec.Use,
		Short: spec.Short,
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			item, err := spec.Load(ctx, runtime, args[0])
			if err != nil {
				return err
			}

			return spec.Write(command, options, item)
		},
	}
	if spec.Configure != nil {
		spec.Configure(command)
	}
	addCommandWithSafety(root, CommandSafetyRead, command)

	return command
}

func runReadListCommand[Page any, Item any](
	ctx context.Context,
	command *cobra.Command,
	args []string,
	options *rootOptions,
	limit int,
	loader readListLoader[Page],
	writeOne readListItemWriter[Item],
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	page, err := loader(ctx, runtime, args, limit)
	if err != nil {
		return err
	}
	items := pageItems[Page, Item](page)
	if err := ensureNonEmpty(options, len(items)); err != nil {
		return err
	}
	items, err = sortByJSONField(items, options.sortField, options.sortOrder)
	if err != nil {
		return err
	}
	if err := noteListTruncation(command, options, page, limit); err != nil {
		return err
	}
	if options.json {
		return writePageJSON(command, options, page, items)
	}
	for _, item := range items {
		if err := writeOne(command, options, item); err != nil {
			return err
		}
	}

	return nil
}

// noteListTruncation writes a stderr note when more pages exist. JSON already
// carries has_next_page on stdout. --id-only and pipes stay on stdout, so the
// note goes to stderr. --quiet suppresses it.
func noteListTruncation[Page any](
	command *cobra.Command,
	options *rootOptions,
	page Page,
	limit int,
) error {
	if options.quiet || !pageHasNextPage(page) {
		return nil
	}

	return writeNote(command, "listing capped at %d items; more pages exist", limit)
}

func pageHasNextPage[Page any](page Page) bool {
	value := reflect.ValueOf(page)
	if !value.IsValid() {
		return false
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return false
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return false
	}
	field := value.FieldByName("HasNextPage")
	if !field.IsValid() || field.Kind() != reflect.Bool {
		return false
	}

	return field.Bool()
}

func writePageJSON[Page any, Item any](
	command *cobra.Command,
	options *rootOptions,
	page Page,
	items []Item,
) error {
	reflectedPage, err := pageWithItems(page, items)
	if err != nil {
		return err
	}

	return writeJSONValue(command, options, reflectedPage)
}

func pageWithItems[Page any, Item any](page Page, items []Item) (Page, error) {
	var zero Page
	pageValue := reflect.ValueOf(page)
	if !pageValue.IsValid() {
		return zero, errors.New("list page value is invalid")
	}

	pageType := pageValue.Type()
	field, err := pageCollectionField(pageType)
	if err != nil {
		return zero, err
	}
	itemsValue := reflect.ValueOf(items)
	if field.Type != itemsValue.Type() {
		return zero, fmt.Errorf(
			"list page %s collection field %q has type %s, not %s",
			pageType,
			field.Name,
			field.Type,
			itemsValue.Type(),
		)
	}

	var structValue reflect.Value
	if pageType.Kind() == reflect.Struct {
		structValue = reflect.New(pageType).Elem()
		structValue.Set(pageValue)
	} else {
		if pageValue.IsNil() {
			return zero, fmt.Errorf("list page %s must not be nil", pageType)
		}
		structValue = reflect.New(pageType.Elem()).Elem()
		structValue.Set(pageValue.Elem())
	}
	structValue.FieldByIndex(field.Index).Set(itemsValue)
	if pageType.Kind() == reflect.Pointer {
		result := reflect.New(pageType.Elem())
		result.Elem().Set(structValue)
		reflect.ValueOf(&zero).Elem().Set(result)
		return zero, nil
	}

	reflect.ValueOf(&zero).Elem().Set(structValue)
	return zero, nil
}

// pageItems reads the collection field that pageWithItems writes. It is the
// read twin of that function: both resolve the same field, so no loader has to
// name it. mustCollectionFieldForList already proved the field holds []Item, so
// the copy below cannot mismatch.
func pageItems[Page any, Item any](page Page) []Item {
	field := mustCollectionFieldForList[Page, Item]()
	var items []Item
	reflect.ValueOf(&items).Elem().Set(
		reflect.Indirect(reflect.ValueOf(page)).FieldByIndex(field.Index),
	)

	return items
}

// listCommandSpec describes one read-only list command in the single list
// pipeline: a loader produces the page and its items, and WriteItem renders one
// human line. The pipeline puts sorted items back into the page for JSON output.
type listCommandSpec[Page any, Item any] struct {
	Use       string
	Short     string
	Long      string
	LimitHelp string
	// Limit is the default --limit value. Zero means defaultListLimit.
	Limit     int
	Args      cobra.PositionalArgs
	Load      readListLoader[Page]
	WriteItem readListItemWriter[Item]
}

const defaultListLimit = 50

// addListCommand registers a list command from its spec. The collection-key
// and read-safety annotations are applied at registration time so the static
// command inventory sees them without executing the command.
func addListCommand[Page any, Item any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	spec listCommandSpec[Page, Item],
) {
	limit := spec.Limit
	if limit == 0 {
		limit = defaultListLimit
	}
	command := &cobra.Command{
		Use:   spec.Use,
		Short: spec.Short,
		Long:  spec.Long,
		Args:  spec.Args,
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				spec.Load,
				spec.WriteItem,
			)
		},
	}
	annotateReadCollectionCommand(command, mustCollectionKeyForList[Page, Item]())
	command.Flags().IntVar(&limit, "limit", limit, "maximum "+spec.LimitHelp+" to return")
	root.AddCommand(command)
}
