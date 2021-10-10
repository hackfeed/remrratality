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

func (mr *mongoRepo) UpdateUser(userID string, user domain.User) error {
	updatedUser := primitive.D{
		bson.E{Key: "user_id", Value: user.UserID},
		bson.E{Key: "email", Value: user.Email},
		bson.E{Key: "password", Value: user.Password},
		bson.E{Key: "token", Value: user.Token},
		bson.E{Key: "refresh_token", Value: user.RefreshToken},
		bson.E{Key: "created_at", Value: user.CreatedAt},
		bson.E{Key: "updated_at", Value: user.UpdatedAt},
		bson.E{Key: "files", Value: convertFilesToUser(user.Files)},
	}
	return mr.userClient.Update(updatedUser, "user_id", userID)
}

func convertFilesToDomain(userFiles []user.File) []domain.File {
	convertedFiles := make([]domain.File, len(userFiles))
	for i, file := range userFiles {
		convertedFiles[i] = domain.File{Name: file.Name, UploadedAt: file.UploadedAt}
	}
	return convertedFiles
}

func convertFilesToUser(domainFiles []domain.File) []user.File {
	convertedFiles := make([]user.File, len(domainFiles))
	for i, file := range domainFiles {
		convertedFiles[i] = user.File{Name: file.Name, UploadedAt: file.UploadedAt}
	}
	return convertedFiles
}
