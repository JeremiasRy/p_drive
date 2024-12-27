package database

import (
	"backend/.gen/personal_drive/public/table"
	"database/sql"
)

type UsersDb struct {
	db *sql.DB
}

func NewUsersDb(db *sql.DB) *UsersDb {
	return &UsersDb{db: db}
}

func (db *UsersDb) InsertNewUser(email string) error {
	stmt := table.Users.INSERT(table.Users.Email).VALUES(email)
	_, err := stmt.Exec(db.db)
	return err
}
