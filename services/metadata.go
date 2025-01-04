package services

import (
	"backend/database"
	"database/sql"

	"github.com/google/uuid"
)

type MetadataService struct {
	db *sql.DB
	md *database.MetadataDb
}

func NewMetaDataService(db *sql.DB, md *database.MetadataDb) *MetadataService {
	return &MetadataService{db: db, md: md}
}

func (ms *MetadataService) InsertNewMetadata(name string, folder string, mime string, size int64) error {
	folderUUID := uuid.MustParse(folder)
	err := ms.md.InsertMetadata(folderUUID, name, mime, size, ms.db)
	return err
}
