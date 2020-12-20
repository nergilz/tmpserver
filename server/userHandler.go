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

func (s *Server) hendlerHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "server is run on index, check JWT")
		s.log.Info("hendler Hello is run, check JWT	")
	}
}

// handlerCreateUser run only with token, role: super_user
func (s *Server) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	userFromCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("not user in context: %v", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("not read the boby: %v", err)
		return
	}
	defer r.Body.Close()

	if userFromCtx.Role == "super_user" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("user from context not super_user: %v", err)
		return
	}
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
	hashPassword, err := utils.GetHashPassword(userFromBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("not valid password: %v", err)
		return
	}

	userFromBody.Password = hashPassword

	err = s.us.Create(&userFromBody)
	if err != nil {
		s.log.Errorf("User not create : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("user not create %v", http.StatusText(http.StatusBadRequest))))
		return
	}
	// create token for new user
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
	u, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(http.StatusText(http.StatusForbidden)))
		s.log.Warningf("not get user from context: %v", err)
		return
	}
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		s.log.Warning("cannot get the 'id' parameter from url")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("cannot get the 'id' parameter from url %v", http.StatusText(http.StatusBadRequest))))
		return
	}
	ID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		s.log.Warningf("not int parametr : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("cannot delete User %v", http.StatusText(http.StatusBadRequest))))
		return
	}
	if u.Role != "super_user" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(http.StatusText(http.StatusForbidden)))
		s.log.Warning("cannot delete user, yuo not super_user")
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

// handlerLoginUser login all users & create token
func (s *Server) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Warningf("Bad request body for login User: %v", err)
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

	userFromDB, err := s.us.FindByLogin(userFromBody.Login)
	if err != nil {
		s.log.Warningf("Not find user by login : %v", err)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("Not unmarshal json : %v", http.StatusText(http.StatusForbidden))))
		return
	}
	hashPassUserFromBody, err := utils.GetHashPassword(userFromBody.Password)
	if err != nil {
		s.log.Warningf("No get hash password : %v", err)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("No get hash password : %v", http.StatusText(http.StatusForbidden))))
		return
	}
	if hashPassUserFromBody != userFromDB.Password {
		s.log.Warningf("No get hash password : %v", err)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("No get hash password : %v", http.StatusText(http.StatusForbidden))))
		return
	}
	/*
		тут создается токен юзера из базы
	*/
	JWTtoken, err := utils.CreateJWTtoken(userFromDB, s.us.GetSecret()) // !!!
	if err != nil {
		s.log.Errorf("JWTtoken for login not create : %v", err)
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
	s.log.Info("Create JWT token for Login User")
	w.WriteHeader(http.StatusCreated)
	w.Write(JWTresp)
}
