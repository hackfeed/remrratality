package domain

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type File struct {
	Name       string    `json:"name"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type User struct {
	UserID       string
	Email        *string
	Password     *string
	Token        *string
	RefreshToken *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Files        []File
}

type signedDetails struct {
	Email  string
	UserID string
	jwt.StandardClaims
}

func (u *User) GenerateTokens() (string, string, error) {
	var token, refreshToken string

	claims := &signedDetails{
		Email:  *u.Email,
		UserID: u.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}

	refreshClaims := &signedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(4)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return token, refreshToken, err
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))

	return token, refreshToken, err
}

func (u *User) UpdateTokens(signedToken, signedRefreshToken string) {
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	u.Token = &signedToken
	u.RefreshToken = &signedRefreshToken
	u.UpdatedAt = updatedAt
}

func (u *User) GetExpirationTime() (int64, error) {
	token, err := jwt.ParseWithClaims(
		*u.Token,
		&signedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*signedDetails)
	if !ok {
		return 0, errors.New("token is invalid")
	}

	return claims.ExpiresAt, nil
}

func (u *User) HashPassword() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*u.Password), 14)

	return string(bytes), err
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
}
