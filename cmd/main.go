package main

import (
	"github.com/0LuigiCode0/Library/logger"
	"github.com/nergilz/tmpserver/database"
	"github.com/nergilz/tmpserver/server"
)

func main() {
	log := logger.InitLogger("")
	log.Service("main start")
	configDB := database.New(log)
	server := server.New(configDB)
	err := server.Start()
	if err != nil {
		log.Fatalf("abandon server : %v", err)
	}
}
