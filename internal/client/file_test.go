package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func uploadFileResponseJSON(success bool, hasFile bool) string {
	file := `null`
	if hasFile {
		file = `{"filename":"a.png","contentType":"image/png","size":12,` +
			`"uploadUrl":"https://uploads.example/put","assetUrl":"https://assets.example/a.png",` +
			`"headers":[{"key":"x-amz-meta","value":"v"}]}`
	}
	successJSON := "false"
	if success {
		successJSON = "true"
	}

	return `{"fileUpload":{"success":` + successJSON + `,"uploadFile":` + file + `}}`
}

func Test_PrepareFileUpload_returns_signed_target_on_success(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"fileUpload": uploadFileResponseJSON(true, true),
	})

	upload, err := PrepareFileUpload(context.Background(), graphqlClient, "a.png", "image/png", 12)

	require.NoError(t, err)
	require.Equal(t, "https://assets.example/a.png", upload.AssetURL)
	require.Equal(t, "https://uploads.example/put", upload.UploadURL)
	require.Len(t, upload.Headers, 1)
	require.Equal(t, "x-amz-meta", upload.Headers[0].Key)
}

func Test_PrepareFileUpload_requires_filename(t *testing.T) {
	_, err := PrepareFileUpload(context.Background(), fakeGraphQLClient(map[string]string{}), "", "image/png", 12)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_PrepareFileUpload_requires_content_type(t *testing.T) {
	_, err := PrepareFileUpload(context.Background(), fakeGraphQLClient(map[string]string{}), "a.png", "", 12)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_PrepareFileUpload_requires_positive_size(t *testing.T) {
	_, err := PrepareFileUpload(context.Background(), fakeGraphQLClient(map[string]string{}), "a.png", "image/png", 0)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_PrepareFileUpload_wraps_mutation_error(t *testing.T) {
	_, err := PrepareFileUpload(context.Background(), fakeGraphQLClient(map[string]string{}), "a.png", "image/png", 12)

	require.Error(t, err)
	require.Contains(t, err.Error(), "prepare file upload")
}

func Test_PrepareFileUpload_fails_when_no_upload_target(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"fileUpload": uploadFileResponseJSON(false, true),
	})

	_, err := PrepareFileUpload(context.Background(), graphqlClient, "a.png", "image/png", 12)

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_PrepareFileUpload_fails_when_upload_file_missing(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"fileUpload": uploadFileResponseJSON(true, false),
	})

	_, err := PrepareFileUpload(context.Background(), graphqlClient, "a.png", "image/png", 12)

	require.ErrorIs(t, err, ErrMutationFailed)
}
