package apiserver

import (
	"io"
	"net/http"

	"github.com/0LuigiCode0/Library/logger"

	"github.com/nergilz/tmpserver/internal/app/store"

	"github.com/gorilla/mux"
)

// APIServer ...
type APIServer struct {
	log    *logger.Logger
	config *Config
	router *mux.Router
	store  *store.Store
}

// NewServer ...
func NewServer(config *Config) *APIServer {
	return &APIServer{
		log:    logger.InitLogger(""),
		config: config,
		router: mux.NewRouter(),
	}
}

// Start ...
func (s *APIServer) Start() error {
	s.log.Service("start server")
	s.configureRouter()
	if err := s.configureStore(); err != nil {
		s.log.Fatal("not configure router")
		return err
	}
	err := http.ListenAndServe(s.config.BindAddr, s.router)
	if err != nil {
		s.log.Fatalf("not listen and serve %s", err)
		return err
	}
	s.log.Service("listen end serve")

	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
	s.log.Service("configure router")
}

func (s *APIServer) configureStore() error {
	stnew := store.NewStore(s.config.Store)
	if err := stnew.Open(); err != nil {
		s.log.Fatalf("configure err: %s", err)
		return err
	}
	s.store = stnew
	s.log.Service("configure DB store")
	return nil
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "server is run")
		s.log.Info("hendleHello is run")
	}
}
