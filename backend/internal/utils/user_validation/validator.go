package user_validation

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hackfeed/remrratality/backend/internal/domain"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type signedDetails struct {
	Email  string
	UserID string
	jwt.StandardClaims
}

func GenerateTokens(email, id string) (string, string, error) {
	var token, refreshToken string

	claims := &signedDetails{
		Email:  email,
		UserID: id,
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
		return token, refreshToken, fmt.Errorf("failed create new token, error is: %s", err)
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return token, refreshToken, fmt.Errorf("failed create new refresh token, error is: %s", err)
	}

	return token, refreshToken, nil
}

func UpdateTokens(user *domain.User, signedToken, signedRefreshToken string) {
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	user.Token = signedToken
	user.RefreshToken = signedRefreshToken
	user.UpdatedAt = updatedAt
}

func GetExpirationTime(token string) (int64, error) {
	tk, err := jwt.ParseWithClaims(
		token,
		&signedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get token, error is: %s", err)
	}
	claims, ok := tk.Claims.(*signedDetails)
	if !ok {
		return 0, errors.New("token is invalid")
	}

	return claims.ExpiresAt, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("failed to hash password, error is: %s", err)
	}

	return string(bytes), nil
}

func VerifyPassword(hashed, given string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(given))
}
