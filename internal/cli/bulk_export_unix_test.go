//go:build !windows

package cli

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CommandFlows_issue_bulk_export_reports_a_failed_write_and_leaves_no_file(t *testing.T) {
	dir := t.TempDir()
	var original syscall.Rlimit
	require.NoError(t, syscall.Getrlimit(syscall.RLIMIT_FSIZE, &original))
	limited := original
	limited.Cur = 0
	require.NoError(t, syscall.Setrlimit(syscall.RLIMIT_FSIZE, &limited))
	restore := func() { require.NoError(t, syscall.Setrlimit(syscall.RLIMIT_FSIZE, &original)) }
	t.Cleanup(restore)

	_, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "bulk-export", filepath.Join(dir, "issues.csv")})
	restore()

	require.ErrorIs(t, err, syscall.EFBIG)
	require.ErrorContains(t, err, "encode csv export")
	entries, readErr := os.ReadDir(dir)
	require.NoError(t, readErr)
	require.Empty(t, entries)
}
