package middlewares

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "No Authorization header provided",
			})
			return
		}

		claims, err := validateToken(clientToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Token validation failed",
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
		return nil, err
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
