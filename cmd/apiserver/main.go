package main

import (
	"github.com/0LuigiCode0/Library/logger"
	"github.com/nergilz/tmpserver/internal/app/apiserver"
)

func main() {
	log := logger.InitLogger("")
	log.Infof("main start")
	config := apiserver.NewConfig()
	server := apiserver.NewServer(config)
	err := server.Start()
	if err != nil {
		log.Fatalf("err start server %s", err)
	}
}
