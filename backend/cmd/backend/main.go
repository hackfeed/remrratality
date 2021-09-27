package main

import (
	"log"

	"github.com/hackfeed/remrratality/backend/internal/server"
)

// @title remrratality API
// @version 1.0
// @description API for getting MRR analytics of your app's money flow.

// @contact.name Sergey "hackfeed" Kononenko
// @contact.url https://hackfeed.github.io
// @contact.email hackfeed@vk.com

// @license.name GPL-3.0 License
// @license.url http://www.gnu.org/licenses/gpl-3.0.html

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token

func main() {
	log.Fatalln(server.SetupServer().Run())
}
