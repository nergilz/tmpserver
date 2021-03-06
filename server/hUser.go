package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/nergilz/tmpserver/store"
	"github.com/nergilz/tmpserver/utils"
)

// CreateUser only for super_user
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	userFromCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not user in context: %v", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not read the boby: %v", err)
		return
	}
	defer r.Body.Close()

	if userFromCtx.Role != "super_user" {
		w.WriteHeader(http.StatusForbidden)
		s.log.Warningf("user from context not super_user: %v", err)
		return
	}
	var userFromBody store.UserModel

	err = json.Unmarshal(body, &userFromBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("Not unmarshal json from r.Body : %v", err)
		return
	}
	if err := userFromBody.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("data is not valid : %v", err)
		return
	}
	hashPassword, err := utils.GetHashPassword(userFromBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not valid password: %v", err)
		return
	}
	userFromBody.Password = hashPassword

	err = s.userST.Create(&userFromBody)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(err.Error()))
		s.log.Errorf("User not create : %v", err)
		return
	}

	s.log.Infof("Create User: %v", userFromBody.Login)
	w.WriteHeader(http.StatusCreated)
}

// DeleteUser only for super_user
func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	u, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not get user from context: %v", err)
		return
	}
	vars := mux.Vars(r)
	userID, ok := vars["user_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warning("cannot get the 'id' parameter from url")
		return
	}
	ID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not int parametr : %v", err)
		return
	}
	if u.Role != "super_user" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warning("cannot delete user, yuo not super_user")
		return
	}
	if err = s.userST.Delete(ID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Errorf("cannot delete User : %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete User"))
	s.log.Infof("Delete User: %v, id: %v", u.Login, userID)
}

// LoginUser login all users & create token
func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("Bad request body for login User: %v", err)
		return
	}
	defer r.Body.Close()

	var userFromBody store.UserModel

	err = json.Unmarshal(body, &userFromBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("Not unmarshal json from r.Body : %v", err)
		return
	}

	if err := userFromBody.Validate(); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warningf("data is not valid : %v", err)
		return
	}

	userFromDB, err := s.userST.FindByLogin(userFromBody.Login)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		s.log.Warningf("User not found: %v", err)
		return
	}
	hashPassUserFromBody, err := utils.GetHashPassword(userFromBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("Cannot get hash password : %v", err)
		return
	}
	if hashPassUserFromBody != userFromDB.Password {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warningf("Invalid Password : %v", err)
		return
	}
	//	create token:
	JWTtoken, err := utils.CreateJWTtoken(userFromDB, s.userST.GetSecret())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Errorf("JWTtoken not create : %v", err)
		return
	}
	JWTresp, err := json.Marshal(JWTtoken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Errorf("JWTtoken not marshal in json : %v", err)
		return
	}
	w.Write(JWTresp)
	s.log.Infof("Login user: %v & create JWTtoken", userFromDB.Login)
}
