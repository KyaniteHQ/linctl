package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeHTTPDoer struct {
	status         int
	body           string
	bodyErr        bool
	err            error
	requestContext context.Context
	requestSize    int64
}

func (doer *fakeHTTPDoer) Do(request *http.Request) (*http.Response, error) {
	doer.requestContext = request.Context()
	doer.requestSize = request.ContentLength
	if doer.err != nil {
		return nil, doer.err
	}
	body := io.NopCloser(strings.NewReader(doer.body))
	if doer.bodyErr {
		body = io.NopCloser(errorReader{})
	}

	return &http.Response{StatusCode: doer.status, Body: body}, nil
}

type errorReader struct{}

func (errorReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read boom")
}

type closeErrorTempFile struct {
	name string
}

func (file closeErrorTempFile) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}

func (file closeErrorTempFile) Close() error {
	return errors.New("close boom")
}

func (file closeErrorTempFile) Name() string {
	return file.name
}

func useFileTransferHTTPClient(t *testing.T, doer httpDoer) func() {
	t.Helper()
	original := newFileTransferHTTPClient
	newFileTransferHTTPClient = func(_ *rootOptions) httpDoer {
		return doer
	}

	return func() {
		newFileTransferHTTPClient = original
	}
}

func runFilesCommand(t *testing.T, doer httpDoer, args []string) (string, error) {
	t.Helper()
	restore := useCommandRuntimeWithFiles(t, commandFlowFakeClient{}, doer)
	defer restore()
	restoreFileClient := useFileTransferHTTPClient(t, doer)
	defer restoreFileClient()
	output := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs(args)

	err := command.ExecuteContext(context.Background())

	return output.String(), err
}

func writeUploadFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "note.txt")
	require.NoError(t, os.WriteFile(path, []byte("hello world"), 0o600))

	return path
}

func Test_Files_upload_prints_asset_url(t *testing.T) {
	doer := &fakeHTTPDoer{status: http.StatusOK}
	output, err := runFilesCommand(t, doer, []string{
		"files", "upload", writeUploadFile(t),
	})

	require.NoError(t, err)
	require.Contains(t, output, "https://assets.example/file.txt")
	require.Equal(t, int64(len("hello world")), doer.requestSize)
	require.NotNil(t, doer.requestContext)
}

func Test_Files_upload_honors_output_flags(t *testing.T) {
	path := writeUploadFile(t)

	idOnly, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK}, []string{
		"--id-only", "files", "upload", path,
	})
	require.NoError(t, err)
	require.Equal(t, "https://assets.example/file.txt\n", idOnly)

	jsonOut, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK}, []string{
		"--json", "files", "upload", path, "--content-type", "text/plain",
	})
	require.NoError(t, err)
	require.Contains(t, jsonOut, `"asset_url"`)

	quiet, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK}, []string{
		"--quiet", "files", "upload", path,
	})
	require.NoError(t, err)
	require.Empty(t, quiet)
}

func Test_Files_upload_reports_read_error(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.txt")

	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK}, []string{
		"files", "upload", missing,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "read")
}

func Test_Files_upload_reports_non_regular_file(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK}, []string{
		"files", "upload", t.TempDir(),
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "not a regular file")
}

func Test_Files_upload_reports_open_error(t *testing.T) {
	path := writeUploadFile(t)
	require.NoError(t, os.Chmod(path, 0o000))
	t.Cleanup(func() {
		require.NoError(t, os.Chmod(path, 0o600))
	})

	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK}, []string{
		"files", "upload", path,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "read")
}

func Test_Files_upload_reports_mutation_error(t *testing.T) {
	restore := useCommandRuntimeWithFiles(
		t,
		commandFlowFakeClient{failOperation: "fileUpload"},
		&fakeHTTPDoer{status: http.StatusOK},
	)
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"files", "upload", writeUploadFile(t)})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "prepare file upload")
}

func Test_Files_upload_reports_put_failure(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusInternalServerError}, []string{
		"files", "upload", writeUploadFile(t),
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected status")
}

func Test_Files_upload_reports_transport_error(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{err: errors.New("dial boom")}, []string{
		"files", "upload", writeUploadFile(t),
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "upload to storage")
}

func Test_putFileContents_reports_request_build_error(t *testing.T) {
	err := putFileContents(
		context.Background(),
		&fakeHTTPDoer{status: http.StatusOK},
		client.FileUpload{UploadURL: "://bad", ContentType: "text/plain"},
		strings.NewReader("x"),
		1,
	)

	require.Error(t, err)
}

func Test_Files_download_writes_file(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.txt")

	output, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK, body: "payload bytes"}, []string{
		"files", "download", "https://assets.example/file.txt", "--output", out,
	})

	require.NoError(t, err)
	require.Contains(t, output, out)
	data, readErr := os.ReadFile(out)
	require.NoError(t, readErr)
	require.Equal(t, "payload bytes", string(data))
}

func Test_Files_download_honors_output_flags(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.txt")

	jsonOut, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK, body: "data"}, []string{
		"--json", "files", "download", "https://assets.example/file.txt", "--output", out,
	})
	require.NoError(t, err)
	require.Contains(t, jsonOut, `"path"`)

	quiet, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK, body: "data"}, []string{
		"--quiet", "files", "download", "https://assets.example/file.txt", "--output", out,
	})
	require.NoError(t, err)
	require.Empty(t, quiet)
}

func Test_Files_download_requires_output(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK, body: "data"}, []string{
		"files", "download", "https://assets.example/file.txt",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "--output is required")
}

func Test_Files_download_reports_bad_url(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK, body: "data"}, []string{
		"files", "download", "://bad", "--output", filepath.Join(t.TempDir(), "x"),
	})

	require.Error(t, err)
}

func Test_Files_download_reports_http_status(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusNotFound}, []string{
		"files", "download", "https://assets.example/file.txt", "--output", filepath.Join(t.TempDir(), "x"),
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected status")
}

func Test_Files_download_reports_transport_error(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{err: errors.New("dial boom")}, []string{
		"files", "download", "https://assets.example/file.txt", "--output", filepath.Join(t.TempDir(), "x"),
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "download")
}

func Test_Files_download_does_not_build_command_runtime(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.txt")
	restoreFileClient := useFileTransferHTTPClient(t, &fakeHTTPDoer{status: http.StatusOK, body: "data"})
	defer restoreFileClient()
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	output := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{
		"files", "download", "https://assets.example/file.txt", "--output", out,
	})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), out)
	data, readErr := os.ReadFile(out)
	require.NoError(t, readErr)
	require.Equal(t, "data", string(data))
}

func Test_Files_download_reports_body_read_error(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK, bodyErr: true}, []string{
		"files", "download", "https://assets.example/file.txt", "--output", filepath.Join(t.TempDir(), "x"),
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "read")
}

func Test_Files_download_reports_write_error(t *testing.T) {
	unwritable := filepath.Join(t.TempDir(), "missing-dir", "file.txt")

	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK, body: "data"}, []string{
		"files", "download", "https://assets.example/file.txt", "--output", unwritable,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "write")
}

func Test_writeDownloadedFile_removes_temp_file_on_copy_error(t *testing.T) {
	directory := t.TempDir()
	output := filepath.Join(directory, "got.txt")

	_, err := writeDownloadedFile(errorReader{}, output)

	require.Error(t, err)
	entries, readErr := os.ReadDir(directory)
	require.NoError(t, readErr)
	require.Empty(t, entries)
}

func Test_writeDownloadedFile_reports_close_error(t *testing.T) {
	original := createDownloadTempFile
	createDownloadTempFile = func(directory string, pattern string) (downloadTempFile, error) {
		return closeErrorTempFile{name: filepath.Join(directory, pattern)}, nil
	}
	defer func() {
		createDownloadTempFile = original
	}()

	_, err := writeDownloadedFile(strings.NewReader("data"), filepath.Join(t.TempDir(), "got.txt"))

	require.Error(t, err)
	require.Contains(t, err.Error(), "close boom")
}

func Test_writeDownloadedFile_removes_temp_file_on_rename_error(t *testing.T) {
	directory := t.TempDir()
	output := filepath.Join(directory, "existing-dir")
	require.NoError(t, os.Mkdir(output, 0o700))

	_, err := writeDownloadedFile(strings.NewReader("data"), output)

	require.Error(t, err)
	entries, readErr := os.ReadDir(directory)
	require.NoError(t, readErr)
	require.Len(t, entries, 1)
	require.Equal(t, "existing-dir", entries[0].Name())
}

func Test_commandRuntime_fileHTTPClient_uses_default_client(t *testing.T) {
	require.Equal(t, http.DefaultClient, commandRuntime{}.fileHTTPClient())
}

func Test_newFileTransferHTTPClient_uses_root_timeout(t *testing.T) {
	httpClient, ok := newFileTransferHTTPClient(&rootOptions{timeout: 3 * time.Second}).(*http.Client)

	require.True(t, ok)
	require.Equal(t, 3*time.Second, httpClient.Timeout)
}

func Test_inferContentType_uses_extension_or_default(t *testing.T) {
	require.Contains(t, inferContentType("photo.png"), "image/png")
	require.Equal(t, "application/octet-stream", inferContentType("mystery-file"))
}
