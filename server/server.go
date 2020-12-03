package server

import (
	"io"
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
}

// New server
func New(config *database.Config) *Server {
	return &Server{
		BindAddr: ":8080",
		dbconf:   config,
		log:      logger.InitLogger(""),
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

	err = db.Init()
	if err != nil {
		s.log.Errorf("not configure DB: %v", err)
	}

	err = http.ListenAndServe(s.BindAddr, s.router)
	if err != nil {
		s.log.Errorf("Server is abandon : %v", err)
		return err
	}
	s.log.Service("Listen And Serve")

	return nil
}

func (s *Server) configureRoute() {
	s.router.HandleFunc("/hello", s.hendlerHello())
	s.log.Service("configure route")
}

func (s *Server) hendlerHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "server is run on index")
		s.log.Info("hendler Hello is run")
	}
}

// func (s *Server) configureDB() error {
// 	return nil
// }
