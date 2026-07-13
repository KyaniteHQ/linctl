// Package main starts the linctl binary.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/KyaniteHQ/linctl/internal/cli"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := cli.Execute(ctx, cli.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
	if err != nil {
		os.Exit(1)
	}
}
