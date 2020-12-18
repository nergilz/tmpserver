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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("not read the boby %v", err)
		return
	}
	defer r.Body.Close()

	var msgFromBody store.MsgModel
	userFromCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("user not auth %v", err)
		return
	}
	err = json.Unmarshal(body, &msgFromBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("cannot body unmarshal json %v", err)
		return
	}
	// TODO
	// get validate body data
	// get id user_to (from r.body)
	msgForCreate := &store.MsgModel{
		Description: msgFromBody.Description,
		MsgText:     msgFromBody.MsgText,
		OwnerID:     userFromCtx.ID,
		//UserToID:
	}
	err = s.ms.CreateMsg(msgForCreate)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("cannot create msg %v", err)
		return
	}
	s.log.Info("create message")
	w.WriteHeader(http.StatusCreated)
}
