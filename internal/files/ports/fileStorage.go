package ports

import "context"

type FileStorage interface {
	UploadFile(context context.Context, tenant, fileName string, file []byte) error
	GetExpiringUrl(tenant, fileName string, expiresInNumberOfMinutes int, insecure bool) (string, error)
	DeleteFile(context context.Context, tenant, fileName string) error
}
