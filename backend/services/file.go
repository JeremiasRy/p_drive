package services

import (
	"context"
	"log"
	"mime/multipart"
	"os"

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
	log.Println("Succesfully created file server client")
	return &FileService{client: client}, err
}

func (s *FileService) UploadFile(ctx context.Context, file multipart.File, name string, contentType string) error {
	client := s.client
	log.Printf("%v\n", BUCKET_NAME)

	bucketErr := s.checkBucket(ctx)

	if bucketErr != nil {
		return bucketErr
	}

	info, err := client.PutObject(ctx, BUCKET_NAME, name, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", name, info.Size)
	return nil
}

func (fs *FileService) checkBucket(ctx context.Context) error {
	client := fs.client
	exists, errBucketExists := client.BucketExists(ctx, BUCKET_NAME)

	if errBucketExists != nil {
		return errBucketExists
	}

	if exists {
		return nil
	}

	log.Printf("Bucket %s does not exist, creating it...", BUCKET_NAME)
	errCreateBucket := client.MakeBucket(ctx, BUCKET_NAME, minio.MakeBucketOptions{Region: REGION})
	if errCreateBucket != nil {
		return errCreateBucket
	}

	log.Printf("Successfully created bucket %s", BUCKET_NAME)
	return nil
}
