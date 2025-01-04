package database

import (
	"github.com/go-jet/jet/v2/qrm"
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

func (md *MetadataDb) InsertMetadata(newMetaData NewMetadata, db qrm.Executable) {

}

func (md *MetadataDb) GetFilesFromFolder(folder string, db qrm.Queryable) {

}
