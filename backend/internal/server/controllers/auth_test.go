package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hackfeed/remrratality/backend/internal/server/models"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
	internalTesting "github.com/hackfeed/remrratality/backend/internal/utils/testing"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	type testInput struct {
		keys map[string]interface{}
		body interface{}
	}
	type testWant struct {
		code    int
		message string
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{keys: map[string]interface{}{
				"user_repo": "invalidRepo",
			}},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "{\"message\":\"Failed to get user_repo\"}",
			},
		},
		{
			input: testInput{keys: map[string]interface{}{
				"user_repo": &userrepo.UserRepositoryMock{},
			}},
			want: testWant{
				code:    http.StatusBadRequest,
				message: "{\"message\":\"Failed to parse request body\"}",
			},
		},
		{
			input: testInput{
				keys: map[string]interface{}{
					"user_repo": &userrepo.UserRepositoryMock{},
				},
				body: models.User{
					Email:    "takenEmail",
					Password: "somePass",
				},
			},
			want: testWant{
				code:    http.StatusBadRequest,
				message: "{\"message\":\"This email is already taken\"}",
			},
		},
		{
			input: testInput{
				keys: map[string]interface{}{
					"user_repo": &userrepo.UserRepositoryMock{},
				},
				body: models.User{
					Email:    "errorGetUser",
					Password: "somePass",
				},
			},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "{\"message\":\"Failed to create new user. Please, try again later\"}",
			},
		},
		{
			input: testInput{
				keys: map[string]interface{}{
					"user_repo": &userrepo.UserRepositoryMock{},
				},
				body: models.User{
					Email:    "someEmail",
					Password: "somePass",
				},
			},
			want: testWant{
				code:    http.StatusOK,
				message: "User created",
			},
		},
	}

	for _, test := range tests {
		c, w := internalTesting.CreateGinContext(test.input.keys, test.input.body, nil)
		SignUp(c)
		assert.Equal(t, test.want.code, w.Code)
		assert.Equal(t, true, strings.Contains(w.Body.String(), test.want.message))
	}
}

func TestLogin(t *testing.T) {
	type testInput struct {
		keys map[string]interface{}
		body interface{}
	}
	type testWant struct {
		code    int
		message string
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{keys: map[string]interface{}{
				"user_repo": "invalidRepo",
			}},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "{\"message\":\"Failed to get user_repo\"}",
			},
		},
		{
			input: testInput{keys: map[string]interface{}{
				"user_repo": &userrepo.UserRepositoryMock{},
			}},
			want: testWant{
				code:    http.StatusBadRequest,
				message: "{\"message\":\"Failed to parse request body\"}",
			},
		},
		{
			input: testInput{
				keys: map[string]interface{}{
					"user_repo": &userrepo.UserRepositoryMock{},
				},
				body: models.User{
					Email:    "errorGetUser",
					Password: "somePass",
				},
			},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "{\"message\":\"User with given email doesn't exist\"}",
			},
		},
		{
			input: testInput{
				keys: map[string]interface{}{
					"user_repo": &userrepo.UserRepositoryMock{},
				},
				body: models.User{
					Email:    "userWithWrongPass",
					Password: "fakePass",
				},
			},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "{\"message\":\"Password is incorrect\"}",
			},
		},
		{
			input: testInput{
				keys: map[string]interface{}{
					"user_repo": &userrepo.UserRepositoryMock{},
				},
				body: models.User{
					Email:    "verifiedEmail",
					Password: "somePass",
				},
			},
			want: testWant{
				code:    http.StatusOK,
				message: "Login success",
			},
		},
	}

	for _, test := range tests {
		c, w := internalTesting.CreateGinContext(test.input.keys, test.input.body, nil)
		Login(c)
		fmt.Println(test)
		assert.Equal(t, test.want.code, w.Code)
		assert.Equal(t, true, strings.Contains(w.Body.String(), test.want.message))
	}
}
