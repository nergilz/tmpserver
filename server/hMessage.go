package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nergilz/tmpserver/store"
)

func (s *Server) hSendMsgInChat(w http.ResponseWriter, r *http.Request) {
	uCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		s.log.Warningf("cannot get user from context: %v", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.log.Warningf("cannot read the body %v", err)
		return
	}
	defer r.Body.Close()

	var msgRequest store.MsgRequestModel
	if err = json.Unmarshal(body, &msgRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.log.Warningf("cannot body unmarshal json %v", err)
		return
	}
	msgRequest.UserID = uCtx.ID

	if err := s.ms.CreateMsg(&msgRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.log.Warningf("cannot create msg, %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	s.log.Infof("send msg in chat: %v", msgRequest.ChatID)
}

func (s *Server) hDeleteMsg(w http.ResponseWriter, r *http.Request) {
	uCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		s.log.Warningf("not get user from context: %v", err)
		return
	}
	vars := mux.Vars(r)
	msgIDfromURL, ok := vars["msg_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.log.Warning("cannot get msg_id parameter from url")
		return
	}
	msgID, err := strconv.ParseInt(msgIDfromURL, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.log.Warningf("parameter not int: %v", err)
		return
	}
	// var msg *store.MsgRequestModel
	// msg, err = s.ms.FindMsgByID(msgID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("cannot find msg by id"))
	// 	s.log.Warningf("cannot find msg by id: %v", err)
	// 	return
	// }
	// if msg.ID == msgID {
	if err = s.ms.DeleteMsg(msgID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.log.Errorf("cannot delete message : %v", err)
		return
	}
	// } else {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	s.log.Errorf("not message id: %v", err)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	s.log.Infof("Delete Message: %v, id: %v", uCtx.Login, msgID)
}

func (s *Server) hGetAllMsgFromChat(w http.ResponseWriter, r *http.Request) {
	uCtx, err := GetUserFromContext(r, СtxKeyUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		s.log.Warningf("not get user from context: %v", err)
		return
	}
	vars := mux.Vars(r)
	msgIDfromURL, ok := vars["chat_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.log.Warning("cannot get chat_id parameter from url")
		return
	}
	chatID, err := strconv.ParseInt(msgIDfromURL, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.log.Warningf("invalid parameter from url: %v", err)
		return
	}

	allMsg, err := s.ms.FindAllMsgFromChat(chatID, uCtx.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot find all msg from chat"))
		s.log.Warningf("cannot find all msg from chat: %v", err)
		return
	}
	res, err := json.Marshal(allMsg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.log.Warningf("cannot marshal json from allmsg: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
	s.log.Info("get all msg from chat")
}

// // created msg in db
// func (s *Server) handlerSendMsg(w http.ResponseWriter, r *http.Request) {
// 	userFromCtx, err := GetUserFromContext(r, СtxKeyUser)
// 	if err != nil {
// 		w.WriteHeader(http.StatusForbidden)
// 		w.Write([]byte(err.Error()))
// 		s.log.Warningf("user not auth %v", err)
// 		return
// 	}
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 		s.log.Warningf("not read the boby from request %v", err)
// 		return
// 	}
// 	defer r.Body.Close()

// 	var msgFromBody store.SendMsgRequestModel
// 	err = json.Unmarshal(body, &msgFromBody)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 		s.log.Warningf("cannot body unmarshal json %v", err)
// 		return
// 	}
// 	if err = msgFromBody.SendValidate(); err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 		s.log.Warningf("not validate message %v", err)
// 		return
// 	}
// 	user, err := s.us.FindByLogin(msgFromBody.ID)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 		s.log.Warningf("not find user by login %v", err)
// 		return
// 	}
// 	msgForCreate := &store.MsgModel{
// 		ToChatID:   user.ID,
// 		FromUserID: userFromCtx.ID,
// 		Content:    msgFromBody.Content,
// 	}
// 	err = s.ms.CreateMsg(msgForCreate)
// 	if err != nil {
// 		w.WriteHeader(http.StatusConflict)
// 		w.Write([]byte(err.Error()))
// 		s.log.Warningf("cannot create msg in db %v", err)
// 		return
// 	}
// 	w.WriteHeader(http.StatusCreated)
// 	s.log.Info("create message")
// }
