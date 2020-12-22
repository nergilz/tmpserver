package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nergilz/tmpserver/store"
)

/*
TODO тут все переделать
	принять сообщение msgModel от userFromContext
	в chat записать users ids
	привязать chat к message в db
*/

func (s *Server) sendMessage(w http.ResponseWriter, r *http.Request) {
	uCtx, err := GetUserFromContext(r, СtxKeyUser)
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
		s.log.Warningf("not read the body %v", err)
		return
	}
	defer r.Body.Close()

	var chatFromBody store.ChatModel

	if err = json.Unmarshal(body, &chatFromBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("cannot body unmarshal json %v", err)
		return
	}

	// TODO validateChat()

	userids := make([]int64, 0)
	userids = append(userids, uCtx.ID)
	chatForCreate := &store.ChatModel{
		MsgIDs:  chatFromBody.MsgIDs,
		UserIDs: userids,
		Private: chatFromBody.Private,
	}
	if err = s.cs.CreateChat(chatForCreate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		s.log.Warningf("cannot create chat %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	s.log.Info("create chat")
}

func (s *Server) updateChat(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) deleteChat(w http.ResponseWriter, r *http.Request) {

}
