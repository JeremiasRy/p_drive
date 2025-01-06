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

func (ms *MetadataService) InsertNewMetadata(name string, folder string, mime string, size int64) (*model.FileMetaData, error) {
	newMetadata, err := ms.md.InsertMetadata(folder, name, mime, size, ms.db)
	return newMetadata, err
}

func (ms *MetadataService) GetFilesFromFolder(folder string) []*model.FileMetaData {
	files := ms.md.GetFilesFromFolder(folder, ms.db)
	return files
}

func (ms *MetadataService) GetFileById(id string) *model.FileMetaData {
	file := ms.md.GetFileById(id, ms.db)
	return file
}
