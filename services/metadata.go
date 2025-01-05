package services

import (
	"backend/.gen/personal_drive/public/model"
	"backend/database"
	"database/sql"
)

type MetadataService struct {
	db *sql.DB
	md *database.MetadataDb
}

func NewMetaDataService(db *sql.DB, md *database.MetadataDb) *MetadataService {
	return &MetadataService{db: db, md: md}
}

func (ms *MetadataService) InsertNewMetadata(name string, folder string, mime string, size int64) error {
	err := ms.md.InsertMetadata(folder, name, mime, size, ms.db)
	return err
}

func (ms *MetadataService) GetFilesFromFolder(folder string) []*model.FileMetaData {
	files := ms.md.GetFilesFromFolder(folder, ms.db)
	return files
}
