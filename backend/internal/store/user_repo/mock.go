package userrepo

import (
	"errors"

	"github.com/hackfeed/remrratality/backend/internal/utils/user_validation"

	"github.com/hackfeed/remrratality/backend/internal/domain"
)

type UserRepositoryMock struct{}

func (urm *UserRepositoryMock) AddUser(email, password string) (domain.User, error) {
	if email == "errorGetUser" {
		return domain.User{}, errors.New("error while adding user")
	}
	if email == "someEmail" {
		id := "id"
		token, refreshToken, _ := user_validation.GenerateTokens(email, id)
		hashedPassword, _ := user_validation.HashPassword(password)
		return domain.User{
			Email:        email,
			Password:     hashedPassword,
			Token:        token,
			RefreshToken: refreshToken,
			Files:        make([]domain.File, 10),
		}, nil
	}
	return domain.User{
		Email:    email,
		Password: password,
		Files:    make([]domain.File, 10),
	}, nil
}

func (urm *UserRepositoryMock) GetUser(email string) (domain.User, error) {
	if email == "errorGetUser" {
		return domain.User{}, errors.New("user not exist")
	}
	if email == "errorToken" || email == "someEmail" {
		return domain.User{}, nil
	}
	id := "id"
	token, refreshToken, _ := user_validation.GenerateTokens(email, id)
	hashedPassword, _ := user_validation.HashPassword("somePass")
	return domain.User{
		Email:        email,
		Password:     hashedPassword,
		Token:        token,
		RefreshToken: refreshToken,
		Files:        make([]domain.File, 10),
	}, nil
}

func (urm *UserRepositoryMock) UpdateUser(userID string, _ domain.User) error {
	if userID == "errorUpdateUser" {
		return errors.New("error while updating user")
	}
	return nil
}
