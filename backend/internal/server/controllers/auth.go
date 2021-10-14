package controllers

import (
	"net/http"

	"github.com/hackfeed/remrratality/backend/internal/utils/user_validation"

	"github.com/gin-gonic/gin"
	"github.com/hackfeed/remrratality/backend/internal/server/models"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
	log "github.com/sirupsen/logrus"
)

// SignUp godoc
// @Summary Signing user up
// @Description Signing user up by adding him to the database
// @Tags signup
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ResponseSuccessAuth
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Param request body models.User true "User's email and password"
// @Router /signup [post]
func SignUp(c *gin.Context) {
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		log.Errorf("failed to get user_repo from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get user_repo",
		})
		return
	}

	var req models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("failed to parse request body, error is: %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "Failed to parse request body",
		})
		return
	}

	existingUser, _ := userRepo.GetUser(req.Email)
	if existingUser.Email != "" {
		log.Infof("user with email %s already exists", existingUser.Email)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "This email is already taken",
		})
		return
	}

	user, err := userRepo.AddUser(req.Email, req.Password)
	if err != nil {
		log.Errorf("failed to add user with email %s, error is: %s", req.Email, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to create new user. Please, try again later",
		})
		return
	}

	expiresAt, err := user_validation.GetExpirationTime(user.Token)
	if err != nil {
		log.Errorf("failed to get token expiration time for user %s, error is: %s", user.UserID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get token expiration time",
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseSuccessAuth{
		Message:   "User created",
		IDToken:   user.Token,
		LocalID:   user.UserID,
		ExpiresAt: expiresAt,
	})
}

// Login godoc
// @Summary Logging user in
// @Description Logging user in by retrieving his data from the database
// @Tags login
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ResponseSuccessAuth
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Param request body models.User true "User's email and password"
// @Router /login [post]
func Login(c *gin.Context) {
	userRepo, ok := c.MustGet("user_repo").(userrepo.UserRepository)
	if !ok {
		log.Errorf("failed to get user_repo from gin.Context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get user_repo",
		})
		return
	}

	var req models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("failed to parse request body, error is: %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "Failed to parse request body",
		})
		return
	}

	user, err := userRepo.GetUser(req.Email)
	if err != nil {
		log.Errorf("failed to get user with email %s, error is: %s", req.Email, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "User with given email doesn't exist",
		})
		return
	}

	if err = user_validation.VerifyPassword(user.Password, req.Password); err != nil {
		log.Errorf("failed to verify password for user %s, error is: %s", user.UserID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Password is incorrect",
		})
		return
	}

	token, refreshToken, err := user_validation.GenerateTokens(user.Email, user.UserID)
	if err != nil {
		log.Errorf("failed to generate tokens for user %s, error is: %s", user.UserID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to generate tokens",
		})
		return
	}

	user_validation.UpdateTokens(&user, token, refreshToken)

	if err = userRepo.UpdateUser(user.UserID, user); err != nil {
		log.Errorf("failed to update user %s, error is: %s", user.UserID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to update user data",
		})
		return
	}

	expiresAt, err := user_validation.GetExpirationTime(user.Token)
	if err != nil {
		log.Errorf("failed to get token expiration time for user %s, error is: %s", user.UserID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get token expiration time",
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseSuccessAuth{
		Message:   "Login success",
		IDToken:   user.Token,
		LocalID:   user.UserID,
		ExpiresAt: expiresAt,
	})
}
