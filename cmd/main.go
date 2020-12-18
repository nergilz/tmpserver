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

/*
Scenario
1.sending mail withing the service
2.sending and reciving mail from SMTP server

TODO
переписать создание user по инвайтной ссылке (сейчас не правильно, просто логиниться)

Questions
как найти user in /delete если нет id,
всегда ли мы создаем claims со всеми полями

*/
