package user_validation

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/hackfeed/remrratality/backend/internal/domain"

	"github.com/stretchr/testify/assert"
)

var (
	user      domain.User
	realToken string
)

func TestMain(m *testing.M) {
	user = domain.User{
		UserID:    "1",
		Email:     "test@test.com",
		Password:  "pass",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Files:     make([]domain.File, 0),
	}

	hashedPassword, _ := HashPassword(user.Password)
	user.Password = hashedPassword

	token, refreshToken, _ := GenerateTokens(user.Email, user.UserID)
	realToken = token
	user.Token = realToken
	user.RefreshToken = refreshToken

	os.Exit(m.Run())
}

func TestGenerateTokens(t *testing.T) {
	token, refreshToken, err := GenerateTokens(user.Email, user.Password)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotNil(t, refreshToken)
}

func TestUpdateTokens(t *testing.T) {
	type testInput struct {
		token, refreshToken string
	}

	tests := []struct {
		input testInput
	}{
		{
			input: testInput{
				token:        "fakeToken",
				refreshToken: "fakeRefreshToken",
			},
		},
	}

	for _, test := range tests {
		oldToken := user.Token
		oldRefreshToken := user.RefreshToken
		oldUpdatedAt := user.UpdatedAt
		UpdateTokens(&user, test.input.token, test.input.refreshToken)
		assert.NotEqual(t, oldToken, user.Token)
		assert.NotEqual(t, oldRefreshToken, user.RefreshToken)
		assert.NotEqual(t, oldUpdatedAt, user.UpdatedAt)
	}
}

func TestGetExpirationTime(t *testing.T) {
	type testInput struct {
		token string
	}
	type testWant struct {
		err error
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				token: "fakeToken",
			},
			want: testWant{
				err: errors.New("failed to get token, error is: token contains an invalid number of segments"),
			},
		},
		{
			input: testInput{
				token: realToken,
			},
			want: testWant{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		_, err := GetExpirationTime(test.input.token)
		assert.Equal(t, test.want.err, err)
	}
}

func TestHashPassword(t *testing.T) {
	type testInput struct {
		password string
	}
	type testWant struct {
		err error
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				password: "newPass",
			},
			want: testWant{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		_, err := HashPassword(test.input.password)
		assert.Equal(t, test.want.err, err)
	}
}

func TestVerifyPassword(t *testing.T) {
	type testInput struct {
		password string
	}
	type testWant struct {
		err error
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				password: "notUserPass",
			},
			want: testWant{
				err: errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password"),
			},
		},
	}

	for _, test := range tests {
		err := VerifyPassword(user.Password, test.input.password)
		assert.Equal(t, test.want.err, err)
	}
}
