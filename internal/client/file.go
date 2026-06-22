package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// FileUploadHeader is one HTTP header required for the storage PUT request.
type FileUploadHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// FileUpload is a prepared upload target: where to PUT the bytes and the
// permanent asset URL the file is reachable at afterwards.
type FileUpload struct {
	Filename    string             `json:"filename"`
	ContentType string             `json:"content_type"`
	Size        int                `json:"size"`
	UploadURL   string             `json:"upload_url"`
	AssetURL    string             `json:"asset_url"`
	Headers     []FileUploadHeader `json:"headers"`
}

// PrepareFileUpload asks Linear for a pre-signed upload target for a file. It is
// a workspace-level asset operation, not a target-pinned write: the returned
// asset URL is attached to an issue or project through the existing guarded
// attachment commands.
func PrepareFileUpload(
	ctx context.Context,
	graphqlClient graphql.Client,
	filename string,
	contentType string,
	size int,
) (FileUpload, error) {
	if filename == "" {
		return FileUpload{}, fmt.Errorf("%w: filename is required", ErrWriteInvalid)
	}
	if contentType == "" {
		return FileUpload{}, fmt.Errorf("%w: content type is required", ErrWriteInvalid)
	}
	if size <= 0 {
		return FileUpload{}, fmt.Errorf("%w: file size must be positive", ErrWriteInvalid)
	}
	result, err := fileUpload(ctx, graphqlClient, contentType, filename, size)
	if err != nil {
		return FileUpload{}, fmt.Errorf("prepare file upload: %w", err)
	}
	if !result.FileUpload.Success || result.FileUpload.UploadFile == nil {
		return FileUpload{}, fmt.Errorf("%w: fileUpload returned no upload target", ErrMutationFailed)
	}

	file := result.FileUpload.UploadFile
	headers := make([]FileUploadHeader, 0, len(file.Headers))
	for _, header := range file.Headers {
		headers = append(headers, FileUploadHeader(header))
	}

	return FileUpload{
		Filename:    file.Filename,
		ContentType: file.ContentType,
		Size:        file.Size,
		UploadURL:   file.UploadUrl,
		AssetURL:    file.AssetUrl,
		Headers:     headers,
	}, nil
}
