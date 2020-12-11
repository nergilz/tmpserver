package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/nergilz/tmpserver/store"
	"github.com/nergilz/tmpserver/utils"
)

func (s *Server) configureRoute() {
	// s.router.HandleFunc("/login", s.)
	s.router.HandleFunc("/user/create", s.handlerCreateUser)
	s.router.HandleFunc("/user/delete", s.handlerDeleteUser).Queries("id", "{id:[0-9]+}")
	s.log.Service("configure Route")
}

func (s *Server) configureHello() {
	s.router.HandleFunc("/hello", s.hendlerHello())
	s.router.Use(s.authMiddleware)
	s.log.Service("configure hello router")
}

func (s *Server) hendlerHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "server is run on index, check JWT")
		s.log.Info("hendler Hello is run, check JWT	")
	}
}

func (s *Server) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Warningf("Bad request body for create User: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Bad body request , %v", http.StatusText(http.StatusBadRequest))))
		return
	}
	defer r.Body.Close()

	var userFromBody store.UserModel

	err = json.Unmarshal(body, &userFromBody)
	if err != nil {
		s.log.Warningf("Not unmarshal json from r.Body : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Not unmarshal json : %v", http.StatusText(http.StatusBadRequest))))
		return
	}

	if err := userFromBody.Validate(); err != nil {
		s.log.Warningf("data is not valid : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("data is not valid : %v", http.StatusText(http.StatusBadRequest))))
		return
	}

	hashPassword, err := store.HashPassword(userFromBody.Password)
	if err != nil {
		s.log.Warningf("No hash password : %v", err)
	}
	userFromBody.Password = hashPassword

	err = s.us.Create(&userFromBody)
	if err != nil {
		s.log.Errorf("User not create : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("user not create %v", http.StatusText(http.StatusBadRequest))))
		return
	}

	JWTtoken, err := utils.CreateJWTtoken(&userFromBody, s.us.GetSecret())
	if err != nil {
		s.log.Errorf("JWTtoken not create : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("JWTtoken not create %v", http.StatusText(http.StatusBadRequest))))
		return
	}
	JWTresp, err := json.Marshal(JWTtoken)
	if err != nil {
		s.log.Errorf("JWTtoken not marshal josn : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("JWTtoken not marshal josn %v", http.StatusText(http.StatusBadRequest))))
		return
	}

	s.log.Info("create User, create JWT token")
	w.WriteHeader(http.StatusCreated)
	w.Write(JWTresp)
}

func (s *Server) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		s.log.Warning("cannot get the 'id' parameter from url")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("cannot get the 'id' parameter from url %v", http.StatusText(http.StatusBadRequest))))
		return
	}
	ID, err := strconv.Atoi(userID)
	if err != nil {
		s.log.Warningf("not int parametr : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("cannot delete User %v", http.StatusText(http.StatusBadRequest))))
		return
	}
	if err = s.us.Delete(ID); err != nil {
		s.log.Errorf("cannot delete User : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("cannot delete User %v", http.StatusText(http.StatusBadRequest))))
		return
	}
	s.log.Infof("delete User id : %v", userID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("delete User, %v", http.StatusText(http.StatusOK))))
}
