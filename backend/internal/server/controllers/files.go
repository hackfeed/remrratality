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
	"github.com/hackfeed/remrratality/backend/internal/server/models"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
)

func LoadFiles(c *gin.Context) {
	email, ok := c.MustGet("email").(string)
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

	files, err := loadFiles(userRepo, email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch user files",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files are loaded",
		"files":   files,
	})
}

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

	var req models.File

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	err = deleteFile(userRepo, email, userID, req.Name)
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
	filepath := fmt.Sprintf("%v/%v", dir, filename)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
	}

	err = c.SaveUploadedFile(file, filepath)
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

	csvFile, err := os.Open(filepath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to find data with given id",
		})
		return
	}
	defer csvFile.Close()

	invoices := []*models.Invoice{}

	err = gocsv.UnmarshalFile(csvFile, &invoices)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse given CSV file",
		})
		return
	}

	err = uploadFile()
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

func loadFiles(userRepo userrepo.UserRepository, email string) ([]map[string]interface{}, error) {
	user, err := userRepo.GetUser(email)

	return user.Files, err
}

func deleteFile(userRepo userrepo.UserRepository, email, userID, filename string) error {
	err := os.Remove(fmt.Sprintf("static/%v/%v", userID, filename))
	if err != nil {
		return err
	}

	user, err := userRepo.GetUser(email)
	if err != nil {
		return err
	}

	newFiles := []map[string]interface{}{}
	for _, file := range user.Files {
		if file["name"] != filename {
			newFiles = append(newFiles, file)
		}
	}

	user.Files = newFiles

	return userRepo.UpdateUser(userID, user)
}

func updateFiles(userRepo userrepo.UserRepository, email, userID, filename string) error {
	user, err := userRepo.GetUser(email)
	if err != nil {
		return err
	}

	uploadedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Files = append(user.Files, map[string]interface{}{"name": filename, "uploaded_at": uploadedAt})

	return userRepo.UpdateUser(userID, user)
}

func uploadFile() error {
	return nil
}
