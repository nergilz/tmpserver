package store

import (
	"errors"

	"github.com/nergilz/tmpserver/database"
)

// MsgStore message store
type MsgStore struct {
	db *database.DB
}

// MsgModel message model
type MsgModel struct {
	ID      int64  `json:"id"`
	OwnerID int64  `json:"owner"`   // who request
	UserTo  string `json:"user_to"` // who responce
	Title   string `json:"title"`
	MsgText string `json:"text"`
}

// InitMsgStore initialization message store
func InitMsgStore(db *database.DB) *MsgStore {
	ms := new(MsgStore)
	ms.db = db
	return ms
}

// CreateMsg create message in database
func (ms *MsgStore) CreateMsg(msg *MsgModel) error {
	var id int64
	q := `INSERT INTO messages (owner_id, user_to, title, text) VALUES ($1,$2,$3,$4) RETURNING id`
	err := ms.db.Conn().QueryRow(q, msg.OwnerID, msg.UserTo, msg.Title, msg.MsgText).Scan(&id)
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
	q := `SELECT owner_id, user_to, title, text FROM messages WHERE id=$1 VALUES ($1)`
	if err := ms.db.Conn().QueryRow(q, msgID).Scan(
		&msg.OwnerID,
		&msg.UserTo,
		&msg.Title,
		&msg.MsgText,
	); err != nil {
		return nil, err
	}
	return msg, nil
}

// Validate ..
func (msg *MsgModel) Validate() error {
	if msg.Title == "" {
		return errors.New("title cannot be empty")
	}
	if msg.UserTo == "" {
		return errors.New("login cannot be empty")
	}
	if msg.MsgText == "" {
		return errors.New("text cannot be empty")
	}
	return nil
}

// TODO
// return request & response msg
