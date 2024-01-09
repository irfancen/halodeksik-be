package appcloud

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"halodeksik-be/app/applogger"
	"halodeksik-be/app/env"
	"io"
	"mime/multipart"
	"path/filepath"
)

var AppFileUploader FileUploader

type FileUploader interface {
	SendToBucket(ctx context.Context, file multipart.File, object, path string) error
	Upload(ctx context.Context, fileHeader any, folderName string) (string, error)
}

func SetAppFileUploader(uploader FileUploader) {
	AppFileUploader = uploader
}

type FileUploaderImpl struct {
	client     *storage.Client
	projectId  string
	bucketName string
	cloudUrl   string
}

func NewFileUploaderImpl() *FileUploaderImpl {
	credentialFile := env.Get("GCLOUD_CREDENTIAL_FILE")
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(credentialFile))
	if err != nil {
		applogger.Log.Errorf("failed to create file uploader client: %v", err)
	}

	projectId := env.Get("GCLOUD_STORAGE_PROJECT_ID")
	bucketName := env.Get("GCLOUD_STORAGE_BUCKET_NAME")
	cloudUrl := env.Get("GCLOUD_STORAGE_CDN")

	return &FileUploaderImpl{
		client:     client,
		projectId:  projectId,
		bucketName: bucketName,
		cloudUrl:   cloudUrl,
	}
}

func (f *FileUploaderImpl) SendToBucket(ctx context.Context, file multipart.File, path, name string) error {
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

func (f *FileUploaderImpl) Upload(ctx context.Context, fileHeader any, folderName string) (string, error) {
	parsedFileHeader := fileHeader.(*multipart.FileHeader)
	file, err := parsedFileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	extension := filepath.Ext(parsedFileHeader.Filename)
	updateUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s%s", updateUUID.String(), extension)

	err = f.SendToBucket(ctx, file, fmt.Sprintf("%s/", folderName), fileName)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/%s/%s", f.cloudUrl, folderName, fileName)
	return url, nil
}
