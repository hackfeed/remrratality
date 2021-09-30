package userrepo

import (
	"fmt"
	"time"

	"github.com/hackfeed/remrratality/backend/internal/db/user"
	"github.com/hackfeed/remrratality/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoRepo struct {
	userClient user.MongoClient
}

func NewMongoRepo(userClient user.MongoClient) UserRepository {
	return &mongoRepo{
		userClient: userClient,
	}
}

func (mr *mongoRepo) AddUser(email, password string) (domain.User, error) {
	internalUser := domain.User{
		Email:    &email,
		Password: &password,
	}

	hashedPassword, err := internalUser.HashPassword()
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password for email %s, error is: %s", email, err)
	}
	createdAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	mappedUser := user.User{}
	mappedUser.ID = primitive.NewObjectID()
	mappedUser.UserID = mappedUser.ID.Hex()
	internalUser.UserID = mappedUser.UserID
	mappedUser.Email = &email
	mappedUser.Password = &hashedPassword
	token, refreshToken, _ := internalUser.GenerateTokens()
	mappedUser.Token = &token
	mappedUser.RefreshToken = &refreshToken
	mappedUser.CreatedAt = createdAt
	mappedUser.UpdatedAt = updatedAt
	mappedUser.Files = []user.File{}

	_, err = mr.userClient.Create(mappedUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to insert user with email %s, error is: %s", email, err)
	}

	internalUser.Password = mappedUser.Password
	internalUser.Token = mappedUser.Token
	internalUser.RefreshToken = mappedUser.RefreshToken
	internalUser.CreatedAt = mappedUser.CreatedAt
	internalUser.UpdatedAt = mappedUser.UpdatedAt
	internalFiles := convertFilesToDomain(mappedUser.Files)
	internalUser.Files = internalFiles

	return internalUser, nil
}

func (mr *mongoRepo) GetUser(email string) (domain.User, error) {
	user, err := mr.userClient.Read("email", email)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user with email %s, error is: %s", email, err)
	}

	mappedFiles := convertFilesToDomain(user.Files)

	mappedUser := domain.User{
		UserID:       user.UserID,
		Email:        user.Email,
		Password:     user.Password,
		Token:        user.Token,
		RefreshToken: user.RefreshToken,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Files:        mappedFiles,
	}

	return mappedUser, nil
}

func (mr *mongoRepo) UpdateUser(user_id string, user domain.User) error {
	updatedUser := primitive.D{
		bson.E{"user_id", user.UserID},
		bson.E{"email", user.Email},
		bson.E{"password", user.Password},
		bson.E{"token", user.Token},
		bson.E{"refresh_token", user.RefreshToken},
		bson.E{"created_at", user.CreatedAt},
		bson.E{"updated_at", user.UpdatedAt},
		bson.E{"files", convertFilesToUser(user.Files)},
	}
	return mr.userClient.Update(updatedUser, "user_id", user_id)
}

func convertFilesToDomain(userFiles []user.File) []domain.File {
	convertedFiles := []domain.File{}
	for _, file := range userFiles {
		convertedFiles = append(convertedFiles, domain.File{Name: file.Name, UploadedAt: file.UploadedAt})
	}
	return convertedFiles
}

func convertFilesToUser(domainFiles []domain.File) []user.File {
	convertedFiles := []user.File{}
	for _, file := range domainFiles {
		convertedFiles = append(convertedFiles, user.File{Name: file.Name, UploadedAt: file.UploadedAt})
	}
	return convertedFiles
}
