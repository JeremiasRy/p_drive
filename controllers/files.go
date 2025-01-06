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

type FileViewData struct {
	File      *model.FileMetaData
	Uploading bool
}

func NewFileController(service *services.FileService, ms *services.MetadataService, db *sql.DB) *FileController {
	return &FileController{service: service, ms: ms, db: db}
}

func (fc *FileController) HandleFiles(w http.ResponseWriter, r *http.Request, u *model.Users) {

	switch r.Method {
	case http.MethodGet:
		{
			path := strings.TrimPrefix(r.URL.Path, "/files/")
			split := strings.Split(path, "/")

			if split[0] == "folder" {
				fc.handleGetFiles(w, r, u, split[len(split)-1])
			} else if split[0] == "poll" {
				fc.handleGetPollFile(w, r, u, split[len(split)-1])
			} else {
				fc.handleGetFile(w, r, u, split[0])
			}

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

func (fc *FileController) handleGetPollFile(w http.ResponseWriter, r *http.Request, u *model.Users, fileId string) {
	file := fc.ms.GetFileById(fileId)
	tmpl, err := template.ParseFiles(filepath.Join("views", "templates", "file", "file-partials.html"))

	if err != nil {
		log.Printf("Failed to parse template file %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return

	}

	tmpl.ExecuteTemplate(w, "file", FileViewData{File: file, Uploading: *file.Status == model.FileStatus_Uploading})
}

func (fc *FileController) handleGetFile(w http.ResponseWriter, r *http.Request, u *model.Users, fileId string) {
	log.Printf("Get file by id %s", fileId)
	file := fc.ms.GetFileById(fileId)
	tmpl, err := template.ParseFiles(filepath.Join("views", "templates", "layout.html"), filepath.Join("views", "templates", "file", "file.html"), filepath.Join("views", "templates", "file", "file-partials.html"))

	if err != nil {
		log.Printf("Failed to parse template file %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", FileViewData{File: file, Uploading: *file.Status == model.FileStatus_Uploading})
}

func (fc *FileController) handleGetFiles(w http.ResponseWriter, r *http.Request, u *model.Users, folder string) {
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

	mime := handler.Header.Get("Content-Type")
	folder := r.PostFormValue("folder-path")
	name := strings.Join([]string{folder, handler.Filename}, "/")
	size := handler.Size

	metadata, err := fc.ms.InsertNewMetadata(handler.Filename, folder, mime, size)

	if err != nil {
		log.Printf("Failed to save metadata for file %s, %v", name, err)
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}

	go fc.service.UploadFile(context.WithoutCancel(r.Context()), file, name, mime, metadata.ID.String())

	tmpl, err := template.ParseFiles(filepath.Join("views", "templates", "file", "file-partials.html"))

	if err != nil {
		log.Printf("Failed to execute template %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	folders := fc.ms.GetFilesFromFolder(folder)
	tmpl.ExecuteTemplate(w, "file-list", folders)
}
