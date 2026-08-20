package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

var openUploadFile = os.Open

func addFilesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	filesCommand := newGroupCommand("files", "Upload and download Linear file assets")
	addFilesUploadCommand(ctx, filesCommand, options)
	addFilesDownloadCommand(ctx, filesCommand, options)
	addCommandWithSafety(root, CommandSafetyWrite, filesCommand)
}

func addFilesUploadCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	contentType := ""
	command := &cobra.Command{
		Use:   "upload PATH",
		Short: "Upload a file and show the Linear asset URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runFileUpload(ctx, command, options, args[0], contentType)
		},
	}
	command.Flags().StringVar(
		&contentType, "content-type", "",
		"MIME type of the file, which linctl infers from the file extension when empty",
	)
	addWriteCommand(root, WriteEffectUnguarded, command)
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
	file, err := openUploadFile(path)
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
		runtime.config.Target,
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
	maxSize := ""
	command := &cobra.Command{
		Use:   "download URL",
		Short: "Download a file asset to a local path",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			var limit *downloadMaxSize
			if command.Flags().Changed("max-size") {
				bytes, err := parseDownloadMaxSize(maxSize)
				if err != nil {
					return err
				}
				limit = &downloadMaxSize{bytes: bytes, flag: maxSize}
			}

			return runFileDownload(ctx, command, options, args[0], output, limit)
		},
	}
	command.Flags().StringVar(&output, "output", "", "local path to write the downloaded file")
	command.Flags().StringVar(
		&maxSize, "max-size", "",
		"fail if the download exceeds this size (bytes, or KB, MB, GB); unset means no limit",
	)
	addWriteCommand(root, WriteEffectLocal, command)
}

func runFileDownload(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	url string,
	output string,
	limit *downloadMaxSize,
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
	size, err := writeDownloadedFile(response.Body, output, limit)
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

// downloadMaxSize is the opt-in byte cap for files download. A nil
// *downloadMaxSize means --max-size was absent and the copy stays unbounded,
// so the absent case needs no separate flag alongside the value.
type downloadMaxSize struct {
	bytes int64
	flag  string
}

func parseDownloadMaxSize(value string) (int64, error) {
	numberText, multiplier := splitDownloadSizeSuffix(strings.TrimSpace(value))
	number, err := strconv.ParseInt(numberText, 10, 64)
	if err != nil || number < 0 {
		return 0, invalidDownloadMaxSize(value)
	}
	if multiplier > 1 && number > math.MaxInt64/multiplier {
		return 0, invalidDownloadMaxSize(value)
	}

	return number * multiplier, nil
}

var downloadSizeMultipliers = []struct {
	suffix     string
	multiplier int64
}{
	{"GB", 1024 * 1024 * 1024},
	{"MB", 1024 * 1024},
	{"KB", 1024},
}

func splitDownloadSizeSuffix(value string) (string, int64) {
	upper := strings.ToUpper(value)
	for _, unit := range downloadSizeMultipliers {
		if strings.HasSuffix(upper, unit.suffix) {
			return strings.TrimSuffix(upper, unit.suffix), unit.multiplier
		}
	}

	return upper, 1
}

func invalidDownloadMaxSize(value string) error {
	return fmt.Errorf(
		"%w: --max-size must be a byte count or KB, MB, or GB, got %q",
		client.ErrWriteInvalid,
		value,
	)
}

func writeDownloadedFile(body io.Reader, output string, limit *downloadMaxSize) (int64, error) {
	var size int64
	err := writeFileAtomically(output, func(writer io.Writer) error {
		var copyErr error
		size, copyErr = io.Copy(writer, limitedDownloadBody(body, limit))
		if copyErr != nil {
			return copyErr
		}
		if downloadExceedsMaxSize(size, limit) {
			return fmt.Errorf("%w: download exceeds --max-size %s", client.ErrWriteInvalid, limit.flag)
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return size, nil
}

func limitedDownloadBody(body io.Reader, limit *downloadMaxSize) io.Reader {
	if limit == nil || limit.bytes == math.MaxInt64 {
		return body
	}

	return io.LimitReader(body, limit.bytes+1)
}

func downloadExceedsMaxSize(size int64, limit *downloadMaxSize) bool {
	return limit != nil && size > limit.bytes
}

func writeFileAtomically(path string, write func(io.Writer) error) error {
	directory := filepath.Dir(path)
	pattern := "." + filepath.Base(path) + ".tmp-*"
	file, err := createDownloadTempFile(directory, pattern)
	if err != nil {
		return err
	}
	tempPath := file.Name()
	keepTemp := false
	defer func() {
		if !keepTemp {
			_ = os.Remove(tempPath) //nolint:errcheck // temp cleanup is best effort after a failed write.
		}
	}()

	writeErr := write(file)
	closeErr := file.Close()
	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return closeErr
	}
	if err := os.Rename(tempPath, path); err != nil {
		return err
	}
	keepTemp = true

	return nil
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
