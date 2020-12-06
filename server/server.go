package server

import (
	"net/http"

	"github.com/0LuigiCode0/Library/logger"
	"github.com/gorilla/mux"
	"github.com/nergilz/tmpserver/database"
)

// Server struct
type Server struct {
	BindAddr string
	dbconf   *database.Config
	log      *logger.Logger
	router   *mux.Router
	db       *database.DB
}

// New server
func New(config *database.Config, log *logger.Logger) *Server {
	return &Server{
		BindAddr: ":8080",
		dbconf:   config,
		log:      log,
		router:   mux.NewRouter(),
	}
}

// Start ..
func (s *Server) Start() error {
	s.log.Service("Server Start")
	s.configureRoute()

	db, err := database.Connect(s.dbconf)
	if err != nil {
		s.log.Errorf("error connect DB : %v", err)
		return err
	}
	if err = db.Init(); err != nil {
		s.log.Errorf("not configure DB: %v", err)
	}
	s.db = db

	return http.ListenAndServe(s.BindAddr, s.router)
}

// func (s *Server) configureDB() error {
// 	return nil
// }
