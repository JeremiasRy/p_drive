package database

import (
	"backend/.gen/personal_drive/public/model"
	"backend/.gen/personal_drive/public/table"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type UsersDb struct{}

func NewUsersDb() *UsersDb {
	return &UsersDb{}
}

func (ud *UsersDb) InsertNewUser(email string, db qrm.Executable) (int64, error) {
	stmt := table.Users.INSERT(table.Users.Email).VALUES(email).ON_CONFLICT(table.Users.Email).DO_NOTHING()
	res, err := stmt.Exec(db)

	if err != nil {
		return -1, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

func (ud *UsersDb) GetUserByEmail(email string, db qrm.Queryable) (*model.Users, error) {
	var user model.Users
	stmt := table.Users.SELECT(table.Users.AllColumns).WHERE(table.Users.Email.EQ(postgres.String(email)))

	err := stmt.Query(db, &user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (fds *UsersDb) GetUserByID(id string, db qrm.Queryable) (*model.Users, error) {
	var user model.Users
	stmt := table.Users.SELECT(table.Users.AllColumns).WHERE(table.Users.ID.EQ(postgres.UUID(uuid.MustParse(id))))

	err := stmt.Query(db, &user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
