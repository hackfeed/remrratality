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
// @Failure 500 {object} models.ResponseFailLoadFiles
// @Security ApiKeyAuth
// @Router /files/load [get]
func LoadFiles(c *gin.Context) {
	email, ok := c.MustGet("email").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ResponseFailLoadFiles{
			Message: "Unable to determine logged in user",
		})
		return
	}
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ResponseFailLoadFiles{
			Message: "Failed to get user_repo",
		})
		return
	}

	files, err := loadFiles(userRepo, email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ResponseFailLoadFiles{
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
// @Security ApiKeyAuth
// @Router /files/delete/{filename} [delete]
func DeleteFile(c *gin.Context) {
	email, ok := c.MustGet("email").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to determine logged in user",
		})
		return
	}
	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to determine logged in user",
		})
		return
	}
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user_repo",
		})
		return
	}
	storageRepo, ok := c.MustGet("storage_repo").(storagerepo.StorageRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get storage_repo",
		})
		return
	}

	filename := c.Param("filename")

	err := deleteFile(userRepo, storageRepo, email, userID, filename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to delete file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted",
	})
}

// SaveFile godoc
// @Summary Saving user's file
// @Description Saving file locally on the server and parsing its content to database
// @Tags files
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /files/upload [post]
func SaveFile(c *gin.Context) {
	email, ok := c.MustGet("email").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to determine logged in user",
		})
		return
	}
	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to determine logged in user",
		})
		return
	}
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user_repo",
		})
		return
	}
	storageRepo, ok := c.MustGet("storage_repo").(storagerepo.StorageRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get storage_repo",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	fext := filepath.Ext(file.Filename)
	if fext != ".csv" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Wrong file format. Please provide CSV file",
		})
		return
	}

	dir := fmt.Sprintf("static/%v", userID)
	filename := fmt.Sprintf("%v%v", uuid.New(), fext)
	filepth := fmt.Sprintf("%v/%v", dir, filename)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
	}

	err = c.SaveUploadedFile(file, filepth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	err = updateFiles(userRepo, email, userID, filename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to update user in db",
		})
		return
	}

	csvFile, err := os.Open(filepth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to find data with given id",
		})
		return
	}
	defer csvFile.Close()

	var invoices []*Invoice

	err = gocsv.UnmarshalFile(csvFile, &invoices)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse given CSV file",
		})
		return
	}

	err = uploadFile(storageRepo, userID, filename, invoices)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to upload data to database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded",
		"filename": filename,
	})
}

func loadFiles(userRepo userrepo.UserRepository, email string) ([]domain.File, error) {
	user, err := userRepo.GetUser(email)

	return user.Files, err
}

func deleteFile(userRepo userrepo.UserRepository, storageRepo storagerepo.StorageRepository, email, userID, filename string) error {
	err := os.Remove(fmt.Sprintf("static/%v/%v", userID, filename))
	if err != nil {
		return err
	}

	user, err := userRepo.GetUser(email)
	if err != nil {
		return err
	}

	newFiles := make([]domain.File, 0)
	for _, file := range user.Files {
		if file.Name != filename {
			newFiles = append(newFiles, file)
		}
	}

	user.Files = newFiles
	err = userRepo.UpdateUser(userID, user)
	if err != nil {
		return err
	}

	return storageRepo.DeleteInvoices(userID, filename)
}

func updateFiles(userRepo userrepo.UserRepository, email, userID, filename string) error {
	user, err := userRepo.GetUser(email)
	if err != nil {
		return err
	}

	uploadedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Files = append(user.Files, domain.File{Name: filename, UploadedAt: uploadedAt})

	return userRepo.UpdateUser(userID, user)
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

	_, err := storageRepo.AddInvoices(mappedInvoices)

	return err
}
