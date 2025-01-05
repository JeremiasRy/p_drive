package database

import (
	"backend/.gen/personal_drive/public/model"
	"backend/.gen/personal_drive/public/table"
	"log"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type MetadataDb struct {
}

func NewMetaDataDb() *MetadataDb {
	return &MetadataDb{}
}

func (md *MetadataDb) InsertMetadata(folder string, name string, mime string, size_bytes int64, db qrm.Executable) error {
	stmt := table.FileMetaData.INSERT(table.FileMetaData.FolderID, table.FileMetaData.Name, table.FileMetaData.Mime, table.FileMetaData.SizeBytes).VALUES(folder, name, mime, size_bytes)
	_, err := stmt.Exec(db)

	return err
}

func (md *MetadataDb) GetFilesFromFolder(folder string, db qrm.Queryable) []*model.FileMetaData {
	results := []*model.FileMetaData{}

	stmt := table.FileMetaData.SELECT(table.FileMetaData.AllColumns).WHERE(table.FileMetaData.FolderID.EQ(postgres.UUID(uuid.MustParse(folder))))

	err := stmt.Query(db, &results)

	if err != nil {
		log.Printf("failed to query for files from database: %v\n", err)
	}

	return results
}
