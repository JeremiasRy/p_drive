package controllers

import (
	"backend/.gen/personal_drive/public/model"
	"backend/services"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
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
			fc.handleGetFiles(w, r, u)
		}
	case http.MethodPost:
		{
			fc.handlePostUpload(w, r, u)
		}
	default:
		{
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}

}

func (fc *FileController) handleGetFiles(w http.ResponseWriter, r *http.Request, u *model.Users) {
	folderPath := strings.Split(r.URL.Path, "/")
	folder := folderPath[len(folderPath)-1]
	files := fc.ms.GetFilesFromFolder(folder)

	for _, file := range files {
		fc.service.GetFilesSignedLink(r.Context(), file)
		log.Printf("File URL: %s\n", *file.SignedLink)
	}

	tmpl, err := template.ParseFiles(filepath.Join("views", "templates", "file", "file-partials.html"))

	if err != nil {
		log.Printf("Failed to parse template file %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "file-list", files)
}

func (fc *FileController) handlePostUpload(w http.ResponseWriter, r *http.Request, u *model.Users) {
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
