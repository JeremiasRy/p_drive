package main

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIoClient struct {
	client *minio.Client
}

func (m *MinIoClient) uploadFile(ctx context.Context, file multipart.File, name string, contentType string) error {
	bucketName := "test"
	location := "test"

	err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := m.client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own hello%s\n", bucketName)
		} else {
			return err
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	info, err := m.client.PutObject(ctx, "my-bucker", name, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", name, info.Size)
	return nil
}

type TemplateData struct {
	BackendURL string
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	templatePath := filepath.Join("static", "index.html")
	tmpl, err := template.ParseFiles(templatePath)
	d := TemplateData{BackendURL: os.Getenv("BACKEND_URL")}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	client, err := minio.New(os.Getenv("MINIO_URL"), &minio.Options{
		Creds: credentials.NewStaticV4(os.Getenv("MINIO_ROOT_USER"), os.Getenv("MINIO_ROOT_PASSWORD"), ""),
	})

	if err != nil {
		log.Fatalf("Failed to start file server client, %s", err)
	}

	minIoClient := MinIoClient{client: client}
	http.HandleFunc("/", templateHandler)
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		err := r.ParseMultipartForm(10 << 24)

		if err != nil {
			http.Error(w, "Failed to validate payload", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")

		if err != nil {
			http.Error(w, "Failed to validate payload", http.StatusBadRequest)
			return
		}
		defer file.Close()

		err = minIoClient.uploadFile(r.Context(), file, handler.Filename, handler.Header.Get("Content-Type"))
		if err != nil {
			fmt.Fprintf(w, "Somethings up: %s", err)
		}
	})

	log.Println("Starting a server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
