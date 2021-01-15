package server

import (
	"encoding/json"

	"io/ioutil"
	"net/http"

	"github.com/nergilz/tmpserver/store"
)

func (s *Server) hCreateChat(w http.ResponseWriter, r *http.Request) {
	uCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warningf("cannot get user from context: %v", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("cannot read the body %v", err)
		return
	}
	defer r.Body.Close()

	var reqChat store.ChatModel
	err = json.Unmarshal(body, &reqChat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("cannot body unmarshal json %v", err)
		return
	}
	reqChat.CreatorID = uCtx.ID
	reqChat.UsersIDs = make([]int64, 0)
	reqChat.UsersIDs = append(reqChat.UsersIDs, uCtx.ID)

	err = s.cs.CreateChat(&reqChat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("cannot create chat %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	s.log.Info("create chat: %v", reqChat.Name)
}

func (s *Server) hGetAllChats(w http.ResponseWriter, r *http.Request) {
	uCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		s.log.Warningf("not get user from context: %v", err)
		return
	}
	chats, err := s.cs.GetAllChats(uCtx.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Warningf("cannot get all chats: %v", err)
		return
	}
	resp, err := json.Marshal(chats)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.log.Errorf("cannot marshal message in json : %v", err)
		return
	}

	w.Write(resp)
	s.log.Infof("get all chats for %v", uCtx.Login)
}
