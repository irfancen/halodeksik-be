package appcloud

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"halodeksik-be/app/applogger"
	"halodeksik-be/app/env"
	"io"
	"mime/multipart"
)

var AppFileUploader FileUploader

type FileUploader interface {
	UploadFile(ctx context.Context, file multipart.File, object, path string) error
}

func SetAppFileUploader(uploader FileUploader) {
	AppFileUploader = uploader
}

type FileUploaderImpl struct {
	client     *storage.Client
	projectId  string
	bucketName string
}

func NewFileUploaderImpl() *FileUploaderImpl {
	credentialFile := env.Get("GCLOUD_CREDENTIAL_FILE")
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(credentialFile))
	if err != nil {
		applogger.Log.Errorf("failed to create file uploader client: %v", err)
	}

	projectId := env.Get("GCLOUD_STORAGE_PROJECT_ID")
	bucketName := env.Get("GCLOUD_STORAGE_BUCKET_NAME")

	return &FileUploaderImpl{
		client:     client,
		projectId:  projectId,
		bucketName: bucketName,
	}
}

func (f *FileUploaderImpl) UploadFile(ctx context.Context, file multipart.File, path, name string) error {
	bucketObject := f.client.Bucket(f.bucketName).Object(path + name)
	wc := bucketObject.NewWriter(ctx)
	wc.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}
