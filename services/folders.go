package services

import (
	"backend/.gen/personal_drive/public/model"
	"backend/database"
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type FoldersService struct {
	db *sql.DB
	fd *database.FoldersDatabase
}

func NewFoldersService(db *sql.DB, fd *database.FoldersDatabase) *FoldersService {
	return &FoldersService{db: db, fd: fd}
}

func (fs *FoldersService) CreateNewFolder(name string, parent uuid.UUID) error {
	err := fs.fd.CreateSubFolder(name, parent, fs.db)

	if err != nil {
		log.Printf("Failed to create sub folder %v\n", err)
		return err
	}

	return nil
}

func (fs *FoldersService) GetFoldersFromNode(id string) []model.Folders {
	return fs.fd.GetFoldersFromNode(id, fs.db)
}
