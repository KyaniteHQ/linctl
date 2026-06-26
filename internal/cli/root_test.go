package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
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
	for _, flagName := range []string{
		"json",
		"compact",
		"fields",
		"id-only",
		"quiet",
		"fail-on-empty",
		"sort",
		"order",
		"format",
		"profile",
		"org",
		"team",
		"team-id",
		"project",
		"timeout",
	} {
		require.NotNil(t, flags.Lookup(flagName), "missing persistent flag %s", flagName)
	}
	require.Equal(t, "pinned Linear team key", flags.Lookup("team").Usage)
	require.Equal(t, "pinned Linear team id", flags.Lookup("team-id").Usage)
}

func Test_RootCommand_rejects_non_positive_limit_before_runtime(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, nil
	}
	defer func() {
		buildCommandRuntime = original
	}()

	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"issue", "list", "--limit", "0"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "invalid --limit 0")
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

func Test_Execute_runs_root_command_without_args(t *testing.T) {
	// When
	err := Execute(context.Background(), BuildInfo{})

	// Then
	require.NoError(t, err)
}

func Test_execute_runs_with_injected_streams(t *testing.T) {
	// Given
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	// When
	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr, []string{"usage"})

	// Then
	require.NoError(t, err)
	require.Contains(t, stdout.String(), "linctl is a schema-aligned Linear CLI")
	require.Empty(t, stderr.String())
}

func Test_execute_emits_error_envelope_on_failure(t *testing.T) {
	// Given
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	// When
	err := execute(
		context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr, []string{"not-a-real-command"},
	)

	// Then
	require.Error(t, err)
	var envelope errorEnvelope
	require.NoError(t, json.Unmarshal(stderr.Bytes(), &envelope))
	require.Equal(t, "INTERNAL", envelope.ErrorCode)
	require.NotEmpty(t, envelope.Message)
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
	require.Contains(t, stdout.String(), "initiative-to-project")
	require.Contains(t, stdout.String(), "roadmap and roadmap-to-project are legacy read compatibility")
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

func Test_RoadmapHelp_marks_legacy_and_points_to_initiatives(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name: "roadmap",
			args: []string{"roadmap", "--help"},
			expected: []string{
				"deprecated planning surface",
				"linctl initiative",
				"List visible legacy Linear roadmaps",
			},
		},
		{
			name: "roadmap projects",
			args: []string{"roadmap", "projects", "--help"},
			expected: []string{
				"deprecated planning surface",
				"linctl initiative projects",
				"List projects associated with one legacy roadmap",
			},
		},
		{
			name: "roadmap to project",
			args: []string{"roadmap-to-project", "--help"},
			expected: []string{
				"deprecated planning association surface",
				"linctl initiative-to-project",
				"List visible legacy Roadmap-to-Project associations",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdout := bytes.Buffer{}
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&stdout)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			for _, expected := range test.expected {
				require.Contains(t, stdout.String(), expected)
			}
		})
	}
}
