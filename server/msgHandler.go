package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nergilz/tmpserver/store"
)

// created msg in db
func (s *Server) handlerCreateMsg(w http.ResponseWriter, r *http.Request) {
	userFromCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("user not auth %v", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("not read the boby %v", err)
		return
	}
	defer r.Body.Close()

	var msgFromBody store.MsgModel

	err = json.Unmarshal(body, &msgFromBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("cannot body unmarshal json %v", err)
		return
	}
	if err = msgFromBody.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("cannot validate message %v", err)
		return
	}
	msgForCreate := &store.MsgModel{
		MsgText:  msgFromBody.MsgText,
		SenderID: userFromCtx.ID,
		ChatID:   msgFromBody.ChatID,
	}
	err = s.ms.CreateMsg(msgForCreate)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("cannot create msg %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	s.log.Info("create message")
}

// delete message in db, get msg_id from URL
func (s *Server) hendlerDeleteMsg(w http.ResponseWriter, r *http.Request) {
	uCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("not get user from context: %v", err)
		return
	}
	vars := mux.Vars(r)
	msgIDfromURL, ok := vars["msg_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warning("cannot get msg_id parameter from url")
		return
	}
	msgID, err := strconv.ParseInt(msgIDfromURL, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("parameter not int: %v", err)
		return
	}
	if err = s.ms.DeleteMsg(msgID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Errorf("cannot delete message : %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete Message"))
	s.log.Infof("Delete Message: %v, id: %v", uCtx.Login, msgID)
}
