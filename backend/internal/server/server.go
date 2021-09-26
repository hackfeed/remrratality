package server

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hackfeed/remrratality/backend/internal/db/cache"
	"github.com/hackfeed/remrratality/backend/internal/db/storage"
	"github.com/hackfeed/remrratality/backend/internal/db/user"
	"github.com/hackfeed/remrratality/backend/internal/server/controllers"
	"github.com/hackfeed/remrratality/backend/internal/server/middlewares"
	cacherepo "github.com/hackfeed/remrratality/backend/internal/store/cache_repo"
	storagerepo "github.com/hackfeed/remrratality/backend/internal/store/storage_repo"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
	"github.com/joho/godotenv"
)

var (
	userRepo    userrepo.UserRepository
	storageRepo storagerepo.StorageRepository
	cacheRepo   cacherepo.CacheRepository
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Failed to load .env file")
	}

	ctx := context.Background()

	userClient, err := user.NewMongoClient(ctx, &user.Options{
		Host:     os.Getenv("MONGO_HOST"),
		Port:     os.Getenv("MONGO_PORT"),
		User:     os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASS"),
		DB:       os.Getenv("MONGO_DB"),
	})
	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln(err)
	}
	storageRepo = storagerepo.NewPostgresRepo(*storageClient)

	cacheClient, err := cache.NewRedisClient(ctx, &cache.Options{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       os.Getenv("REDIS_DB"),
	})
	if err != nil {
		log.Fatalln(err)
	}
	cacheRepo = cacherepo.NewRedisRepo(*cacheClient, 1*time.Hour)
}

func SetupServer() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"token"}
	r.Use(cors.New(config))

	r.Use(middlewares.UserRepo(userRepo))
	r.Use(middlewares.StorageRepo(storageRepo))
	r.Use(middlewares.CacheRepo(cacheRepo))

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	r.Use(middlewares.Auth())

	files := r.Group("/files")
	{
		files.GET("/load", controllers.LoadFiles)
		files.POST("/upload", controllers.SaveFile)
		files.POST("/delete", controllers.DeleteFile)
	}

	r.POST("/analytics", controllers.GetAnalytics)

	return r
}
