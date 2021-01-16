package store

import (
	"errors"

	"github.com/0LuigiCode0/Library/logger"
	"github.com/nergilz/tmpserver/database"
)

// MsgStore message store
type MsgStore struct {
	db *database.DB
}

// MsgModel message model
type MsgModel struct {
	ID         int64  `json:"id"`
	ToChatID   int64  `json:"to_chat_id"`   // кому
	FromUserID int64  `json:"from_user_id"` // от кого
	Content    string `json:"content"`
}

// MsgRequestModel send msg from user to chat
type MsgRequestModel struct {
	ID      int64
	ChatID  int64  `json:"chat_id"`
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
}

// InitMsgStore initialization message store
func InitMsgStore(db *database.DB, log *logger.Logger) *MsgStore {
	ms := new(MsgStore)
	ms.db = db
	log.Service("init message store")
	return ms
}

// CreateMsg create message in database
func (ms *MsgStore) CreateMsg(msg *MsgRequestModel) error {
	var id int64
	q := `INSERT INTO messages (chat_id, user_id, content) VALUES ($1,$2,$3) RETURNING id`
	err := ms.db.Conn().QueryRow(q, msg.ChatID, msg.UserID, msg.Content).Scan(&id)
	if err != nil {
		return err
	}
	msg.ID = id
	return nil
}

// DeleteMsg delete message in database
func (ms *MsgStore) DeleteMsg(msgID int64) error {
	q := `DELETE FROM messages WHERE id = $1`
	if err := ms.db.Conn().QueryRow(q, msgID).Err(); err != nil {
		return err
	}
	return nil
}

// FindMsgByID return msg by id
func (ms *MsgStore) FindMsgByID(msgID int64) (*MsgRequestModel, error) {
	msg := &MsgRequestModel{}
	q := `SELECT id, chat_id, user_id, content FROM messages WHERE id=$1 VALUES($1)`
	if err := ms.db.Conn().QueryRow(q, msgID).Scan(
		&msg.ChatID,
		&msg.UserID,
		&msg.Content,
	); err != nil {
		return nil, err
	}
	return msg, nil
}

// FindAllMsgFromChat ...
func (ms *MsgStore) FindAllMsgFromChat(chatID, userID int64) ([]*MsgRequestModel, error) {
	messages := []*MsgRequestModel{}
	q := `SELECT * FROM messages where chat_id=$1 AND user_id=$2`
	rows, err := ms.db.Conn().Query(q, chatID, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := new(MsgRequestModel)
		if err := rows.Scan(&item.ID, &item.ChatID, &item.UserID, &item.Content); err != nil {
			return nil, err
		}
		messages = append(messages, item)
	}
	rows.Close()

	return messages, nil
}

// Validate ..
func (msg *MsgRequestModel) Validate() error {
	if msg.Content == "" {
		return errors.New("text cannot be empty")
	}
	return nil
}

// SendValidate ..
func (msg *MsgRequestModel) SendValidate() error {
	if msg.Content == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}

//-----------------------------------------------------------------------------------------------------

// FindAllIncomingMsg все принятые пользователем
func (ms *MsgStore) FindAllIncomingMsg(userID int64) ([]*MsgModel, error) {
	messages := []*MsgModel{}
	q := `SELECT * FROM messages where to_id=$1`
	rows, err := ms.db.Conn().Query(q, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := new(MsgModel)
		if err := rows.Scan(&item.ID, &item.ToChatID, &item.FromUserID, &item.Content); err != nil {
			return nil, err
		}
		messages = append(messages, item)
	}
	rows.Close()

	return messages, nil
}

// FindAllOutgoingMsg все отправленные пользователем
func (ms *MsgStore) FindAllOutgoingMsg(userID int64) ([]*MsgModel, error) {
	messages := []*MsgModel{}
	q := `SELECT * FROM messages where from_id=$1`
	rows, err := ms.db.Conn().Query(q, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := new(MsgModel)
		if err := rows.Scan(&item.ID, &item.ToChatID, &item.FromUserID, &item.Content); err != nil {
			return nil, err
		}
		messages = append(messages, item)
	}
	rows.Close()

	return messages, nil
}
