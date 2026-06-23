package cli

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// CommandInfo is the normalized command metadata used by generators and
// drift checks that need the public Cobra surface without re-walking it.
type CommandInfo struct {
	Path    string
	UseLine string
	Short   string
	Aliases []string
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
		Path:    CommandPath(command),
		UseLine: command.UseLine(),
		Short:   command.Short,
		Aliases: aliases,
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
