package userrepo

import (
	"errors"

	"github.com/hackfeed/remrratality/backend/internal/domain"
)

type UserRepositoryMock struct{}

func (urm *UserRepositoryMock) AddUser(email, password string) (domain.User, error) {
	return domain.User{Email: &email, Password: &password, Files: make([]domain.File, 10)}, nil
}

func (urm *UserRepositoryMock) GetUser(email string) (domain.User, error) {
	if email == "errorGetUser" {
		return domain.User{}, errors.New("user not exist")
	}
	return domain.User{Email: &email, Files: make([]domain.File, 10)}, nil
}

func (urm *UserRepositoryMock) UpdateUser(userID string, _ domain.User) error {
	if userID == "errorUpdateUser" {
		return errors.New("error while updating user")
	}
	return nil
}
