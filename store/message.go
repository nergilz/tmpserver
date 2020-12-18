package store

import (
	"github.com/nergilz/tmpserver/database"
)

// MsgStore message store
type MsgStore struct {
	db *database.DB
}

// MsgModel message model
type MsgModel struct {
	ID          int64  `json:"id"`
	OwnerID     int64  `json:"owner"` // who request
	UserToID    int64  `json:"user"`  // who responce
	Description string `json:"desc"`
	MsgText     string `json:"text"`
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
	q := `INSERT INTO messages (owner_id, user_id, description, text) VALUES ($1,$2,$3,$4) RETURNING id`
	err := ms.db.Conn().QueryRow(q, msg.OwnerID, msg.UserToID, msg.Description, msg.MsgText).Scan(&id)
	if err != nil {
		return err
	}
	msg.ID = id
	return nil
}

// DeleteMsg delete message in database
func (ms *MsgStore) DeleteMsg(msgID int64) error {
	q := `DELETE FROM messages WHERE id=$1`
	if err := ms.db.Conn().QueryRow(q, msgID).Err(); err != nil {
		return err
	}
	return nil
}
