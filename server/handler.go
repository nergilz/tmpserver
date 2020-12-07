package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nergilz/tmpserver/store"
)

func (s *Server) configureRoute() {
	s.router.HandleFunc("/hello", s.hendlerHello())
	s.router.HandleFunc("/user/create", s.handlerCreateUser)
	// s.router.HandleFunc("/user/find", s.handlerFindByEmail)
	s.log.Service("configure Route")
}

func (s *Server) hendlerHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "server is run on index")
		s.log.Info("hendler Hello is run")
	}
}

func (s *Server) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	us := store.InitUserStore(s.db)
	s.log.Service("init user store")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Warningf("Bad request body for create User: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Bad body request , %v", http.StatusText(http.StatusBadRequest))))
		return
	}
	var userFromBody store.UserModel

	err = json.Unmarshal(body, &userFromBody)
	if err != nil {
		s.log.Warningf("Not unmarshal json from r.Body : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Not unmarshal json %v", http.StatusText(http.StatusBadRequest))))
		return
	}
	// TODO: validation data

	err = us.Create(&userFromBody)
	if err != nil {
		s.log.Errorf("user not create : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("user not create %v", http.StatusText(http.StatusBadRequest))))
		return
	}

	s.log.Info("create User")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("create User %v***", http.StatusText(http.StatusOK))))
}
