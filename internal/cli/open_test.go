package cli

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func swapGOOS(goos string) func() {
	original := currentGOOS
	currentGOOS = goos

	return func() { currentGOOS = original }
}

func swapOpenExecutor(fn func(context.Context, string, []string) error) func() {
	original := openExecutor
	openExecutor = fn

	return func() { openExecutor = original }
}

func Test_openCommand_maps_platforms(t *testing.T) {
	name, args := openCommand("linux", "https://x")
	require.Equal(t, "xdg-open", name)
	require.Equal(t, []string{"https://x"}, args)

	name, _ = openCommand("darwin", "https://x")
	require.Equal(t, "open", name)

	name, args = openCommand("windows", "https://x")
	require.Equal(t, "rundll32", name)
	require.Contains(t, args, "https://x")

	name, args = openCommand("plan9", "https://x")
	require.Empty(t, name)
	require.Nil(t, args)
}

func Test_openURL_runs_resolved_opener(t *testing.T) {
	var opened string
	defer swapOpenExecutor(func(_ context.Context, _ string, args []string) error {
		opened = args[len(args)-1]

		return nil
	})()
	defer swapGOOS("linux")()

	require.NoError(t, openURL(context.Background(), "https://linear.app/x"))
	require.Equal(t, "https://linear.app/x", opened)
}

func Test_openURL_rejects_unsupported_platform(t *testing.T) {
	defer swapGOOS("plan9")()

	err := openURL(context.Background(), "https://linear.app/x")

	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported platform")
}

func Test_runOpenExecutable_runs_and_reports(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name, args := "true", []string(nil)
		if runtime.GOOS == "windows" {
			name, args = "cmd", []string{"/c", "exit", "0"}
		}

		require.NoError(t, runOpenExecutable(context.Background(), name, args))
	})

	t.Run("error", func(t *testing.T) {
		err := runOpenExecutable(context.Background(), "linctl-nonexistent-binary-xyz", nil)

		require.Error(t, err)
	})
}

func runOpenFlow(t *testing.T, opener func(context.Context, string, []string) error, args []string) (string, error) {
	t.Helper()
	defer swapOpenExecutor(opener)()
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	output := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs(args)

	err := command.ExecuteContext(context.Background())

	return output.String(), err
}

func noopOpener(_ context.Context, _ string, _ []string) error { return nil }

func Test_CommandFlows_issue_open_prints_url(t *testing.T) {
	output, err := runOpenFlow(t, noopOpener, []string{"issue", "open", "LIT-1"})

	require.NoError(t, err)
	require.Contains(t, output, "https://linear.app/kyanite/issue/LIT-1")
}

func Test_CommandFlows_project_open_prints_url(t *testing.T) {
	output, err := runOpenFlow(t, noopOpener, []string{"project", "open", "project-id"})

	require.NoError(t, err)
	require.Contains(t, output, "https://linear.app/kyanite/project/project-id")
}

func Test_CommandFlows_issue_open_honors_output_flags(t *testing.T) {
	idOnly, err := runOpenFlow(t, noopOpener, []string{"--id-only", "issue", "open", "LIT-1"})
	require.NoError(t, err)
	require.Equal(t, "https://linear.app/kyanite/issue/LIT-1\n", idOnly)

	jsonOut, err := runOpenFlow(t, noopOpener, []string{"--json", "issue", "open", "LIT-1"})
	require.NoError(t, err)
	require.Contains(t, jsonOut, `"url"`)

	quiet, err := runOpenFlow(t, noopOpener, []string{"--quiet", "issue", "open", "LIT-1"})
	require.NoError(t, err)
	require.Empty(t, quiet)
}

func Test_CommandFlows_issue_open_surfaces_opener_error(t *testing.T) {
	_, err := runOpenFlow(t, func(context.Context, string, []string) error {
		return errors.New("opener boom")
	}, []string{"issue", "open", "LIT-1"})

	require.Error(t, err)
	require.Contains(t, err.Error(), "opener boom")
}

func Test_CommandFlows_issue_open_surfaces_resolve_error(t *testing.T) {
	defer swapOpenExecutor(noopOpener)()
	restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue"})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"issue", "open", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
}

func Test_CommandFlows_project_open_surfaces_resolve_error(t *testing.T) {
	defer swapOpenExecutor(noopOpener)()
	restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "project"})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"project", "open", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
}
