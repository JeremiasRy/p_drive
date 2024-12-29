package database

import (
	"backend/.gen/personal_drive/public/model"
	"backend/.gen/personal_drive/public/table"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type UsersDb struct {
	db *sql.DB
}

func NewUsersDb(db *sql.DB) *UsersDb {
	return &UsersDb{db: db}
}

func (db *UsersDb) InsertNewUser(email string, authStrategy string) error {
	stmt := table.Users.INSERT(table.Users.Email, table.Users.AuthStrategy).VALUES(email, authStrategy)
	_, err := stmt.Exec(db.db)
	return err
}

func (db *UsersDb) GetUserByEmail(email string) (*model.Users, error) {
	var user model.Users
	stmt := table.Users.SELECT(table.Users.AllColumns).WHERE(table.Users.Email.EQ(postgres.String(email)))

	err := stmt.Query(db.db, &user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *UsersDb) GetUserByID(id string) (*model.Users, error) {
	var user model.Users
	stmt := table.Users.SELECT(table.Users.AllColumns).WHERE(table.Users.ID.EQ(postgres.UUID(uuid.MustParse(id))))

	err := stmt.Query(db.db, &user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
