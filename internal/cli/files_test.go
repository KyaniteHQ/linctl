package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeHTTPDoer struct {
	status      int
	err         error
	requested   bool
	requestSize int64
}

func (doer *fakeHTTPDoer) Do(request *http.Request) (*http.Response, error) {
	doer.requested = true
	doer.requestSize = request.ContentLength
	if doer.err != nil {
		return nil, doer.err
	}

	return &http.Response{StatusCode: doer.status, Body: http.NoBody}, nil
}

type recordingGraphQLClient struct {
	client commandFlowFakeClient
	calls  []string
}

func (client *recordingGraphQLClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	client.calls = append(client.calls, request.OpName)

	return client.client.MakeRequest(ctx, request, response)
}

type closeErrorTempFile struct {
	name string
}

type writeErrorTempFile struct {
	downloadTempFile
	err error
}

func (file writeErrorTempFile) Write(_ []byte) (int, error) {
	return 0, file.err
}

func stubAtomicTempWriteError(t *testing.T, writeErr error) {
	t.Helper()
	original := createDownloadTempFile
	createDownloadTempFile = func(directory string, pattern string) (downloadTempFile, error) {
		file, err := os.CreateTemp(directory, pattern)
		if err != nil {
			return nil, err
		}

		return writeErrorTempFile{downloadTempFile: file, err: writeErr}, nil
	}
	t.Cleanup(func() {
		createDownloadTempFile = original
	})
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

func runFilesCommand(t *testing.T, doer httpDoer, args []string) (string, error) {
	t.Helper()
	restore := useCommandRuntimeWithFiles(t, commandFlowFakeClient{}, doer)
	defer restore()
	output := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs(args)

	err := command.ExecuteContext(context.Background())

	return output.String(), err
}

func serveDownload(t *testing.T, handler http.HandlerFunc) string {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return server.URL + "/file.txt"
}

func downloadBody(t *testing.T, body string) http.HandlerFunc {
	t.Helper()

	return func(writer http.ResponseWriter, _ *http.Request) {
		if _, err := io.WriteString(writer, body); err != nil {
			t.Errorf("write download body: %v", err)
		}
	}
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

func Test_Files_upload_sends_the_content_type_to_linear(t *testing.T) {
	dir := t.TempDir()
	image := filepath.Join(dir, "photo.png")
	bare := filepath.Join(dir, "mystery-file")
	require.NoError(t, os.WriteFile(image, []byte("png"), 0o600))
	require.NoError(t, os.WriteFile(bare, []byte("raw"), 0o600))
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "extension", args: []string{"files", "upload", image}, want: "image/png"},
		{name: "no extension", args: []string{"files", "upload", bare}, want: "application/octet-stream"},
		{name: "flag", args: []string{"files", "upload", bare, "--content-type", "text/plain"}, want: "text/plain"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fake := &directWriteCaptureClient{
				operation: "fileUpload",
				variables: []expectedWriteVariable{{path: []string{"contentType"}, value: test.want}},
				delegate:  commandFlowFakeClient{},
			}
			restore := useCommandRuntimeWithFiles(t, fake, &fakeHTTPDoer{status: http.StatusOK})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&bytes.Buffer{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Equal(t, 1, fake.calls)
		})
	}
}

func Test_Files_upload_reports_read_error(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.txt")

	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK}, []string{
		"files", "upload", missing,
	})

	require.ErrorContains(t, err, "read")
}

func Test_Files_upload_reports_non_regular_file(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusOK}, []string{
		"files", "upload", t.TempDir(),
	})

	require.ErrorContains(t, err, "not a regular file")
}

func Test_Files_upload_refuses_an_unreadable_file_before_any_upload(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows file modes do not deny reads")
	}
	if os.Geteuid() == 0 {
		t.Skip("root reads a file whatever its mode")
	}
	path := writeUploadFile(t)
	require.NoError(t, os.Chmod(path, 0o000))
	doer := &fakeHTTPDoer{status: http.StatusOK}

	_, err := runFilesCommand(t, doer, []string{"files", "upload", path})

	require.ErrorIs(t, err, os.ErrPermission)
	require.ErrorContains(t, err, "read "+path)
	require.False(t, doer.requested)
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

	require.ErrorContains(t, err, "prepare file upload")
}

func Test_Files_upload_target_mismatch_stops_before_file_upload_and_put(t *testing.T) {
	graphqlClient := &recordingGraphQLClient{client: commandFlowFakeClient{}}
	doer := &fakeHTTPDoer{status: http.StatusOK}
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		mismatched := testCommandRuntime(graphqlClient)
		mismatched.config.Target.OrgID = "other-org"
		mismatched.fileClient = doer
		return mismatched, nil
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"files", "upload", writeUploadFile(t)})

	err := command.ExecuteContext(context.Background())

	require.ErrorIs(t, err, client.ErrTargetMismatch)
	require.NotContains(t, graphqlClient.calls, "fileUpload")
	require.False(t, doer.requested)
}

func Test_Files_upload_reports_put_failure(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{status: http.StatusInternalServerError}, []string{
		"files", "upload", writeUploadFile(t),
	})

	require.ErrorContains(t, err, "unexpected status")
}

func Test_Files_upload_reports_transport_error(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{err: errors.New("dial boom")}, []string{
		"files", "upload", writeUploadFile(t),
	})

	require.ErrorContains(t, err, "upload to storage")
}

func Test_putFileContents_reports_request_build_error(t *testing.T) {
	err := putFileContents(
		context.Background(),
		&fakeHTTPDoer{status: http.StatusOK},
		client.FileUpload{UploadURL: "://bad", ContentType: "text/plain"},
		http.NoBody,
		1,
	)

	require.Error(t, err)
}

func Test_Files_download_writes_file(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.txt")

	output, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", serveDownload(t, downloadBody(t, "payload bytes")), "--output", out,
	})

	require.NoError(t, err)
	require.Contains(t, output, out)
	data, readErr := os.ReadFile(out)
	require.NoError(t, readErr)
	require.Equal(t, "payload bytes", string(data))
}

func Test_Files_download_honors_output_flags(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.txt")
	url := serveDownload(t, downloadBody(t, "data"))

	jsonOut, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"--json", "files", "download", url, "--output", out,
	})
	require.NoError(t, err)
	require.Contains(t, jsonOut, `"path"`)

	quiet, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"--quiet", "files", "download", url, "--output", out,
	})
	require.NoError(t, err)
	require.Empty(t, quiet)
}

func Test_Files_download_requires_output(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", "https://assets.example/file.txt",
	})

	require.ErrorContains(t, err, "--output is required")
}

func Test_Files_download_reports_bad_url(t *testing.T) {
	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", "://bad", "--output", filepath.Join(t.TempDir(), "x"),
	})

	require.Error(t, err)
}

func Test_Files_download_reports_http_status(t *testing.T) {
	url := serveDownload(t, func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
	})

	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", url, "--output", filepath.Join(t.TempDir(), "x"),
	})

	require.ErrorContains(t, err, "unexpected status 404")
}

func Test_Files_download_reports_transport_error(t *testing.T) {
	server := httptest.NewServer(downloadBody(t, "data"))
	url := server.URL + "/file.txt"
	server.Close()

	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", url, "--output", filepath.Join(t.TempDir(), "x"),
	})

	require.ErrorContains(t, err, "download "+url)
}

func Test_Files_download_stops_at_the_timeout_flag(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.txt")
	url := serveDownload(t, func(writer http.ResponseWriter, request *http.Request) {
		select {
		case <-request.Context().Done():
			return
		case <-time.After(2 * time.Second):
		}
		if _, err := io.WriteString(writer, "late"); err != nil {
			t.Errorf("write download body: %v", err)
		}
	})

	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"--timeout", "50ms", "files", "download", url, "--output", out,
	})

	var timeoutErr interface{ Timeout() bool }
	require.ErrorAs(t, err, &timeoutErr)
	require.True(t, timeoutErr.Timeout())
	require.NoFileExists(t, out)
}

func Test_Files_download_does_not_build_command_runtime(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.txt")
	url := serveDownload(t, downloadBody(t, "data"))
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
	command.SetArgs([]string{"files", "download", url, "--output", out})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), out)
	data, readErr := os.ReadFile(out)
	require.NoError(t, readErr)
	require.Equal(t, "data", string(data))
}

func Test_Files_download_reports_body_read_error(t *testing.T) {
	url := serveDownload(t, func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Length", "100")
		if _, err := io.WriteString(writer, "short"); err != nil {
			t.Errorf("write download body: %v", err)
		}
	})

	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", url, "--output", filepath.Join(t.TempDir(), "x"),
	})

	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func Test_Files_download_reports_write_error(t *testing.T) {
	unwritable := filepath.Join(t.TempDir(), "missing-dir", "file.txt")

	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", serveDownload(t, downloadBody(t, "data")), "--output", unwritable,
	})

	require.ErrorContains(t, err, "write")
}

func Test_writeDownloadedFile_reports_close_error(t *testing.T) {
	original := createDownloadTempFile
	createDownloadTempFile = func(directory string, pattern string) (downloadTempFile, error) {
		return closeErrorTempFile{name: filepath.Join(directory, pattern)}, nil
	}
	defer func() {
		createDownloadTempFile = original
	}()

	_, err := writeDownloadedFile(
		http.NoBody,
		filepath.Join(t.TempDir(), "got.txt"),
		nil,
	)

	require.ErrorContains(t, err, "close boom")
}

func Test_Files_download_to_a_directory_fails_and_leaves_no_temp_file(t *testing.T) {
	directory := t.TempDir()
	output := filepath.Join(directory, "existing-dir")
	require.NoError(t, os.Mkdir(output, 0o700))

	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", serveDownload(t, downloadBody(t, "data")), "--output", output,
	})

	require.ErrorContains(t, err, "write "+output)
	entries, readErr := os.ReadDir(directory)
	require.NoError(t, readErr)
	require.Len(t, entries, 1)
	require.Equal(t, "existing-dir", entries[0].Name())
}

func Test_Files_download_without_max_size_writes_full_body(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.bin")
	body := strings.Repeat("x", 2048)

	output, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", serveDownload(t, downloadBody(t, body)), "--output", out,
	})

	require.NoError(t, err)
	require.Contains(t, output, out)
	data, readErr := os.ReadFile(out)
	require.NoError(t, readErr)
	require.Equal(t, body, string(data))
}

func Test_Files_download_rejects_body_over_max_size_and_leaves_no_file(t *testing.T) {
	directory := t.TempDir()
	out := filepath.Join(directory, "got.bin")
	body := strings.Repeat("x", 2048)

	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", serveDownload(t, downloadBody(t, body)),
		"--output", out, "--max-size", "1KB",
	})

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Equal(t, "INVALID_WRITE", errorCode(err))
	require.ErrorContains(t, err, "1KB")
	_, statErr := os.Stat(out)
	require.ErrorIs(t, statErr, os.ErrNotExist)
	entries, readErr := os.ReadDir(directory)
	require.NoError(t, readErr)
	require.Empty(t, entries)
}

func Test_Files_download_accepts_body_equal_to_max_size(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.bin")
	body := strings.Repeat("x", 1024)

	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", serveDownload(t, downloadBody(t, body)),
		"--output", out, "--max-size", "1KB",
	})

	require.NoError(t, err)
	data, readErr := os.ReadFile(out)
	require.NoError(t, readErr)
	require.Equal(t, body, string(data))
}

func Test_Files_download_treats_the_largest_max_size_as_unbounded(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.txt")

	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", serveDownload(t, downloadBody(t, "data")),
		"--output", out, "--max-size", "9223372036854775807",
	})

	require.NoError(t, err)
	data, readErr := os.ReadFile(out)
	require.NoError(t, readErr)
	require.Equal(t, "data", string(data))
}

func Test_Files_download_rejects_invalid_max_size_before_request(t *testing.T) {
	out := filepath.Join(t.TempDir(), "got.bin")
	url := serveDownload(t, func(http.ResponseWriter, *http.Request) {
		t.Error("download must not be requested")
	})

	_, err := runFilesCommand(t, &fakeHTTPDoer{}, []string{
		"files", "download", url, "--output", out, "--max-size", "1TB",
	})

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Equal(t, "INVALID_WRITE", errorCode(err))
	_, statErr := os.Stat(out)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func Test_parseDownloadMaxSize_accepts_plain_bytes_and_suffixes(t *testing.T) {
	cases := []struct {
		raw  string
		want int64
	}{
		{raw: "0", want: 0},
		{raw: "13", want: 13},
		{raw: "1KB", want: 1024},
		{raw: "1kb", want: 1024},
		{raw: " 2MB ", want: 2 * 1024 * 1024},
		{raw: "1GB", want: 1024 * 1024 * 1024},
		{raw: strconv.FormatInt(math.MaxInt64, 10), want: math.MaxInt64},
	}
	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := parseDownloadMaxSize(tc.raw)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func Test_parseDownloadMaxSize_rejects_invalid_and_overflowing_values(t *testing.T) {
	for _, raw := range []string{"", "  ", "nope", "-1", "1TB", "1.5MB", "8589934592GB"} {
		t.Run(raw, func(t *testing.T) {
			_, err := parseDownloadMaxSize(raw)
			require.ErrorIs(t, err, client.ErrWriteInvalid)
			require.ErrorContains(t, err, raw)
		})
	}
}
