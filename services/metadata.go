package services

import (
	"backend/database"
	"context"
	"database/sql"
	"log"
)

type MetadataService struct {
	db *sql.DB
	md *database.MetadataDb
}

func NewMetaDataService(db *sql.DB, md *database.MetadataDb) *MetadataService {
	return &MetadataService{db: db, md: md}
}

func (ms *MetadataService) InsertNewMetadata(ctx context.Context, metadata database.NewMetadata) error {
	db, err := ms.db.BeginTx(ctx, nil)
	defer db.Commit().Error()

	if err != nil {
		log.Printf("Failed to start transaction %v", err)
		return err
	}

	ms.md.InsertMetadata(metadata, db)
	return nil
}
