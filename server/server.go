package server

import (
	"errors"
	"net/http"

	"github.com/nergilz/tmpserver/store"

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
	us       *store.UserStore
	ms       *store.MsgStore
}

// New server
func New(dbconfig *database.Config, log *logger.Logger) *Server {
	return &Server{
		BindAddr: ":8080",
		dbconf:   dbconfig,
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
	s.log.Service("connect DB")

	if err = db.Init(); err != nil {
		return err
	}
	s.db = db
	s.us = store.InitUserStore(db)
	s.ms = store.InitMsgStore(db)
	s.log.Service("Init DB, userStrore, msgStore")

	return http.ListenAndServe(s.BindAddr, s.router)
}

func (s *Server) configureRoute() {
	s.router.HandleFunc("/hello", s.hendlerHello())

	userRouter := s.router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", s.handlerLoginUser)
	userRouter.HandleFunc("/create", s.handlerCreateUser)
	userRouter.HandleFunc("/delete", s.handlerDeleteUser).Queries("id", "{id:[0-9]+}")

	MsgRouter := s.router.PathPrefix("/message").Subrouter()
	MsgRouter.HandleFunc("/create", s.handlerCreateMsg)

	s.router.Use(s.authMiddleware)
	s.log.Service("configure Route with authMiddleware")
}

// GetUserFromContext get UserModel from Context
func GetUserFromContext(r *http.Request, CtxKeyUser CtxKey) (*store.UserModel, error) {
	valueCtx := r.Context().Value(CtxKeyUser)
	if valueCtx == nil {
		return nil, errors.New("Context is empty")
	}
	userCtx := valueCtx.(*store.UserModel)
	return userCtx, nil
}
