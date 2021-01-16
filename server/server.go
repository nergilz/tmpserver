package server

import (
	"errors"
	"net/http"

	"github.com/nergilz/tmpserver/utils"

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
	userST   *store.UserStore
	msgST    *store.MsgStore
	chatST   *store.ChatStore
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

	passwordSuperUser, err := utils.GetHashPassword("q1w2e3r4t5y6")
	if err != nil {
		s.log.Warningf("Cannot get hash superuser pass: %v", err)
	}
	if err = db.Init(passwordSuperUser); err != nil {
		return err
	}
	s.db = db
	s.userST = store.InitUserStore(db, s.log)
	s.msgST = store.InitMsgStore(db, s.log)
	s.chatST = store.InitChartStore(db, s.log)

	return http.ListenAndServe(s.BindAddr, s.router)
}

func (s *Server) configureRoute() {
	s.router.HandleFunc("/api/login", s.LoginUser)

	userRouter := s.router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/create", s.CreateUser)
	userRouter.HandleFunc("/delete", s.DeleteUser).Queries("user_id", "{user_id:[0-9]+}")

	chatRouter := s.router.PathPrefix("/chat").Subrouter()
	chatRouter.HandleFunc("/create", s.CreateChat)
	chatRouter.HandleFunc("/getall", s.GetAllChats)

	msgRouter := s.router.PathPrefix("/message").Subrouter()
	msgRouter.HandleFunc("/send", s.SendMsgToChat).Queries("chat_id", "{chat_id:[0-9]+}")
	msgRouter.HandleFunc("/getall", s.GetAllMsgFromChat).Queries("chat_id", "{chat_id:[0-9]+}")
	msgRouter.HandleFunc("/delete", s.DeleteMsg).Queries("msg_id", "{msg_id:[0-9]+}")

	userRouter.Use(s.authMiddleware)
	msgRouter.Use(s.authMiddleware)
	chatRouter.Use(s.authMiddleware)
	s.log.Service("configure Route with authMiddleware")
}

// GetUserFromContext get UserModel from Context
func GetUserFromContext(r *http.Request, CtxKeyUser CtxKey) (*store.UserModel, error) {
	valueCtx := r.Context().Value(CtxKeyUser)
	if valueCtx == nil {
		return nil, errors.New("cannot get user from ctx, context is empty")
	}
	userCtx := valueCtx.(*store.UserModel)
	return userCtx, nil
}
