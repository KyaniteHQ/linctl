package gitctx

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GitContextScenarios_resolve_or_report_issue_references(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("fakes git/jj via POSIX shell scripts on PATH; not portable to Windows")
	}

	t.Run("non git directory returns current branch error", func(t *testing.T) {
		_, err := CurrentIssueIdentifier(context.Background(), t.TempDir())

		require.Error(t, err)
		require.Contains(t, err.Error(), "git branch --show-current")
	})

	t.Run("jj description parse is used when git branch lacks issue", func(t *testing.T) {
		dir := t.TempDir()
		runGit(t, dir, "init")
		runGit(t, dir, "checkout", "-b", "main")

		jjPath := filepath.Join(dir, "jj")
		require.NoError(t, os.WriteFile(jjPath, []byte("#!/usr/bin/env bash\nprintf 'Work item\\n\\nLinear-Issue: LIT-77\\n'\n"), 0o700))
		t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

		identifier, err := CurrentIssueIdentifier(context.Background(), dir)

		require.NoError(t, err)
		require.Equal(t, "LIT-77", identifier)
	})

	t.Run("empty git branch and empty jj description are reported", func(t *testing.T) {
		dir := t.TempDir()
		runGit(t, dir, "init")
		jjPath := filepath.Join(dir, "jj")
		require.NoError(t, os.WriteFile(jjPath, []byte("#!/usr/bin/env bash\nexit 0\n"), 0o700))
		t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

		_, err := CurrentIssueIdentifier(context.Background(), dir)

		require.Error(t, err)
		require.Contains(t, err.Error(), `git branch "master"`)
		require.Contains(t, err.Error(), "jj description empty")
	})

	t.Run("jj description without identifier is reported", func(t *testing.T) {
		dir := t.TempDir()
		runGit(t, dir, "init")
		runGit(t, dir, "checkout", "-b", "main")
		jjPath := filepath.Join(dir, "jj")
		require.NoError(t, os.WriteFile(jjPath, []byte("#!/usr/bin/env bash\nprintf 'no issue here\\n'\n"), 0o700))
		t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

		_, err := CurrentIssueIdentifier(context.Background(), dir)

		require.Error(t, err)
		require.Contains(t, err.Error(), "jj description has no identifier")
	})

	t.Run("empty git branch is reported", func(t *testing.T) {
		dir := t.TempDir()
		gitPath := filepath.Join(dir, "git")
		require.NoError(t, os.WriteFile(gitPath, []byte("#!/usr/bin/env bash\nexit 0\n"), 0o700))
		jjPath := filepath.Join(dir, "jj")
		require.NoError(t, os.WriteFile(jjPath, []byte("#!/usr/bin/env bash\nprintf 'no issue here\\n'\n"), 0o700))
		t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

		_, err := CurrentIssueIdentifier(context.Background(), dir)

		require.Error(t, err)
		require.Contains(t, err.Error(), "git branch empty")
	})
}
