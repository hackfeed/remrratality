package server

import (
	"context"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/hackfeed/remrratality/backend/docs"
	"github.com/hackfeed/remrratality/backend/internal/db/cache"
	"github.com/hackfeed/remrratality/backend/internal/db/storage"
	"github.com/hackfeed/remrratality/backend/internal/db/user"
	"github.com/hackfeed/remrratality/backend/internal/server/controllers"
	"github.com/hackfeed/remrratality/backend/internal/server/middlewares"
	cacherepo "github.com/hackfeed/remrratality/backend/internal/store/cache_repo"
	storagerepo "github.com/hackfeed/remrratality/backend/internal/store/storage_repo"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	userRepo    userrepo.UserRepository
	storageRepo storagerepo.StorageRepository
	cacheRepo   cacherepo.CacheRepository
)

func init() {
	ctx := context.Background()

	userClient, err := user.NewMongoClient(ctx, &user.Options{
		Host:     os.Getenv("MONGO_HOST"),
		Port:     os.Getenv("MONGO_PORT"),
		User:     os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASS"),
		DB:       os.Getenv("MONGO_DB"),
	})
	if err != nil {
		log.Fatalf("failed to create mongo client, error is: %s", err)
	}
	userRepo = userrepo.NewMongoRepo(*userClient)

	storageClient, err := storage.NewPostgresClient(ctx, &storage.Options{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASS"),
		DB:       os.Getenv("POSTGRES_DB"),
	})
	if err != nil {
		log.Fatalf("failed to create postgres client, error is: %s", err)
	}
	storageRepo = storagerepo.NewPostgresRepo(*storageClient)

	cacheClient, err := cache.NewRedisClient(ctx, &cache.Options{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       os.Getenv("REDIS_DB"),
	})
	if err != nil {
		log.Fatalf("failed to create redis client, error is: %s", err)
	}
	cacheRepo = cacherepo.NewRedisRepo(*cacheClient, 1*time.Hour)
}

func SetupServer() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"token", "Origin", "X-Requested-With", "Content-Type", "Accept"}
	config.AllowMethods = []string{"GET", "POST", "DELETE"}
	r.Use(cors.New(config))

	r.Use(middlewares.UserRepo(userRepo))
	r.Use(middlewares.StorageRepo(storageRepo))
	r.Use(middlewares.CacheRepo(cacheRepo))

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("auth")
		{
			auth.POST("/signup", controllers.SignUp)
			auth.POST("/login", controllers.Login)
		}
		files := v1.Group("/files", middlewares.Auth())
		{
			files.GET("/load", controllers.LoadFiles)
			files.POST("/upload", controllers.SaveFile)
			files.DELETE("/delete/:filename", controllers.DeleteFile)
		}
		analytics := v1.Group("analytics", middlewares.Auth())
		{
			analytics.POST("/mrr", controllers.GetAnalytics)
		}
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}
