package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hackfeed/remrratality/backend/internal/server/models"
	log "github.com/sirupsen/logrus"
)

type signedDetails struct {
	Email  string
	UserID string
	jwt.StandardClaims
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			log.Errorf("failed to get authorization header")
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Message: "No Authorization header provided",
			})
			return
		}

		claims, err := validateToken(clientToken)
		if err != nil {
			log.Errorf("failed to validate %s token", clientToken)
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Message: "Token validation failed",
			})
			return
		}

		c.Set("email", claims.Email)
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func validateToken(signedToken string) (*signedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&signedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain token, error is: %s", err)
	}

	claims, ok := token.Claims.(*signedDetails)
	if !ok {
		return nil, errors.New("token is invalid")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}
