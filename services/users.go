package services

import (
	"backend/.gen/personal_drive/public/model"
	"backend/database"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type UserService struct {
	db *sql.DB
	ud *database.UsersDb
	fd *database.FoldersDatabase
}

func NewUserService(db *sql.DB, ud *database.UsersDb, fd *database.FoldersDatabase) *UserService {
	return &UserService{db: db, ud: ud, fd: fd}
}

func (us *UserService) NewUser(ctx context.Context, email string) error {
	tx, err := us.db.BeginTx(ctx, nil)

	defer tx.Rollback()

	if err != nil {
		log.Printf("Failed to start transaction: %v\n", err)
		return err
	}

	err = us.ud.InsertNewUser(email, tx)

	if err != nil {
		fmt.Printf("Failed to insert new user: %v\n", err)
		tx.Rollback()
		return err
	}

	user, err := us.ud.GetUserByEmail(email, tx)

	if err != nil {
		fmt.Printf("Failed to retrieve new user %v\n", err)
		tx.Rollback()
		return err
	}

	err = us.fd.CreateRootFolder(user.ID.String(), tx)

	if err != nil {
		fmt.Printf("Failed to create root folder: %v\n", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (us *UserService) GetUserByEmail(email string) (*model.Users, error) {
	return us.ud.GetUserByEmail(email, us.db)
}

func (us *UserService) GetUserByID(id string) (*model.Users, error) {
	return us.ud.GetUserByID(id, us.db)
}
