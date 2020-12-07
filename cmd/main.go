package main

import (
	"github.com/0LuigiCode0/Library/logger"
	"github.com/nergilz/tmpserver/database"
	"github.com/nergilz/tmpserver/server"
)

func main() {
	log := logger.InitLogger("")
	DBconfig := database.New(log)
	server := server.New(DBconfig, log)
	if err := server.Start(); err != nil {
		log.Fatalf("Server Fatal : %v", err)
	} else {
		log.Service("start server")
	}
}
