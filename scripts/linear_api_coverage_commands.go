//go:build ignore

package main

import (
	"context"

	"github.com/KyaniteHQ/linctl/internal/cli"
)

func commandInventoryByName(domainCommands []domainCommand) map[string]cli.CommandInfo {
	root := cli.NewRootCommand(context.Background(), cli.BuildInfo{})
	inventory := cli.EnrichCommandInventory(
		cli.CommandInventory(root),
		domainCommandBacking(domainCommands),
	)
	commands := map[string]cli.CommandInfo{}
	for _, command := range inventory {
		commands[command.Path] = command
		for _, alias := range command.Aliases {
			commands[alias] = command
		}
	}
	if command, ok := commands["next"]; ok {
		command.Safety = cli.CommandSafetyRead
		command.TargetScope = "`--dry-run` read-only"
		commands["next --dry-run"] = command
	}

	return commands
}

func domainCommandBacking(commands []domainCommand) map[string]cli.CommandBacking {
	backingByPath := map[string]cli.CommandBacking{}
	for _, command := range commands {
		references := domainCommandReferences(command)
		roots := make([]cli.CommandGraphQLRoot, 0, len(references))
		for _, reference := range references {
			roots = append(roots, cli.CommandGraphQLRoot{
				Kind:      reference.Kind,
				Field:     reference.Field,
				Operation: reference.OperationName,
			})
		}
		backingByPath[command.Command] = cli.CommandBacking{
			OperationBacking: command.Backing,
			TargetScope:      command.Scope,
			GraphQLRoots:     roots,
		}
	}
	return backingByPath
}
