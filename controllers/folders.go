package controllers

import (
	"backend/.gen/personal_drive/public/model"
	"backend/services"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type FoldersController struct {
	fs *services.FoldersService
}

func NewFoldersController(fs *services.FoldersService) *FoldersController {
	return &FoldersController{fs: fs}
}

func (fc *FoldersController) HandleFolders(w http.ResponseWriter, r *http.Request, u *model.Users) {
	switch r.Method {
	case http.MethodGet:
		{
			fc.handleGetFolders(w, u)
		}
	case http.MethodPost:
		{
			fc.handlePostFolders(w, r, u)
		}
	default:
		{
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (fc *FoldersController) handlePostFolders(w http.ResponseWriter, r *http.Request, u *model.Users) {
	err := r.ParseForm()

	if err != nil {
		log.Printf("Failed to parse form data %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	newFolderName := r.PostFormValue("folder-name")
	err = fc.fs.CreateNewFolder(newFolderName, u.ID)

	if err != nil {
		log.Printf("Failed to create new folder %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	folders := fc.fs.GetFoldersFromNode(u.ID.String())

	template, err := template.ParseFiles(filepath.Join("views", "templates", "folder-list.html"))

	if err != nil {
		log.Printf("Failed to parse template HTML %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(w, "folder-list", folders)

	if err != nil {
		log.Printf("Dailed to execute HTML template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (fc *FoldersController) handleGetFolders(w http.ResponseWriter, u *model.Users) {
	template, err := template.ParseFiles(filepath.Join("views", "templates", "folders.html"), filepath.Join("views", "templates", "folder-list.html"), filepath.Join("views", "templates", "folder-input.html"))

	if err != nil {
		log.Printf("Failed to load HTML template %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	folders := fc.fs.GetFoldersFromNode(u.ID.String())

	err = template.ExecuteTemplate(w, "folders", folders)

	if err != nil {
		log.Printf("Failed to execute HTML template %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
