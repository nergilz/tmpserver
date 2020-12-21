package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
		Title:   msgFromBody.Title,
		MsgText: msgFromBody.MsgText,
		OwnerID: userFromCtx.ID,
		UserTo:  msgFromBody.UserTo,
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
