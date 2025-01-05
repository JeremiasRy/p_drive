package services

import (
	"context"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"strings"
	"time"

	"backend/.gen/personal_drive/public/model"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileService struct {
	client *minio.Client
}

var (
	MINIO_URL        = os.Getenv("MINIO_URL")
	MINIO_ACCESS_KEY = os.Getenv("MINIO_ROOT_USER")
	MINIO_SECRET     = os.Getenv("MINIO_ROOT_PASSWORD")
	BUCKET_NAME      = os.Getenv("MINIO_BUCKET")
	REGION           = os.Getenv("MINIO_REGION")
)

func NewFileservice() (*FileService, error) {
	client, err := minio.New(MINIO_URL, &minio.Options{
		Creds: credentials.NewStaticV4(MINIO_ACCESS_KEY, MINIO_SECRET, ""),
	})

	if err != nil {
		return nil, err
	}

	wait := time.Second
	retries := 0
	for {
		if retries > 5 {
			log.Fatal("Failed to connect to MINIO to check bucket")
		}
		exists, errBucketExists := client.BucketExists(context.Background(), BUCKET_NAME)

		if errBucketExists != nil {
			log.Printf("Failed to fecth bucket infromation: %v\n", errBucketExists)
			log.Printf("Retrying in %s seconds...", wait)
			retries++
			time.Sleep(wait)
			wait = wait << 1
			continue
		}

		if exists {
			log.Println("Succesfully created file server client")
			return &FileService{client: client}, nil
		}
		break
	}

	log.Printf("Bucket %s does not exist, creating it...", BUCKET_NAME)
	errCreateBucket := client.MakeBucket(context.Background(), BUCKET_NAME, minio.MakeBucketOptions{Region: REGION})
	if errCreateBucket != nil {
		log.Fatalf("Failed to create bucket: %v\n", errCreateBucket)
	}

	log.Printf("Successfully created bucket %s", BUCKET_NAME)
	log.Println("Succesfully created file server client")
	return &FileService{client: client}, nil
}

func (fs *FileService) UploadFile(ctx context.Context, file multipart.File, name string, contentType string) (minio.UploadInfo, error) {
	client := fs.client

	info, err := client.PutObject(ctx, BUCKET_NAME, name, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return minio.UploadInfo{}, err
	}

	log.Printf("Successfully uploaded %s of size %d\n", name, info.Size)
	return info, nil
}

func (fs *FileService) GetFilesSignedLink(ctx context.Context, file *model.FileMetaData) {
	client := fs.client

	reqParams := url.Values{}
	reqParams.Add("response-content-type", file.Mime)

	objectKey := strings.Join([]string{file.FolderID.String(), file.Name}, "/")

	link, err := client.PresignedGetObject(ctx, BUCKET_NAME, objectKey, time.Hour, reqParams)

	if err != nil {
		log.Printf("Failed to generate presigned URL for file: %v, ERROR: %s", file, err)
	}

	url := link.String()
	file.SignedLink = &url

}
