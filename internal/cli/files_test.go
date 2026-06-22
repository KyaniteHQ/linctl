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

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeHTTPDoer struct {
	status  int
	body    string
	bodyErr bool
	err     error
}

func (doer *fakeHTTPDoer) Do(_ *http.Request) (*http.Response, error) {
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

func swapFileHTTPClient(doer httpDoer) func() {
	original := fileHTTPClient
	fileHTTPClient = doer

	return func() { fileHTTPClient = original }
}

func runFilesCommand(t *testing.T, doer httpDoer, args []string) (string, error) {
	t.Helper()
	restoreDoer := swapFileHTTPClient(doer)
	defer restoreDoer()
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
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
	output, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK}, []string{
		"files", "upload", writeUploadFile(t),
	})

	require.NoError(t, err)
	require.Contains(t, output, "https://assets.example/file.txt")
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

func Test_Files_upload_reports_mutation_error(t *testing.T) {
	restoreDoer := swapFileHTTPClient(&fakeHTTPDoer{status: http.StatusOK})
	defer restoreDoer()
	restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "fileUpload"})
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
	err := putFileContents(context.Background(), client.FileUpload{UploadURL: "://bad", ContentType: "text/plain"}, []byte("x"))

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

func Test_inferContentType_uses_extension_or_default(t *testing.T) {
	require.Contains(t, inferContentType("photo.png"), "image/png")
	require.Equal(t, "application/octet-stream", inferContentType("mystery-file"))
}
