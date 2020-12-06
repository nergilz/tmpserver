package server

import (
	"io"
	"net/http"

	"github.com/nergilz/tmpserver/model"
)

func (s *Server) configureRoute() {
	s.router.HandleFunc("/hello", s.hendlerHello())
	s.router.HandleFunc("/user/create", s.handlerCreateUser)
	// s.router.HandleFunc("/user/find", s.handlerFimdByEmail)
	s.log.Service("configure route")
}

func (s *Server) hendlerHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "server is run on index")
		s.log.Info("hendler Hello is run")
	}
}

func (s *Server) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	us := model.InitUserStore(s.db)
	user := &model.UserModel{
		Email:    "testemail1@mail.ru",
		Password: "testpass1",
		Role:     "user",
	}
	s.log.Service("init user store")
	err := us.Create(user)
	if err != nil {
		s.log.Errorf("user not create : %v", err)
		w.Write([]byte("user not create"))
	} else {
		s.log.Service("create user")
	}
	w.WriteHeader(http.StatusOK)
	s.log.Info("test user create")
}

// func (s *Server) handlerFimdByEmail(w http.ResponseWriter, r *http.Request) {

// }
