package cli

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

// currentGOOS is the platform used to resolve the URL opener; a var for tests.
var currentGOOS = runtime.GOOS

// openExecutor runs the resolved open command; injectable for tests.
var openExecutor = runOpenExecutable

// openedURL is the structured confirmation of an opened entity.
type openedURL struct {
	URL string `json:"url"`
}

func runOpenExecutable(ctx context.Context, name string, args []string) error {
	// name is a fixed platform opener and args carry the Linear entity URL as a
	// discrete argv argument (no shell), so there is no command-injection surface.
	//nolint:gosec // G204: fixed opener launched with the entity URL as an explicit argv arg.
	command := exec.CommandContext(ctx, name, args...)
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("open: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

// openCommand maps a platform to its URL opener and arguments.
func openCommand(goos string, url string) (string, []string) {
	switch goos {
	case "linux":
		return "xdg-open", []string{url}
	case "darwin":
		return "open", []string{url}
	case "windows":
		return "rundll32", []string{"url.dll,FileProtocolHandler", url}
	default:
		return "", nil
	}
}

func openURL(ctx context.Context, url string) error {
	name, args := openCommand(currentGOOS, url)
	if name == "" {
		return fmt.Errorf("unsupported platform %q for open", currentGOOS)
	}

	return openExecutor(ctx, name, args)
}

func runOpenEntity(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	id string,
	resolveURL func(commandRuntime, string) (string, error),
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	url, err := resolveURL(runtime, id)
	if err != nil {
		return err
	}
	if err := openURL(ctx, url); err != nil {
		return err
	}

	return writeOpenedURL(command, options, url)
}

func writeOpenedURL(command *cobra.Command, options *rootOptions, url string) error {
	if wrote, err := writeIDOnly(command, options, url); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, openedURL{URL: url})
	}

	return render.WriteLine(command.OutOrStdout(), "%s", url)
}

func resolveIssueURL(ctx context.Context) func(commandRuntime, string) (string, error) {
	return func(runtime commandRuntime, id string) (string, error) {
		issue, err := client.GetIssueByID(ctx, runtime.graphqlClient, id)

		return issue.URL, err
	}
}

func resolveProjectURL(ctx context.Context) func(commandRuntime, string) (string, error) {
	return func(runtime commandRuntime, id string) (string, error) {
		project, err := client.GetProjectByID(ctx, runtime.graphqlClient, id)

		return project.URL, err
	}
}

func addIssueOpenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "open ISSUE_ID",
		Short: "Open an issue in the default browser",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runOpenEntity(ctx, command, options, args[0], resolveIssueURL(ctx))
		},
	})
}

func addProjectOpenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "open PROJECT_ID",
		Short: "Open a project in the default browser",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runOpenEntity(ctx, command, options, args[0], resolveProjectURL(ctx))
		},
	})
}
