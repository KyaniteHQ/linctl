package cli

import (
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

// CommandInfo is the normalized command metadata used by generators and
// drift checks that need the public Cobra surface without re-walking it.
type CommandInfo struct {
	Path          string
	UseLine       string
	Short         string
	Aliases       []string
	Entity        string
	TargetArgs    []string
	Safety        CommandSafety
	CollectionKey string
	DocCategory   string
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

func annotateReadCollectionCommand(command *cobra.Command, collectionKey string) {
	annotateCommand(command, commandCollectionKeyAnnotation, collectionKey)
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
	path := " " + CommandPath(command) + " "
	for _, action := range []string{
		" archive ",
		" bulk-export ",
		" create ",
		" delete ",
		" done ",
		" download ",
		" import ",
		" next ",
		" resolve ",
		" restore ",
		" retire ",
		" unarchive ",
		" unresolve ",
		" update ",
		" upload ",
	} {
		if strings.Contains(path, action) {
			return CommandSafetyWrite
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
	for pageType.Kind() == reflect.Pointer {
		pageType = pageType.Elem()
	}
	if pageType.Kind() != reflect.Struct {
		return ""
	}

	var collectionKey string
	for index := range pageType.NumField() {
		field := pageType.Field(index)
		if field.PkgPath != "" || field.Type.Kind() != reflect.Slice {
			continue
		}
		key := jsonFieldName(field)
		if key == "" {
			continue
		}
		if collectionKey != "" {
			return ""
		}
		collectionKey = key
	}

	return collectionKey
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

// CollectionKeys returns the explicit collection field names that list-page
// envelopes may project over with --fields.
func CollectionKeys() []string {
	keys := make([]string, len(collectionKeys))
	copy(keys, collectionKeys)

	return keys
}

// collectionKeys is the explicit allowlist of collection field names that list
// pages emit (each list page carries exactly one such array plus scalar
// pagination/context fields). It is deliberately NOT generic top-level []any
// detection: some detail responses embed an incidental array that is not a
// collection (for example a TimeScheduleSummary's "entries", or a
// ProjectSummary's "teams"), so treating "the single top-level array" as the
// collection would wrongly project per-element instead of over the object.
// Equally, list pages are not all paginated (AuditEntryTypeList,
// SemanticSearchList, SLAConfigurationList, TemplateList carry no
// has_next_page), so a pagination marker cannot stand in for the allowlist
// either. Multi-array responses (IssueDependencyGraph) and detail objects fall
// through to whole-object projection.
var collectionKeys = []string{
	"issues",
	"associations",
	"cycles",
	"projects",
	"members",
	"comments",
	"updates",
	"milestones",
	"documents",
	"labels",
	"teams",
	"users",
	"memberships",
	"drafts",
	"initiatives",
	"notifications",
	"notification_subscriptions",
	"release_pipelines",
	"release_stages",
	"releases",
	"history",
	"links",
	"release_notes",
	"customers",
	"customer_needs",
	"customer_statuses",
	"customer_tiers",
	"relations",
	"roadmaps",
	"time_schedules",
	"triage_responsibilities",
	"sla_configurations",
	"results",
	"templates",
	"workflow_states",
	"agent_activities",
	"agent_skills",
	"external_users",
	"audit_entry_types",
	"favorites",
	"emojis",
	"attachments",
	"custom_views",
	"project_labels",
	"project_statuses",
	"spans",
	"git_automation_states",
}
