package middlewares

import (
	"github.com/gin-gonic/gin"
	storagerepo "github.com/hackfeed/remrratality/backend/internal/store/storage_repo"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
)

func UserRepo(userRepo userrepo.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_repo", userRepo)
		c.Next()
	}
}

func StorageRepo(storageRepo storagerepo.StorageRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("storage_repo", storageRepo)
		c.Next()
	}
}
