package server

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hackfeed/remrratality/backend/internal/db/storage"
	"github.com/hackfeed/remrratality/backend/internal/db/user"
	"github.com/hackfeed/remrratality/backend/internal/server/controllers"
	"github.com/hackfeed/remrratality/backend/internal/server/middlewares"
	userrepo "github.com/hackfeed/remrratality/backend/internal/store/user_repo"
)

var (
	userRepo userrepo.UserRepository
)

func init() {
	ctx := context.Background()

	userClient, err := user.NewMongoClient(ctx, &user.Options{
		Host: "localhost",
		Port: 27017,
	})
	if err != nil {
		log.Fatalln(err)
	}
	userRepo = userrepo.NewMongoRepo(*userClient)

	_, err = storage.NewPostgresClient(ctx, &storage.Options{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "pass",
		DB:       "postgres",
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func SetupServer() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"token"}
	r.Use(cors.New(config))

	r.Use(middlewares.UserRepo(userRepo))

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	r.Use(middlewares.Auth())

	files := r.Group("/files")
	{
		files.GET("/load", controllers.LoadFiles)
		files.POST("/upload", controllers.SaveFile)
		files.POST("/delete", controllers.DeleteFile)
	}

	// r.POST("/analytics", controllers.GetAnalytics)

	return r
}
