package cloudstorage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

const bucketName = "dmd-we-care"

var gcsLogger = log.New(os.Stdout, "[GCS] ", log.LstdFlags)

type ICloudStorageService interface {
	UploadImage(file *multipart.FileHeader) (*string, error)
}

type service struct {
	client *storage.Client
}

func NewService(gcsClient *storage.Client) *service {
	return &service{gcsClient}
}

func (s *service) UploadImage(file *multipart.FileHeader) (*string, error) {
	src, err := file.Open()
	if err != nil {
		gcsLogger.Println("error opening file")
		return nil, errors.New("error opening file")
	}
	defer src.Close()

	filename := fmt.Sprintf("uploads/%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))

	ctx := context.Background()
	storageWriter := s.client.Bucket(bucketName).Object(filename).NewWriter(ctx)
	storageWriter.ContentType = file.Header.Get("Content-Type") // Set the correct content type

	if _, err := io.Copy(storageWriter, src); err != nil {
		er := fmt.Errorf("error writing to GCS: %v", err)
		gcsLogger.Println(er.Error())
		return nil, er
	}
	if err := storageWriter.Close(); err != nil {
		er := fmt.Errorf("error closing GCS writer %v", err)
		gcsLogger.Println(er.Error())
		return nil, er
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
	return &publicURL, nil
}
