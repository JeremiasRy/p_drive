package database

import (
	"backend/.gen/personal_drive/public/table"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type MetadataDb struct {
}

type NewMetadata struct {
	Folder    string
	Name      string
	Mime      string
	SizeBytes int64
}

func NewMetaDataDb() *MetadataDb {
	return &MetadataDb{}
}

func (md *MetadataDb) InsertMetadata(folder uuid.UUID, name string, mime string, size_bytes int64, db qrm.Executable) error {
	stmt := table.FileMetaData.INSERT(table.FileMetaData.FolderID, table.FileMetaData.Name, table.FileMetaData.Mime, table.FileMetaData.SizeBytes).VALUES(folder, name, mime, size_bytes)
	_, err := stmt.Exec(db)

	return err
}

func (md *MetadataDb) GetFilesFromFolder(folder string, db qrm.Queryable) {

}
