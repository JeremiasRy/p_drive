package controllers

import (
	"backend/.gen/personal_drive/public/model"
	"backend/services"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type FileController struct {
	service *services.FileService
	ms      *services.MetadataService
	db      *sql.DB
}

func NewFileController(service *services.FileService, ms *services.MetadataService, db *sql.DB) *FileController {
	return &FileController{service: service, ms: ms, db: db}
}

func (fc *FileController) HandleFiles(w http.ResponseWriter, r *http.Request, u *model.Users) {

	switch r.Method {
	case http.MethodGet:
		{
			fc.HandleGetFiles(w, r, u)
		}
	case http.MethodPost:
		{
			fc.HandlePostUpload(w, r, u)
		}
	default:
		{
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}

}

func (fc *FileController) HandleGetFiles(w http.ResponseWriter, r *http.Request, u *model.Users) {
	folderPath := strings.Split(strings.TrimPrefix(r.URL.Path, "folders"), "/")
	folder := folderPath[len(folderPath)-1]

	log.Printf("Folder %s", folder)
}

func (fc *FileController) HandlePostUpload(w http.ResponseWriter, r *http.Request, u *model.Users) {
	err := r.ParseMultipartForm(10 << 24)

	if err != nil {
		log.Fatalf("Payload error: %s", err)
		http.Error(w, "Failed to validate payload", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	folder := r.PostFormValue("folder-path")

	name := strings.Join([]string{folder, handler.Filename}, "/")

	if err != nil {
		log.Fatalf("Payload error: %s", err)
		http.Error(w, "Failed to validate payload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	info, err := fc.service.UploadFile(r.Context(), file, name, handler.Header.Get("Content-Type"))

	if err != nil {
		fmt.Fprintf(w, "Failed to upload file to file server: %s", err)
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	err = fc.ms.InsertNewMetadata(handler.Filename, folder, handler.Header.Get("Content-Type"), info.Size)

	if err != nil {
		log.Printf("Failed to save metadata for file at:  %s", name)
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}
}
