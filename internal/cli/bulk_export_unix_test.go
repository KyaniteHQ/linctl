//go:build !windows

package cli

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func limitFileSizeToZero(t *testing.T) func() {
	t.Helper()
	var original syscall.Rlimit
	require.NoError(t, syscall.Getrlimit(syscall.RLIMIT_FSIZE, &original))
	limited := original
	limited.Cur = 0
	require.NoError(t, syscall.Setrlimit(syscall.RLIMIT_FSIZE, &limited))
	restore := func() { require.NoError(t, syscall.Setrlimit(syscall.RLIMIT_FSIZE, &original)) }
	t.Cleanup(restore)

	return restore
}

func Test_CommandFlows_issue_bulk_export_reports_a_failed_write_and_leaves_no_file(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{name: "issues.csv", message: "encode csv export"},
		{name: "issues.json", message: "write json"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			restore := limitFileSizeToZero(t)

			_, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "bulk-export", filepath.Join(dir, test.name)})
			restore()

			require.ErrorIs(t, err, syscall.EFBIG)
			require.ErrorContains(t, err, test.message)
			entries, readErr := os.ReadDir(dir)
			require.NoError(t, readErr)
			require.Empty(t, entries)
		})
	}
}

func Test_CommandFlows_issue_export_reports_a_failed_write_and_leaves_no_file(t *testing.T) {
	dir := t.TempDir()
	restore := limitFileSizeToZero(t)

	_, _, err := runExportFlow(t, commandFlowFakeClient{}, []string{"issue", "export", "LIT-1", dir})
	restore()

	require.ErrorIs(t, err, syscall.EFBIG)
	require.ErrorContains(t, err, "write "+filepath.Join(dir, "LIT-1.md"))
	entries, readErr := os.ReadDir(dir)
	require.NoError(t, readErr)
	require.Empty(t, entries)
}
