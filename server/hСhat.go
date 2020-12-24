package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nergilz/tmpserver/store"
)

func (s *Server) handlerSendMessage(w http.ResponseWriter, r *http.Request) {
	uCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not get user from context: %v", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not read the body %v", err)
		return
	}
	defer r.Body.Close()

	var chatFromBody store.ChatModel
	if err = json.Unmarshal(body, &chatFromBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("cannot body unmarshal json %v", err)
		return
	}

	// TODO validateChat()
	chatFromBody.UserIDs = append(chatFromBody.UserIDs, uCtx.ID)
	chatForCreate := &store.ChatModel{
		UserIDs: chatFromBody.UserIDs,
		Private: chatFromBody.Private,
	}
	if err = s.cs.CreateChat(chatForCreate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("cannot create chat %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	s.log.Info("create chat")
}

// func (s *Server) handlerCheckChat(w http.ResponseWriter, r *http.Request) {
// 	uCtx, err := GetUserFromContext(r, СtxKeyUser)
// 	if err != nil {
// 		w.WriteHeader(http.StatusForbidden)
// 		w.Write([]byte(err.Error()))
// 		s.log.Warningf("not get user from context: %v", err)
// 		return
// 	}
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 		s.log.Warningf("not read the body %v", err)
// 		return
// 	}
// 	defer r.Body.Close()

// 	var chatFromBody store.ChatModel
// 	if err = json.Unmarshal(body, &chatFromBody); err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 		s.log.Warningf("cannot body unmarshal json %v", err)
// 		return
// 	}
// 	_, err = s.cs.CheckChat()
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 		s.log.Errorf("cannot check chat %v", err)
// 		return
// 	}
// 	w.WriteHeader(http.StatusCreated)
// 	s.log.Info("create chat")

// }

func (s *Server) handlerGetListChats(w http.ResponseWriter, r *http.Request) {

}
