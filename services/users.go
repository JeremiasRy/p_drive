package services

import (
	"backend/.gen/personal_drive/public/model"
	"backend/database"
)

type UserService struct {
	db *database.UsersDb
}

func NewUserService(db *database.UsersDb) *UserService {
	return &UserService{db: db}
}

func (us *UserService) NewUser(email string) error {
	return us.db.InsertNewUser(email)
}

func (us *UserService) GetUserByEmail(email string) (*model.Users, error) {
	return us.db.GetUserByEmail(email)
}
