package main

import (
	"github.com/hackfeed/remrratality/backend/internal/server"
)

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic("Failed to load .env file")
	// }
	server.SetupServer().Run(":8081")
}
