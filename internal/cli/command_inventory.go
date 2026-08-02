package cli

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// Command metadata rides on cobra annotations. The rule is that a command is
// stamped where it is registered: addCommandWithSafety, addWriteCommand, or
// annotateReadCollectionCommand. newGroupCommand is the one exception, because
// a group command is read by construction and has no registration helper of its
// own.
const (
	commandCollectionKeyAnnotation = "linctl.collection_key"
	commandSafetyAnnotation        = "linctl.safety"
	commandIrreversibleAnnotation  = "linctl.irreversible"
	commandWriteEffectAnnotation   = "linctl.write_effect"
)

// CommandWriteEffect names what a write-classified command changes. Safety says
// a command writes; the effect says where the change lands, so generated
// documentation never has to infer it from prose.
type CommandWriteEffect string

// Write effect values used by the command metadata inventory.
const (
	// WriteEffectNone is the effect of a command that does not write.
	WriteEffectNone CommandWriteEffect = ""
	// WriteEffectGuarded changes Linear only after the pinned target matches
	// the target resolved from the live credential.
	WriteEffectGuarded CommandWriteEffect = "guarded"
	// WriteEffectUnguarded changes Linear through an entity that carries no
	// team or project scope to compare against the pinned target.
	WriteEffectUnguarded CommandWriteEffect = "unguarded"
	// WriteEffectLocal writes only on the local machine.
	WriteEffectLocal CommandWriteEffect = "local"
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
	WriteEffect      CommandWriteEffect
	Irreversible     bool
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
		WriteEffect:   commandWriteEffect(command),
		Irreversible:  commandIrreversible(command),
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

// annotateReadCollectionCommand stamps both facts a read list command carries:
// the collection key and the read safety class. List commands that build their
// own RunE use it instead of addCommandWithSafety, because they also need the
// key.
func annotateReadCollectionCommand(command *cobra.Command, collectionKey string) {
	annotateCollectionKey(command, collectionKey)
	annotateCommand(command, commandSafetyAnnotation, string(CommandSafetyRead))
}

// annotateIrreversibleWrite marks a guarded write that linctl cannot undo, so
// generated documentation warns about it instead of a hand-maintained list.
func annotateIrreversibleWrite(command *cobra.Command) {
	annotateCommand(command, commandIrreversibleAnnotation, "true")
}

// addWriteCommand registers a write command with an explicit write effect.
// Commands built through the shared guarded-write pipeline get the effect from
// that pipeline; commands that register their own RunE declare it here.
func addWriteCommand(root *cobra.Command, effect CommandWriteEffect, command *cobra.Command) {
	annotateCommand(command, commandWriteEffectAnnotation, string(effect))
	addCommandWithSafety(root, CommandSafetyWrite, command)
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

// The four readers below share one contract: the command is registered and
// non-nil, and the value is whatever the registrar stamped. Only the annotate
// helpers above write these keys, and they write typed constants, so no reader
// validates.

func commandCollectionKey(command *cobra.Command) string {
	return command.Annotations[commandCollectionKeyAnnotation]
}

func commandWriteEffect(command *cobra.Command) CommandWriteEffect {
	return CommandWriteEffect(command.Annotations[commandWriteEffectAnnotation])
}

func commandIrreversible(command *cobra.Command) bool {
	return command.Annotations[commandIrreversibleAnnotation] == "true"
}

// commandSafety reads the safety class stamped at registration. Every linctl
// command carries the annotation, so the prose of a Short string can never
// reclassify a command. Only the cobra-generated completion subcommands, which
// linctl does not register, reach the path check.
func commandSafety(command *cobra.Command) CommandSafety {
	if safety, stamped := command.Annotations[commandSafetyAnnotation]; stamped {
		return CommandSafety(safety)
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

// mustCollectionFieldForList resolves the single collection field of Page and
// proves it holds []Item. It panics at registration, so no command can reach
// the item pipeline with a page and item type that disagree.
func mustCollectionFieldForList[Page any, Item any]() reflect.StructField {
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

	return field
}

func mustCollectionKeyForList[Page any, Item any]() string {
	return jsonFieldName(mustCollectionFieldForList[Page, Item]())
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
