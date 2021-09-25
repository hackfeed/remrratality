package main

import (
	"github.com/hackfeed/remrratality/backend/internal/server"
)

func main() {
	server.SetupServer().Run()
}
