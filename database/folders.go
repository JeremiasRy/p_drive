package database

import (
	"log"

	"backend/.gen/personal_drive/public/table"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type FoldersDatabase struct{}

func NewFoldersDatabase() *FoldersDatabase {
	return &FoldersDatabase{}
}

func (fd *FoldersDatabase) CreateRootFolder(userId string, db qrm.Executable) error {
	stmt := table.Folders.INSERT(table.Folders.Name).VALUES(userId)
	_, err := stmt.Exec(db)

	if err != nil {
		log.Printf("Failed to insert root folder: %v\n", err)
		return err
	}
	return nil
}

func (fd *FoldersDatabase) CreateSubFolder(name string, parent uuid.UUID, db qrm.Executable) error {
	stmt := table.Folders.INSERT(table.Folders.Name, table.Folders.ParentID).VALUES(name, parent)

	_, err := stmt.Exec(db)

	if err != nil {
		log.Printf("Failed to create sub folder folder: %v\n", err)
		return err
	}
	return nil
}
