package main

import (
	"os"

	"github.com/hackfeed/remrratality/backend/internal/server"
	log "github.com/sirupsen/logrus"
)

// @title remrratality API
// @version 1.0
// @description API for getting MRR analytics of your app's money flow.

// @contact.name Sergey "hackfeed" Kononenko
// @contact.url https://hackfeed.github.io
// @contact.email hackfeed@yandex.ru

// @license.name GPL-3.0 License
// @license.url http://www.gnu.org/licenses/gpl-3.0.html

// @host weblabs.com:8003
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token

func main() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to create logs file, error is: %s", err)
	}

	log.SetOutput(file)

	log.Fatalln(server.SetupServer().Run())
}
