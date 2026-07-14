//go:build windows

package auth

import (
	"context"
	"errors"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	lockFileFailImmediately   = 0x00000001
	lockFileExclusiveLock     = 0x00000002
	lockFileLockViolation     = syscall.Errno(33)
	authFileLockRetryInterval = 10 * time.Millisecond
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	lockFileExProc   = kernel32.NewProc("LockFileEx")
	unlockFileExProc = kernel32.NewProc("UnlockFileEx")
)

func lockStateFile(ctx context.Context, file *os.File) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		overlapped := syscall.Overlapped{}
		result, _, callErr := lockFileExProc.Call(
			file.Fd(),
			lockFileExclusiveLock|lockFileFailImmediately,
			0,
			1,
			0,
			uintptr(unsafe.Pointer(&overlapped)),
		)
		if result != 0 {
			return nil
		}
		if !errors.Is(callErr, lockFileLockViolation) {
			return callErr
		}

		timer := time.NewTimer(authFileLockRetryInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}

func unlockStateFile(file *os.File) error {
	overlapped := syscall.Overlapped{}
	result, _, callErr := unlockFileExProc.Call(
		file.Fd(),
		0,
		1,
		0,
		uintptr(unsafe.Pointer(&overlapped)),
	)
	if result == 0 {
		return callErr
	}

	return nil
}
