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
	ID       int64  `json:"id"`
	SenderID int64  `json:"sender_id"` // who send message
	ChatID   int64  `json:"chat_id"`   // to whon send message
	MsgText  string `json:"text"`
}

// InitMsgStore initialization message store
func InitMsgStore(db *database.DB, log *logger.Logger) *MsgStore {
	ms := new(MsgStore)
	ms.db = db
	log.Service("init message store")
	return ms
}

// CreateMsg create message in database
func (ms *MsgStore) CreateMsg(msg *MsgModel) error {
	var id int64
	q := `INSERT INTO messages (sender_id, chat_id, text) VALUES ($1,$2,$3) RETURNING id`
	err := ms.db.Conn().QueryRow(q, msg.SenderID, msg.ChatID, msg.MsgText).Scan(&id)
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
func (ms *MsgStore) FindMsgByID(msgID int64) (*MsgModel, error) {
	msg := &MsgModel{}
	q := `SELECT owner_id, chat_id, text FROM messages WHERE id=$1 VALUES ($1)`
	if err := ms.db.Conn().QueryRow(q, msgID).Scan(
		&msg.SenderID,
		&msg.ChatID,
		&msg.MsgText,
	); err != nil {
		return nil, err
	}
	return msg, nil
}

// Validate ..
func (msg *MsgModel) Validate() error {
	if msg.MsgText == "" {
		return errors.New("text cannot be empty")
	}
	return nil
}
