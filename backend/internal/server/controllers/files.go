package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/hackfeed/remrratality/backend/internal/domain"
	"github.com/hackfeed/remrratality/backend/internal/server/models"
	storagerepo "github.com/hackfeed/remrratality/backend/internal/store/storage_repo"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
	log "github.com/sirupsen/logrus"
)

type Invoice struct {
	CustomerID  uint32  `csv:"customer_id"`
	PeriodStart string  `csv:"period_start"`
	PaidPlan    string  `csv:"paid_plan"`
	PaidAmount  float32 `csv:"paid_amount"`
	PeriodEnd   string  `csv:"period_end"`
}

// LoadFiles godoc
// @Summary Loading user's files
// @Description Loading files' names, uploaded by user
// @Tags files
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ResponseSuccessLoadFiles
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Router /files/load [get]
func LoadFiles(c *gin.Context) {
	email, ok := c.MustGet("email").(string)
	if !ok {
		log.Errorf("failed to get email from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Unable to determine logged in user",
		})
		return
	}
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		log.Errorf("failed to get user_repo from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get user_repo",
		})
		return
	}

	files, err := loadFiles(userRepo, email)
	if err != nil {
		log.Errorf("failed to load files for email %s, error is: %s", email, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to fetch user files",
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseSuccessLoadFiles{
		Message: "Files are loaded",
		Files:   files,
	})
}

// DeleteFile godoc
// @Summary Deleting user's file
// @Description Deleting file and cleaning database
// @Tags files
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Param filename path string true "File to delete"
// @Router /files/delete/{filename} [delete]
func DeleteFile(c *gin.Context) {
	email, ok := c.MustGet("email").(string)
	if !ok {
		log.Errorf("failed to get email from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Unable to determine logged in user",
		})
		return
	}
	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		log.Errorf("failed to get user_id from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Unable to determine logged in user",
		})
		return
	}
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		log.Errorf("failed to get user_repo from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get user_repo",
		})
		return
	}
	storageRepo, ok := c.MustGet("storage_repo").(storagerepo.StorageRepository)
	if !ok {
		log.Errorf("failed to get storage_repo from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get storage_repo",
		})
		return
	}

	filename := c.Param("filename")

	if err := deleteFile(userRepo, storageRepo, email, userID, filename); err != nil {
		log.Errorf("failed to delete file %s for email %s, user_id %s, error is: %s", filename, email, userID)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "Failed to delete file",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "File deleted",
	})
}

// SaveFile godoc
// @Summary Saving user's file
// @Description Saving file locally on the server and parsing its content to database
// @Tags files
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ResponseSuccessSaveFile
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Param file formData file true "File to upload"
// @Router /files/upload [post]
func SaveFile(c *gin.Context) {
	email, ok := c.MustGet("email").(string)
	if !ok {
		log.Errorf("failed to get email from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Unable to determine logged in user",
		})
		return
	}
	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		log.Errorf("failed to get user_id from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Unable to determine logged in user",
		})
		return
	}
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		log.Errorf("failed to get user_repo from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get user_repo",
		})
		return
	}
	storageRepo, ok := c.MustGet("storage_repo").(storagerepo.StorageRepository)
	if !ok {
		log.Errorf("failed to get storage_repo from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get storage_repo",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("failed to get file from formFile")
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "No file is received",
		})
		return
	}

	fext := filepath.Ext(file.Filename)
	if fext != ".csv" {
		log.Errorf("non csv files are not allowed, given %s", fext)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "Wrong file format. Please provide CSV file",
		})
		return
	}

	dir := fmt.Sprintf("static/%v", userID)
	filename := fmt.Sprintf("%v%v", uuid.New(), fext)
	filepth := fmt.Sprintf("%v/%v", dir, filename)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		log.Infof("%s isn't exist, creating", dir)
		os.Mkdir(dir, 0777)
	}

	if err = c.SaveUploadedFile(file, filepth); err != nil {
		log.Errorf("unable to save file %s, path %s, error is: %s", file, filepth, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Unable to save the file",
		})
		return
	}

	if err = updateFiles(userRepo, email, userID, filename); err != nil {
		log.Errorf("unable to update files for email %s, user_id %s, error is: %s", email, userID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Unable to update user in db",
		})
		return
	}

	csvFile, err := os.Open(filepth)
	if err != nil {
		log.Errorf("unable to open file at %s, error is: %s", filepth, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to find data with given id",
		})
		return
	}
	defer csvFile.Close()

	var invoices []*Invoice

	if err = gocsv.UnmarshalFile(csvFile, &invoices); err != nil {
		log.Errorf("unable to unmarshal %s, error is: %s", csvFile, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "Failed to parse given CSV file",
		})
		return
	}

	if err = uploadFile(storageRepo, userID, filename, invoices); err != nil {
		log.Errorf("unable to upload invoices for email %s, user_id %s, error is: %s", email, userID, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "Failed to upload data to database",
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseSuccessSaveFile{
		Message:  "File is uploaded",
		Filename: filename,
	})
}

func loadFiles(userRepo userrepo.UserRepository, email string) ([]domain.File, error) {
	user, err := userRepo.GetUser(email)
	if err != nil {
		return []domain.File{}, fmt.Errorf("failed to get user, error is: %s", err)
	}

	return user.Files, nil
}

func deleteFile(userRepo userrepo.UserRepository, storageRepo storagerepo.StorageRepository, email, userID, filename string) error {
	if err := os.Remove(fmt.Sprintf("static/%v/%v", userID, filename)); err != nil {
		return fmt.Errorf("failed to remove local file, error is: %s", err)
	}

	user, err := userRepo.GetUser(email)
	if err != nil {
		return fmt.Errorf("failed to get user, error is: %s", err)
	}

	newFiles := make([]domain.File, 0)
	for _, file := range user.Files {
		if file.Name != filename {
			newFiles = append(newFiles, file)
		}
	}

	user.Files = newFiles
	if err = userRepo.UpdateUser(userID, user); err != nil {
		return fmt.Errorf("failed to update user, error is: %s", err)
	}

	if err = storageRepo.DeleteInvoices(userID, filename); err != nil {
		return fmt.Errorf("failed to delete ivoices from db, error is: %s", err)
	}

	return nil
}

func updateFiles(userRepo userrepo.UserRepository, email, userID, filename string) error {
	user, err := userRepo.GetUser(email)
	if err != nil {
		return fmt.Errorf("failed to get user, error is: %s", err)
	}

	uploadedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Files = append(user.Files, domain.File{Name: filename, UploadedAt: uploadedAt})

	if err = userRepo.UpdateUser(userID, user); err != nil {
		return fmt.Errorf("failed to update user, error is: %s", err)
	}

	return nil
}

func uploadFile(storageRepo storagerepo.StorageRepository, userID, fileID string, invoices []*Invoice) error {
	mappedInvoices := make([]domain.Invoice, 0)

	for _, invoice := range invoices {
		mappedInvoice := domain.Invoice{
			UserID:      userID,
			FileID:      fileID,
			CustomerID:  invoice.CustomerID,
			PeriodStart: invoice.PeriodStart,
			PaidPlan:    invoice.PaidPlan,
			PaidAmount:  invoice.PaidAmount,
			PeriodEnd:   invoice.PeriodEnd,
		}
		mappedInvoices = append(mappedInvoices, mappedInvoice)
	}

	if _, err := storageRepo.AddInvoices(mappedInvoices); err != nil {
		return fmt.Errorf("failed upload invoices to db, error is: %s", err)
	}

	return nil
}
