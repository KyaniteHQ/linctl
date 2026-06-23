package cli

import (
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

// httpDoer is the subset of *http.Client used for file transfers.
type httpDoer interface {
	Do(request *http.Request) (*http.Response, error)
}

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
	Bytes int64  `json:"bytes"`
}

type downloadTempFile interface {
	io.Writer
	Close() error
	Name() string
}

var createDownloadTempFile = func(directory string, pattern string) (downloadTempFile, error) {
	return os.CreateTemp(directory, pattern)
}

var newFileTransferHTTPClient = func(options *rootOptions) httpDoer {
	return &http.Client{Timeout: options.timeout}
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
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("read %s: not a regular file", path)
	}
	//nolint:gosec // G304: the upload command's purpose is to read the user-named file.
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	defer closeQuietly(file)
	if contentType == "" {
		contentType = inferContentType(path)
	}
	upload, err := client.PrepareFileUpload(
		ctx,
		runtime.graphqlClient,
		filepath.Base(path),
		contentType,
		int(info.Size()),
	)
	if err != nil {
		return err
	}
	if err := putFileContents(ctx, runtime.fileHTTPClient(), upload, file, info.Size()); err != nil {
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

func putFileContents(
	ctx context.Context,
	httpClient httpDoer,
	upload client.FileUpload,
	content io.Reader,
	size int64,
) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, upload.UploadURL, content)
	if err != nil {
		return err
	}
	request.ContentLength = size
	request.Header.Set("Content-Type", upload.ContentType)
	for _, header := range upload.Headers {
		request.Header.Set(header.Key, header.Value)
	}
	response, err := httpClient.Do(request)
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
	response, err := newFileTransferHTTPClient(options).Do(request)
	if err != nil {
		return fmt.Errorf("download %s: %w", url, err)
	}
	defer closeQuietly(response.Body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("download %s: unexpected status %d", url, response.StatusCode)
	}
	size, err := writeDownloadedFile(response.Body, output)
	if err != nil {
		return fmt.Errorf("write %s: %w", output, err)
	}

	return writeDownloadResult(command, options, output, size)
}

func (runtime commandRuntime) fileHTTPClient() httpDoer {
	if runtime.fileClient != nil {
		return runtime.fileClient
	}

	return http.DefaultClient
}

func writeDownloadedFile(body io.Reader, output string) (int64, error) {
	directory := filepath.Dir(output)
	pattern := "." + filepath.Base(output) + ".tmp-*"
	file, err := createDownloadTempFile(directory, pattern)
	if err != nil {
		return 0, err
	}
	tempPath := file.Name()
	keepTemp := false
	defer func() {
		if !keepTemp {
			_ = os.Remove(tempPath) //nolint:errcheck // temp cleanup is best effort after a failed write.
		}
	}()

	size, copyErr := io.Copy(file, body)
	closeErr := file.Close()
	if copyErr != nil {
		return 0, copyErr
	}
	if closeErr != nil {
		return 0, closeErr
	}
	if err := os.Rename(tempPath, output); err != nil {
		return 0, err
	}
	keepTemp = true

	return size, nil
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

func writeDownloadResult(command *cobra.Command, options *rootOptions, path string, size int64) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, fileDownloadResult{Path: path, Bytes: size})
	}

	return render.WriteLine(command.OutOrStdout(), "%s %d bytes", path, size)
}
