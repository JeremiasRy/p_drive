package controllers

import (
	"backend/.gen/personal_drive/public/model"
	"backend/services"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FoldersController struct {
	service *services.FileService
	fs      *services.FoldersService
	ms      *services.MetadataService
}

func NewFoldersController(fs *services.FoldersService, ms *services.MetadataService, service *services.FileService) *FoldersController {
	return &FoldersController{fs: fs, ms: ms, service: service}
}

func (fc *FoldersController) HandleFolders(w http.ResponseWriter, r *http.Request, u *model.Users) {
	switch r.Method {
	case http.MethodGet:
		{
			path := strings.Split(strings.TrimPrefix(r.URL.Path, "/folders/"), "/")
			folderId := path[0]

			if path[len(path)-1] == "files" {
				fc.handleGetFolderFiles(w, r, u, folderId)
			} else {
				fc.handleGetFolders(w, r, u)
			}
		}
	case http.MethodPost:
		{
			fc.handlePostFolders(w, r)
		}
	default:
		{
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (fc *FoldersController) handleGetFolderFiles(w http.ResponseWriter, r *http.Request, u *model.Users, folder string) {
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

func (fc *FoldersController) handlePostFolders(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Printf("Failed to parse form data %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	newFolderName := r.PostFormValue("folder-name")
	newFolderParent := r.PostFormValue("folder-parent")
	err = fc.fs.CreateNewFolder(newFolderName, uuid.MustParse(newFolderParent))

	if err != nil {
		log.Printf("Failed to create new folder %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	folders := fc.fs.GetFoldersFromNode(newFolderParent)

	template, err := template.ParseFiles(filepath.Join("views", "templates", "folder", "folder-partials.html"))

	if err != nil {
		log.Printf("Failed to parse template HTML %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(w, "folder-list", struct{ Folders []model.Folders }{Folders: folders})

	if err != nil {
		log.Printf("Dailed to execute HTML template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (fc *FoldersController) handleGetFolders(w http.ResponseWriter, r *http.Request, u *model.Users) {
	template, err := template.ParseFiles(filepath.Join("views", "templates", "layout.html"), filepath.Join("views", "templates", "folder", "folder.html"), filepath.Join("views", "templates", "folder", "folder-partials.html"), filepath.Join("views", "templates", "file", "file-partials.html"))

	if err != nil {
		log.Printf("Failed to load HTML template %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s := strings.Split(r.URL.Path, "/")

	var folderNode string

	if s[len(s)-1] == "my-drive" {
		folderNode = u.ID.String()
	} else {
		folderNode = s[len(s)-1]
	}

	folders := fc.fs.GetFoldersFromNode(folderNode)
	breadcrumbs := fc.fs.GetBreadcrumbsFromNode(folderNode)

	err = template.ExecuteTemplate(w, "layout", struct {
		Breadcrumbs  []model.Folders
		Folders      []model.Folders
		Folder       string
		FolderParent string
	}{Folders: folders, Folder: folderNode, FolderParent: folderNode, Breadcrumbs: breadcrumbs})

	if err != nil {
		log.Printf("Failed to execute HTML template %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
