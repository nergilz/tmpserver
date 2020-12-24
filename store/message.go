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
	ID      int64  `json:"id"`
	ToID    int64  `json:"to_id"`   // кому
	FromID  int64  `json:"from_id"` // от кого
	Content string `json:"content"`
}

// SendMsgRequestModel send msg from user to user
type SendMsgRequestModel struct {
	Login   string `json:"login"`
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
func (ms *MsgStore) CreateMsg(msg *MsgModel) error {
	var id int64
	q := `INSERT INTO messages (to_id, from_id, content) VALUES ($1,$2,$3) RETURNING id`
	err := ms.db.Conn().QueryRow(q, msg.ToID, msg.FromID, msg.Content).Scan(&id)
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
		&msg.ToID,
		&msg.FromID,
		&msg.Content,
	); err != nil {
		return nil, err
	}
	return msg, nil
}

// FindMsgByLogin ..
func (ms *MsgStore) FindMsgByLogin(login string) (*MsgModel, error) {
	msg := &MsgModel{}
	q := `SELECT owner_id, chat_id, text FROM messages WHERE id=$1 VALUES ($1)`
	if err := ms.db.Conn().QueryRow(q, login).Scan(
		&msg.ToID,
		&msg.FromID,
		&msg.Content,
	); err != nil {
		return nil, err
	}
	return msg, nil
}

// FindAllIncomingMsg все принятые пользователем
func (ms *MsgStore) FindAllIncomingMsg(userFromCtxID int64) ([]*MsgModel, error) {
	messages := []*MsgModel{}
	q := `SELECT * FROM messages where from_id=$1`
	rows, err := ms.db.Conn().Query(q, userFromCtxID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := new(MsgModel)
		if err := rows.Scan(&item.ID, &item.ToID, &item.FromID, &item.Content); err != nil {
			return nil, err
		}
		messages = append(messages, item)
	}
	rows.Close()

	return messages, nil
}

// func (ms *MsgStore) FindAllOutgoingMsg(fromID int64) ([]*MsgModel, error) {}

// Validate ..
func (msg *MsgModel) Validate() error {
	if msg.Content == "" {
		return errors.New("text cannot be empty")
	}
	return nil
}

// SendValidate ..
func (msg *SendMsgRequestModel) SendValidate() error {
	if msg.Login == "" {
		return errors.New("login cannot be empty")
	}
	if msg.Content == "" {
		return errors.New("text cannot be empty")
	}
	return nil
}
