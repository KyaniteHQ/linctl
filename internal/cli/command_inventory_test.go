package cli

import (
	"context"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CommandInventory_exposes_stable_public_command_surface(t *testing.T) {
	root := NewRootCommand(context.Background(), BuildInfo{})

	inventory := CommandInventory(root)
	require.NotEmpty(t, inventory)

	paths := make([]string, 0, len(inventory))
	pathsByName := map[string]bool{}
	aliasesByName := map[string]bool{}
	for _, command := range inventory {
		require.NotEmpty(t, command.Path)
		require.NotContains(t, command.Path, "linctl ")
		require.NotEmpty(t, command.UseLine)
		require.NotEmpty(t, command.Short)
		paths = append(paths, command.Path)
		pathsByName[command.Path] = true
		for _, alias := range command.Aliases {
			aliasesByName[alias] = true
		}
	}

	require.True(t, slices.IsSorted(paths), "command inventory must remain stable")
	require.True(t, pathsByName["issue list"])
	require.True(t, pathsByName["organization teams"])
	require.False(t, pathsByName["completion"])
	require.False(t, pathsByName["help"])
	require.True(t, aliasesByName["issue get ISSUE_ID"])
}
