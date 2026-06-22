package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

// httpDoer is the subset of *http.Client used for file transfers; it is a var so
// tests can substitute the network round-trip.
type httpDoer interface {
	Do(request *http.Request) (*http.Response, error)
}

var fileHTTPClient httpDoer = http.DefaultClient

// closeQuietly closes a response body; the close error on an already-consumed
// body is not actionable on the upload/download paths.
func closeQuietly(closer io.Closer) {
	_ = closer.Close() //nolint:errcheck // consumed-body close error is not actionable.
}

// fileUploadResult is the structured confirmation of a completed upload.
type fileUploadResult struct {
	AssetURL string `json:"asset_url"`
}

// fileDownloadResult is the structured confirmation of a completed download.
type fileDownloadResult struct {
	Path  string `json:"path"`
	Bytes int    `json:"bytes"`
}

func addFilesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	filesCommand := &cobra.Command{
		Use:   "files",
		Short: "Upload and download Linear file assets",
	}
	addFilesUploadCommand(ctx, filesCommand, options)
	addFilesDownloadCommand(ctx, filesCommand, options)
	root.AddCommand(filesCommand)
}

func addFilesUploadCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	contentType := ""
	command := &cobra.Command{
		Use:   "upload PATH",
		Short: "Upload a file and print its Linear asset URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runFileUpload(ctx, command, options, args[0], contentType)
		},
	}
	command.Flags().StringVar(
		&contentType, "content-type", "",
		"MIME type; inferred from the file extension when empty",
	)
	root.AddCommand(command)
}

func runFileUpload(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	path string,
	contentType string,
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	//nolint:gosec // G304: the upload command's purpose is to read the user-named file.
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if contentType == "" {
		contentType = inferContentType(path)
	}
	upload, err := client.PrepareFileUpload(ctx, runtime.graphqlClient, filepath.Base(path), contentType, len(data))
	if err != nil {
		return err
	}
	if err := putFileContents(ctx, upload, data); err != nil {
		return err
	}

	return writeAssetURL(command, options, upload.AssetURL)
}

func inferContentType(path string) string {
	if contentType := mime.TypeByExtension(filepath.Ext(path)); contentType != "" {
		return contentType
	}

	return "application/octet-stream"
}

func putFileContents(ctx context.Context, upload client.FileUpload, data []byte) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, upload.UploadURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", upload.ContentType)
	for _, header := range upload.Headers {
		request.Header.Set(header.Key, header.Value)
	}
	response, err := fileHTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("upload to storage: %w", err)
	}
	defer closeQuietly(response.Body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("upload to storage: unexpected status %d", response.StatusCode)
	}

	return nil
}

func addFilesDownloadCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	output := ""
	command := &cobra.Command{
		Use:   "download URL",
		Short: "Download a file asset to a local path",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runFileDownload(ctx, command, options, args[0], output)
		},
	}
	command.Flags().StringVar(&output, "output", "", "local path to write the downloaded file")
	root.AddCommand(command)
}

func runFileDownload(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	url string,
	output string,
) error {
	if output == "" {
		return errors.New("--output is required")
	}
	//nolint:gosec // G107: the download command's purpose is to fetch the user-provided URL.
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	response, err := fileHTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("download %s: %w", url, err)
	}
	defer closeQuietly(response.Body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("download %s: unexpected status %d", url, response.StatusCode)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read %s: %w", url, err)
	}
	if err := os.WriteFile(output, data, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", output, err)
	}

	return writeDownloadResult(command, options, output, len(data))
}

func writeAssetURL(command *cobra.Command, options *rootOptions, assetURL string) error {
	if wrote, err := writeIDOnly(command, options, assetURL); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, fileUploadResult{AssetURL: assetURL})
	}

	return render.WriteLine(command.OutOrStdout(), "%s", assetURL)
}

func writeDownloadResult(command *cobra.Command, options *rootOptions, path string, size int) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, fileDownloadResult{Path: path, Bytes: size})
	}

	return render.WriteLine(command.OutOrStdout(), "%s %d bytes", path, size)
}
