package upload

import (
	"context"
	"mime/multipart"
)

type UploadResult struct {
	URL      string
	PublicID string
}

// Uploader is a provider-agnostic interface for handling file uploads.
// When adding a new provider (e.g., AWS S3), implement this interface in a new file (e.g., s3.go)
// and wire it up in cmd/main.go. Domain code should only ever depend on this interface.
type Uploader interface {
	Upload(ctx context.Context, file multipart.File, folder string) (UploadResult, error)
	Delete(ctx context.Context, publicID string) error
}
