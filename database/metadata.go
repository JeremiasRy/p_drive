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

func (md *MetadataDb) InsertMetadata(folder string, name string, mime string, size_bytes int64, db qrm.Queryable) (*model.FileMetaData, error) {
	stmt := table.FileMetaData.INSERT(table.FileMetaData.FolderID, table.FileMetaData.Name, table.FileMetaData.Mime, table.FileMetaData.SizeBytes).VALUES(folder, name, mime, size_bytes).RETURNING(table.FileMetaData.AllColumns)
	var returnValue model.FileMetaData
	err := stmt.Query(db, &returnValue)

	return &returnValue, err
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

func (md *MetadataDb) GetFileById(id string, db qrm.Queryable) *model.FileMetaData {
	var result model.FileMetaData

	stmt := table.FileMetaData.SELECT(table.FileMetaData.AllColumns).WHERE(table.FileMetaData.ID.EQ(postgres.UUID(uuid.MustParse(id))))
	stmt.Query(db, &result)
	return &result
}

func (md *MetadataDb) UpdateFileStatus(id string, status model.FileStatus, db qrm.Executable) {
	stmt := table.FileMetaData.UPDATE(table.FileMetaData.Status).SET(status).WHERE(table.FileMetaData.ID.EQ(postgres.UUID(uuid.MustParse(id))))
	stmt.Exec(db)
}
