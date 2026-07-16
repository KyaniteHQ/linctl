package cli

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

const (
	commandCollectionKeyAnnotation = "linctl.collection_key"
	commandSafetyAnnotation        = "linctl.safety"
)

// CommandSafety classifies a command's behavior for command-surface audits.
type CommandSafety string

// Command safety values used by the command metadata inventory.
const (
	CommandSafetyRead    CommandSafety = "read"
	CommandSafetyWrite   CommandSafety = "write"
	CommandSafetyLocal   CommandSafety = "local"
	CommandSafetyUnknown CommandSafety = "unknown"
)

// CommandGraphQLRoot describes a GraphQL root field that backs a public command.
type CommandGraphQLRoot struct {
	Kind      string
	Field     string
	Operation string
}

// CommandBacking is the domain-map side of a command's control-surface metadata.
type CommandBacking struct {
	OperationBacking string
	TargetScope      string
	GraphQLRoots     []CommandGraphQLRoot
}

// CommandInfo is the normalized command metadata used by generators and
// drift checks that need the public Cobra surface without re-walking it.
type CommandInfo struct {
	Path             string
	UseLine          string
	Short            string
	Aliases          []string
	Entity           string
	TargetArgs       []string
	Safety           CommandSafety
	CollectionKey    string
	DocCategory      string
	OperationBacking string
	TargetScope      string
	GraphQLRoots     []CommandGraphQLRoot
}

// CommandInventory returns available non-help commands in stable path order.
func CommandInventory(root *cobra.Command) []CommandInfo {
	commands := make([]CommandInfo, 0, len(root.Commands()))
	for _, command := range SortedAvailableCommands(root) {
		commands = append(commands, commandInfo(command))
		commands = append(commands, CommandInventory(command)...)
	}

	return commands
}

// EnrichCommandInventory merges domain-map backing into the normalized command
// inventory without making Cobra registration depend on repo documentation.
func EnrichCommandInventory(commands []CommandInfo, backingByPath map[string]CommandBacking) []CommandInfo {
	enriched := make([]CommandInfo, len(commands))
	for index, command := range commands {
		enriched[index] = command
		backing, ok := commandBacking(command, backingByPath)
		if !ok {
			continue
		}
		enriched[index].OperationBacking = backing.OperationBacking
		enriched[index].TargetScope = backing.TargetScope
		enriched[index].GraphQLRoots = append([]CommandGraphQLRoot(nil), backing.GraphQLRoots...)
	}

	return enriched
}

func commandBacking(command CommandInfo, backingByPath map[string]CommandBacking) (CommandBacking, bool) {
	if backing, ok := backingByPath[command.Path]; ok {
		return backing, true
	}
	for _, alias := range command.Aliases {
		if backing, ok := backingByPath[alias]; ok {
			return backing, true
		}
	}
	return CommandBacking{}, false
}

// SortedAvailableCommands returns the available child commands in stable path order.
func SortedAvailableCommands(parent *cobra.Command) []*cobra.Command {
	commands := make([]*cobra.Command, 0, len(parent.Commands()))
	for _, command := range parent.Commands() {
		if !isInventoryCommand(command) {
			continue
		}
		commands = append(commands, command)
	}
	sort.Slice(commands, func(left int, right int) bool {
		return commands[left].CommandPath() < commands[right].CommandPath()
	})

	return commands
}

func commandInfo(command *cobra.Command) CommandInfo {
	aliases := make([]string, 0, 1)
	if alias := commandUseAlias(command); alias != "" {
		aliases = append(aliases, alias)
	}

	return CommandInfo{
		Path:          CommandPath(command),
		UseLine:       command.UseLine(),
		Short:         command.Short,
		Aliases:       aliases,
		Entity:        commandEntity(command),
		TargetArgs:    commandTargetArgs(command),
		Safety:        commandSafety(command),
		CollectionKey: commandCollectionKey(command),
		DocCategory:   commandDocCategory(command),
	}
}

func isInventoryCommand(command *cobra.Command) bool {
	return command.IsAvailableCommand() && command.Name() != "help" && command.Name() != "completion"
}

// CommandPath returns the command path without the binary name prefix.
func CommandPath(command *cobra.Command) string {
	return strings.TrimPrefix(command.CommandPath(), "linctl ")
}

func commandUseAlias(command *cobra.Command) string {
	use := strings.TrimPrefix(command.UseLine(), "linctl ")
	use = strings.TrimSuffix(use, " [flags]")

	return strings.TrimSpace(use)
}

// addCommandWithSafety registers a command with an explicit safety class so
// the static inventory never falls back to the prose heuristic for it.
func addCommandWithSafety(root *cobra.Command, safety CommandSafety, command *cobra.Command) {
	annotateCommand(command, commandSafetyAnnotation, string(safety))
	root.AddCommand(command)
}

// annotateCollectionKey records the JSON collection key that --fields projects
// over, without touching the command's safety class. Write commands that emit
// a collection envelope (issue import) use it directly.
func annotateCollectionKey(command *cobra.Command, collectionKey string) {
	annotateCommand(command, commandCollectionKeyAnnotation, collectionKey)
}

func annotateReadCollectionCommand(command *cobra.Command, collectionKey string) {
	annotateCollectionKey(command, collectionKey)
	annotateCommand(command, commandSafetyAnnotation, string(CommandSafetyRead))
}

func annotateCommand(command *cobra.Command, key string, value string) {
	if value == "" {
		return
	}
	if command.Annotations == nil {
		command.Annotations = map[string]string{}
	}
	command.Annotations[key] = value
}

func commandCollectionKey(command *cobra.Command) string {
	if command == nil || command.Annotations == nil {
		return ""
	}

	return command.Annotations[commandCollectionKeyAnnotation]
}

func commandSafety(command *cobra.Command) CommandSafety {
	if command != nil && command.Annotations != nil {
		switch CommandSafety(command.Annotations[commandSafetyAnnotation]) {
		case CommandSafetyRead:
			return CommandSafetyRead
		case CommandSafetyWrite:
			return CommandSafetyWrite
		case CommandSafetyLocal:
			return CommandSafetyLocal
		case CommandSafetyUnknown:
			return CommandSafetyUnknown
		}
	}
	for _, prefix := range []string{"Get ", "List ", "Read ", "Search ", "Show ", "Check ", "Suggest "} {
		if strings.HasPrefix(command.Short, prefix) {
			return CommandSafetyRead
		}
	}
	if strings.HasPrefix(CommandPath(command), "completion ") {
		return CommandSafetyLocal
	}

	return CommandSafetyUnknown
}

func commandEntity(command *cobra.Command) string {
	fields := strings.Fields(CommandPath(command))
	if len(fields) == 0 {
		return ""
	}

	return fields[0]
}

func commandDocCategory(command *cobra.Command) string {
	fields := strings.Fields(CommandPath(command))
	if len(fields) == 0 {
		return ""
	}

	return fields[0]
}

func commandTargetArgs(command *cobra.Command) []string {
	parts := strings.Fields(command.UseLine())
	targets := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.Trim(part, "[]")
		if strings.HasSuffix(part, "_ID") || strings.HasSuffix(part, "_KEY") || strings.HasSuffix(part, "_URL") {
			targets = append(targets, part)
		}
	}

	return targets
}

func collectionKeyForPage[Page any]() string {
	pageType := reflect.TypeOf((*Page)(nil)).Elem()
	field, err := pageCollectionField(pageType)
	if err != nil {
		return ""
	}

	return jsonFieldName(field)
}

func mustCollectionKeyForList[Page any, Item any]() string {
	pageType := reflect.TypeOf((*Page)(nil)).Elem()
	field, err := pageCollectionField(pageType)
	if err != nil {
		panic(err)
	}
	itemsType := reflect.TypeOf((*[]Item)(nil)).Elem()
	if field.Type != itemsType {
		panic(fmt.Errorf(
			"list page %s collection field %q has type %s, not %s",
			pageType,
			field.Name,
			field.Type,
			itemsType,
		))
	}

	return jsonFieldName(field)
}

func pageCollectionField(pageType reflect.Type) (reflect.StructField, error) {
	originalType := pageType
	if pageType.Kind() == reflect.Pointer {
		pageType = pageType.Elem()
		if pageType.Kind() != reflect.Struct {
			return reflect.StructField{}, fmt.Errorf("list page %s must point directly to a struct", originalType)
		}
	}
	if pageType.Kind() != reflect.Struct {
		return reflect.StructField{}, fmt.Errorf("list page %s must be a struct or pointer to a struct", originalType)
	}

	var collectionField reflect.StructField
	found := false
	for index := range pageType.NumField() {
		field := pageType.Field(index)
		if field.PkgPath != "" || field.Type.Kind() != reflect.Slice {
			continue
		}
		key := jsonFieldName(field)
		if key == "" {
			continue
		}
		if found {
			return reflect.StructField{}, fmt.Errorf(
				"list page %s has multiple exported JSON slice fields %q and %q",
				originalType,
				collectionField.Name,
				field.Name,
			)
		}
		collectionField = field
		found = true
	}
	if !found {
		return reflect.StructField{}, fmt.Errorf("list page %s has no exported JSON slice field", originalType)
	}

	return collectionField, nil
}

func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return ""
	}
	if name, _, ok := strings.Cut(tag, ","); ok {
		return name
	}
	if tag != "" {
		return tag
	}

	return ""
}
