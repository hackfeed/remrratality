package domain

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	user      User
	realToken string
)

func TestMain(m *testing.M) {
	email := "test@test.com"
	password := "pass"

	user = User{
		UserID:    "1",
		Email:     &email,
		Password:  &password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Files:     make([]File, 0),
	}

	hashedPassword, _ := user.HashPassword()
	user.Password = &hashedPassword

	token, refreshToken, _ := user.GenerateTokens()
	realToken = token
	user.Token = &realToken
	user.RefreshToken = &refreshToken

	os.Exit(m.Run())
}

func TestGenerateTokens(t *testing.T) {
	token, refreshToken, err := user.GenerateTokens()
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
		user.UpdateTokens(test.input.token, test.input.refreshToken)
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
		user.Token = &test.input.token
		_, err := user.GetExpirationTime()
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
		user.Password = &test.input.password
		_, err := user.HashPassword()
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
				err: errors.New("crypto/bcrypt: hashedSecret too short to be a bcrypted password"),
			},
		},
	}

	for _, test := range tests {
		err := user.VerifyPassword(test.input.password)
		assert.Equal(t, test.want.err, err)
	}
}
