package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type FoldersDatabase struct {
	db *sql.DB
}

func NewFoldersDatabase(db *sql.DB) *FoldersDatabase {
	return &FoldersDatabase{db: db}
}

func (fd *FoldersDatabase) CreateFolder(name string, parent *uuid.UUID) {

}
