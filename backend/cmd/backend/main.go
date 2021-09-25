package main

import (
	"log"

	"github.com/hackfeed/remrratality/backend/internal/server"
)

func main() {
	log.Fatalln(server.SetupServer().Run())
}
