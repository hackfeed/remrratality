package controllers

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/hackfeed/remrratality/backend/internal/domain"
	storagerepo "github.com/hackfeed/remrratality/backend/internal/store/storage_repo"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
	internalTesting "github.com/hackfeed/remrratality/backend/internal/utils/testing"
	"github.com/stretchr/testify/assert"
)

func TestLoadFilesHandler(t *testing.T) {
	type testInput struct {
		keys map[string]interface{}
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
				"email": 1,
			}},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "Unable to determine logged in user",
			},
		},
		{
			input: testInput{keys: map[string]interface{}{
				"email":     "test@test.com",
				"user_repo": "invalidRepo",
			}},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "Failed to get user_repo",
			},
		},
		{
			input: testInput{keys: map[string]interface{}{
				"email":     "errorGetUser",
				"user_repo": &userrepo.UserRepositoryMock{},
			}},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "Failed to fetch user files",
			},
		},
		{
			input: testInput{keys: map[string]interface{}{
				"email":     "test@test.com",
				"user_repo": &userrepo.UserRepositoryMock{},
			}},
			want: testWant{
				code:    http.StatusOK,
				message: "Files are loaded",
			},
		},
	}

	for _, test := range tests {
		c, w := internalTesting.CreateGinContext(test.input.keys, nil, nil)
		LoadFiles(c)
		assert.Equal(t, test.want.code, w.Code)
		assert.Equal(t, true, strings.Contains(w.Body.String(), test.want.message))
	}
}

func TestDeleteFileContentHandler(t *testing.T) {
	type testInput struct {
		keys map[string]interface{}
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
				"email": 1,
			}},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "Unable to determine logged in user",
			},
		},
		{
			input: testInput{keys: map[string]interface{}{
				"email":   "test@test.com",
				"user_id": 1,
			}},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "Unable to determine logged in user",
			},
		},
		{
			input: testInput{keys: map[string]interface{}{
				"email":     "test@test.com",
				"user_id":   "id",
				"user_repo": "invalidRepo",
			}},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "Failed to get user_repo",
			},
		},
		{
			input: testInput{keys: map[string]interface{}{
				"email":        "test@test.com",
				"user_id":      "id",
				"user_repo":    &userrepo.UserRepositoryMock{},
				"storage_repo": "invalidRepo",
			}},
			want: testWant{
				code:    http.StatusInternalServerError,
				message: "Failed to get storage_repo",
			},
		},
		{
			input: testInput{keys: map[string]interface{}{
				"email":        "errorGetUser",
				"user_id":      "id",
				"user_repo":    &userrepo.UserRepositoryMock{},
				"storage_repo": &storagerepo.StorageRepositoryMock{},
			}},
			want: testWant{
				code:    http.StatusBadRequest,
				message: "Failed to delete file",
			},
		},
		{
			input: testInput{keys: map[string]interface{}{
				"email":        "user",
				"user_id":      "user",
				"user_repo":    &userrepo.UserRepositoryMock{},
				"storage_repo": &storagerepo.StorageRepositoryMock{},
			}},
			want: testWant{
				code:    http.StatusOK,
				message: "File deleted",
			},
		},
	}

	for _, test := range tests {
		c, w := internalTesting.CreateGinContext(test.input.keys, nil, nil)
		DeleteFileContent(c)
		assert.Equal(t, test.want.code, w.Code)
		assert.Equal(t, true, strings.Contains(w.Body.String(), test.want.message))
	}
}

func TestLoadFiles(t *testing.T) {
	type testInput struct {
		email string
	}
	type testWant struct {
		files []domain.File
		err   error
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				email: "errorGetUser",
			},
			want: testWant{
				files: make([]domain.File, 0),
				err:   errors.New("failed to get user, error is: user not exist"),
			},
		},
		{
			input: testInput{
				email: "user",
			},
			want: testWant{
				files: make([]domain.File, 10),
			},
		},
	}

	userMock := &userrepo.UserRepositoryMock{}

	for _, test := range tests {
		files, err := loadFiles(userMock, test.input.email)
		assert.Equal(t, test.want.files, files)
		assert.Equal(t, test.want.err, err)
	}
}

func TestDeleteFileContent(t *testing.T) {
	type testInput struct {
		email, userID string
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
				email: "errorGetUser",
			},
			want: testWant{
				err: errors.New("failed to get user, error is: user not exist"),
			},
		},
		{
			input: testInput{
				email:  "user",
				userID: "errorUpdateUser",
			},
			want: testWant{
				err: errors.New("failed to update user, error is: error while updating user"),
			},
		},
		{
			input: testInput{
				email:  "user",
				userID: "errorDeleteInvoices",
			},
			want: testWant{
				err: errors.New("failed to delete ivoices from db, error is: error while deleting invoices"),
			},
		},
		{
			input: testInput{
				email:  "user",
				userID: "user",
			},
			want: testWant{
				err: nil,
			},
		},
	}

	userMock := &userrepo.UserRepositoryMock{}
	storageMock := &storagerepo.StorageRepositoryMock{}

	for _, test := range tests {
		err := deleteFileContent(userMock, storageMock, test.input.email, test.input.userID, "someFile")
		assert.Equal(t, test.want.err, err)
	}
}

func TestUpdateFiles(t *testing.T) {
	type testInput struct {
		email, userID string
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
				email: "errorGetUser",
			},
			want: testWant{
				err: errors.New("failed to get user, error is: user not exist"),
			},
		},
		{
			input: testInput{
				email:  "user",
				userID: "errorUpdateUser",
			},
			want: testWant{
				err: errors.New("failed to update user, error is: error while updating user"),
			},
		},
		{
			input: testInput{
				email:  "user",
				userID: "user",
			},
			want: testWant{
				err: nil,
			},
		},
	}

	userMock := &userrepo.UserRepositoryMock{}

	for _, test := range tests {
		err := updateFiles(userMock, test.input.email, test.input.userID, "someFile")
		assert.Equal(t, test.want.err, err)
	}
}

func TestUploadFileContent(t *testing.T) {
	type testInput struct {
		invoices []*Invoice
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
				invoices: make([]*Invoice, 0),
			},
			want: testWant{
				err: errors.New("failed upload invoices to db, error is: error while adding invoices"),
			},
		},
		{
			input: testInput{
				invoices: []*Invoice{{}},
			},
			want: testWant{
				err: nil,
			},
		},
	}

	storageMock := &storagerepo.StorageRepositoryMock{}

	for _, test := range tests {
		err := uploadFileContent(storageMock, "", "", test.input.invoices)
		assert.Equal(t, test.want.err, err)
	}
}
