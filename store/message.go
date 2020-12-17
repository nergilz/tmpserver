package store

import (
	"github.com/nergilz/tmpserver/database"
)

// MsgStore message store
type MsgStore struct {
	db *database.DB
}

// MsgModel model request
type MsgModel struct {
	ID          int64  `json:"id"`
	OwnerID     int64  `json:"owner"` // кто отправил
	UserToID    int64  `json:"user"`  // кому отправил
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
func (ms *MsgStore) CreateMsg(m *MsgModel) error {
	var id int64
	q := `INSERT INTO messages (ownerID, userID, desc, text) VALUES ($1,$2,$3,$4) RETURNING id`
	err := ms.db.Conn().QueryRow(q, m.OwnerID, m.UserToID, m.Description, m.MsgText).Scan(&id)
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}
