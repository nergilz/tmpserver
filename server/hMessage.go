package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nergilz/tmpserver/store"
)

// created msg in db
func (s *Server) handlerSendMsg(w http.ResponseWriter, r *http.Request) {
	userFromCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warningf("user not auth %v", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not read the boby from request %v", err)
		return
	}
	defer r.Body.Close()

	var msgFromBody store.SendMsgRequestModel
	err = json.Unmarshal(body, &msgFromBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("cannot body unmarshal json %v", err)
		return
	}
	if err = msgFromBody.SendValidate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not validate message %v", err)
		return
	}
	user, err := s.us.FindByLogin(msgFromBody.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not find user by login %v", err)
		return
	}
	msgForCreate := &store.MsgModel{
		ToID:    user.ID,
		FromID:  userFromCtx.ID,
		Content: msgFromBody.Content,
	}
	err = s.ms.CreateMsg(msgForCreate)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		s.log.Warningf("cannot create msg in db %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	s.log.Info("create message")
}

// hendlerGetMsg ..
func (s *Server) hendlerGetInMsg(w http.ResponseWriter, r *http.Request) {
	userFromCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warningf("user not auth %v", err)
		return
	}
	messages, err := s.ms.FindAllIncomingMsg(userFromCtx.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("cannot find incoming msg %v", err)
		return
	}
	resp, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Errorf("cannot marshal message in json : %v", err)
		return
	}
	w.Write(resp)
	s.log.Infof("get all incoming msg for %v", userFromCtx.Login)
}

// delete message in db, get msg_id from URL
func (s *Server) hendlerDeleteMsg(w http.ResponseWriter, r *http.Request) {
	uCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not get user from context: %v", err)
		return
	}
	vars := mux.Vars(r)
	msgIDfromURL, ok := vars["msg_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warning("cannot get msg_id parameter from url")
		return
	}
	msgID, err := strconv.ParseInt(msgIDfromURL, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("parameter not int: %v", err)
		return
	}
	if err = s.ms.DeleteMsg(msgID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Errorf("cannot delete message : %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	s.log.Infof("Delete Message: %v, id: %v", uCtx.Login, msgID)
}
