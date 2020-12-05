package main

import (
	"github.com/0LuigiCode0/Library/logger"
	"github.com/nergilz/tmpserver/database"
	"github.com/nergilz/tmpserver/server"
)

func main() {
	log := logger.InitLogger("")
	log.Service("start main")
	configDB := database.New(log)
	server := server.New(configDB)
	if err := server.Start(); err != nil {
		log.Fatalf("abandon server : %v", err)
	}
	log.Service("start server")
}
