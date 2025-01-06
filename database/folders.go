package database

import (
	"database/sql"
	"log"

	"backend/.gen/personal_drive/public/model"
	"backend/.gen/personal_drive/public/table"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type FoldersDatabase struct{}

func (fd *FoldersDatabase) GetBreadcrumbsFromFolder(id string, db *sql.DB) []model.Folders {
	sub := postgres.CTE("sub")
	stmt := postgres.WITH_RECURSIVE(
		sub.AS(
			postgres.SELECT(
				table.Folders.AllColumns,
			).FROM(
				table.Folders,
			).WHERE(
				table.Folders.ID.EQ(postgres.UUID(uuid.MustParse(id)))).UNION(
				postgres.SELECT(
					table.Folders.AllColumns,
				).FROM(
					table.Folders.INNER_JOIN(sub, table.Folders.ParentID.From(sub).EQ(table.Folders.ID)),
				),
			),
		),
	)(
		postgres.SELECT(
			sub.AllColumns(),
		).FROM(
			sub,
		),
	)

	results := []model.Folders{}

	stmt.Query(db, &results)

	return results
}

func NewFoldersDatabase() *FoldersDatabase {
	return &FoldersDatabase{}
}

func (fd *FoldersDatabase) CreateRootFolder(userId string, db qrm.Executable) error {
	stmt := table.Folders.INSERT(table.Folders.ID, table.Folders.Name, table.Folders.FolderClientPath).VALUES(userId, "My Drive", "my-drive")
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

func (fd *FoldersDatabase) GetFoldersFromNode(id string, db qrm.Queryable) []model.Folders {
	folders := []model.Folders{}

	stmt := table.Folders.SELECT(table.Folders.AllColumns).WHERE(table.Folders.ParentID.EQ(postgres.UUID(uuid.MustParse(id))))

	err := stmt.Query(db, &folders)

	if err != nil {
		log.Printf("Failed to fetch folders, id: %s, %v", id, err)
	}
	return folders
}
