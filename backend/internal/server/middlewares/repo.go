package middlewares

import (
	"github.com/gin-gonic/gin"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
)

func UserRepo(userRepo userrepo.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_repo", userRepo)
		c.Next()
	}
}
