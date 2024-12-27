package database

import "database/sql"

type MetadataDb struct {
	db *sql.DB
}

func NewMetaDataDb(db *sql.DB) *MetadataDb {
	return &MetadataDb{db: db}
}
