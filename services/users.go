package services

import "backend/database"

type UserService struct {
	db *database.UsersDb
}

func NewUserService(db *database.UsersDb) *UserService {
	return &UserService{db: db}
}

func (us *UserService) NewUser(email string) error {
	return us.db.InsertNewUser(email)
}
