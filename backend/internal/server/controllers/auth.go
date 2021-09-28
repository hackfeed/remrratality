package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hackfeed/remrratality/backend/internal/server/models"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
)

// SignUp godoc
// @Summary Signing user up
// @Description Signing user up by adding him to the database
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ResponseSuccessAuth
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Param request body models.User true "User's email and password"
// @Router /auth/signup [post]
func SignUp(c *gin.Context) {
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get user_repo",
		})
		return
	}

	var req models.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "Failed to parse request body",
		})
		return
	}

	existingUser, _ := userRepo.GetUser(*req.Email)
	if existingUser.Email != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "This email is already taken",
		})
		return
	}

	user, err := userRepo.AddUser(*req.Email, *req.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
		})
		return
	}

	expiresAt, err := user.GetExpirationTime()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get token expiration time",
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseSuccessAuth{
		Message:   "User created",
		IDToken:   *user.Token,
		LocalID:   user.UserID,
		ExpiresAt: expiresAt,
	})
}

// Login godoc
// @Summary Logging user in
// @Description Logging user in by retrieving his data from the database
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ResponseSuccessAuth
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Param request body models.User true "User's email and password"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get user_repo",
		})
		return
	}

	var req models.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "Failed to parse request body",
		})
		return
	}

	user, err := userRepo.GetUser(*req.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "User with given email doesn't exist",
		})
		return
	}

	err = user.VerifyPassword(*req.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Password is incorrect",
		})
		return
	}

	token, refreshToken, err := user.GenerateTokens()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to generate tokens",
		})
		return
	}

	user.Token = &token
	user.RefreshToken = &refreshToken

	err = userRepo.UpdateUser(user.UserID, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to update user data",
		})
		return
	}

	expiresAt, err := user.GetExpirationTime()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get token expiration time",
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseSuccessAuth{
		Message:   "Login success",
		IDToken:   *user.Token,
		LocalID:   user.UserID,
		ExpiresAt: expiresAt,
	})
}
