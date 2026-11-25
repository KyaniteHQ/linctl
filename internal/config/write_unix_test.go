//go:build !windows

package config

import (
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WritePin_removes_partial_pin_when_write_fails(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".linctl.toml")
	target := Target{OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id"}
	var original syscall.Rlimit
	require.NoError(t, syscall.Getrlimit(syscall.RLIMIT_FSIZE, &original))
	limited := original
	limited.Cur = 0
	require.NoError(t, syscall.Setrlimit(syscall.RLIMIT_FSIZE, &limited))
	restore := func() { require.NoError(t, syscall.Setrlimit(syscall.RLIMIT_FSIZE, &original)) }
	t.Cleanup(restore)

	err := WritePin(path, target)
	restore()

	require.ErrorIs(t, err, syscall.EFBIG)
	require.NoFileExists(t, path)
	require.NoError(t, WritePin(path, target))
}
