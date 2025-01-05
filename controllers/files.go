package controllers

import (
	"backend/.gen/personal_drive/public/model"
	"backend/services"
	"context"
	"database/sql"
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
	if err != nil {
		log.Fatalf("Payload error: %s", err)
		http.Error(w, "Failed to validate payload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	mime := handler.Header.Get("Content-Type")
	folder := r.PostFormValue("folder-path")
	name := strings.Join([]string{folder, handler.Filename}, "/")
	size := handler.Size

	err = fc.ms.InsertNewMetadata(handler.Filename, folder, mime, size)

	if err != nil {
		log.Printf("Failed to save metadata for file %s, %v", name, err)
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}

	go fc.service.UploadFile(context.WithoutCancel(r.Context()), file, name, mime)

	tmpl, err := template.ParseFiles(filepath.Join("views", "templates", "file", "file-partials.html"))

	if err != nil {
		log.Printf("Failed to execute template %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	folders := fc.ms.GetFilesFromFolder(folder)
	tmpl.ExecuteTemplate(w, "file-list", folders)
}
