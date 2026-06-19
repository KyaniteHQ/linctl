package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RootCommand_exposes_global_flags_when_created(t *testing.T) {
	// Given
	command := NewRootCommand(context.Background(), BuildInfo{
		Version: "test-version",
		Commit:  "test-commit",
		Date:    "test-date",
	})

	// When
	flags := command.PersistentFlags()

	// Then
	for _, flagName := range []string{"json", "profile", "org", "team", "project", "timeout"} {
		require.NotNil(t, flags.Lookup(flagName), "missing persistent flag %s", flagName)
	}
}

func Test_Execute_prints_version_when_version_flag_is_set(t *testing.T) {
	// Given
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{
		Version: "test-version",
		Commit:  "test-commit",
		Date:    "test-date",
	})
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs([]string{"--version"})

	// When
	err := command.ExecuteContext(context.Background())

	// Then
	require.NoError(t, err)
	require.Equal(t, "linctl test-version (commit test-commit, built test-date)\n", stdout.String())
	require.Empty(t, stderr.String())
}

func Test_Usage_prints_overview_when_called(t *testing.T) {
	// Given
	stdout := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&stdout)
	command.SetArgs([]string{"usage"})

	// When
	err := command.ExecuteContext(context.Background())

	// Then
	require.NoError(t, err)
	require.Contains(t, stdout.String(), "linctl is a schema-aligned Linear CLI")
	require.Contains(t, stdout.String(), "linctl issue usage")
}

func Test_Usage_prints_domain_usage_when_called_from_domain(t *testing.T) {
	// Given
	stdout := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&stdout)
	command.SetArgs([]string{"project", "usage"})

	// When
	err := command.ExecuteContext(context.Background())

	// Then
	require.NoError(t, err)
	require.Contains(t, stdout.String(), "project commands cover the safe Linear project loop")
	require.Contains(t, stdout.String(), "project archive PROJECT_ID")
}

func Test_Usage_prints_json_when_json_flag_is_set(t *testing.T) {
	// Given
	stdout := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&stdout)
	command.SetArgs([]string{"--json", "usage", "issue"})

	// When
	err := command.ExecuteContext(context.Background())

	// Then
	require.NoError(t, err)
	require.Contains(t, stdout.String(), `"topic": "issue"`)
	require.Contains(t, stdout.String(), "safe Linear issue loop")
}
