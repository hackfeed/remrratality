package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hackfeed/remrratality/backend/internal/server/models"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
)

func SignUp(c *gin.Context) {
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user_repo",
		})
		return
	}

	var req models.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	existingUser, _ := userRepo.GetUser(*req.Email)
	if existingUser.Email != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "This email is already taken",
		})
		return
	}

	user, err := userRepo.AddUser(*req.Email, *req.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	expiresAt, err := user.GetExpirationTime()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get token expiration time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "User created",
		"idToken":   *user.Token,
		"localId":   user.UserID,
		"expiresAt": expiresAt,
	})
}

func Login(c *gin.Context) {
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user_repo",
		})
		return
	}

	var req models.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	user, err := userRepo.GetUser(*req.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "User with given email doesn't exist",
		})
		return
	}

	err = user.VerifyPassword(*req.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Password is incorrect",
		})
		return
	}

	token, refreshToken, err := user.GenerateTokens()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate tokens",
		})
		return
	}

	user.Token = &token
	user.RefreshToken = &refreshToken

	err = userRepo.UpdateUser(user.UserID, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user data",
		})
		return
	}

	expiresAt, err := user.GetExpirationTime()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get token expiration time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login success",
		"idToken":   *user.Token,
		"localId":   user.UserID,
		"expiresAt": expiresAt,
	})
}
