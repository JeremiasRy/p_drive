package controllers

import (
	"backend/.gen/personal_drive/public/model"
	"backend/services"
	"fmt"
	"log"
	"net/http"
)

type FileController struct {
	service *services.FileService
}

func NewFileController(service *services.FileService) *FileController {
	return &FileController{service: service}
}

func (fc *FileController) HandlePostUpload(w http.ResponseWriter, r *http.Request, u *model.Users) {
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

	err = fc.service.UploadFile(r.Context(), file, handler.Filename, handler.Header.Get("Content-Type"))
	if err != nil {
		fmt.Fprintf(w, "Failed to upload file to file server: %s", err)
	}
}
