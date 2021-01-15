package store

import (
	"fmt"

	"github.com/0LuigiCode0/Library/logger"
	"github.com/lib/pq"
	"github.com/nergilz/tmpserver/database"
)

// ChatStore chat store
type ChatStore struct {
	db *database.DB
}

// // RequestChat ...
// type RequestChat struct {
// 	Name   string `json:"name"`
// }

// ChatModel chat model
type ChatModel struct {
	ID         int64   `json:"chat_id"`
	Name       string  `json:"chat_name"`
	CreatorID  int64   `json:"creator_id"`
	UsersIDs   []int64 `json:"users_ids"`
	Individual bool    `json:"individual"`
}

// InitChartStore ...
func InitChartStore(db *database.DB, log *logger.Logger) *ChatStore {
	cs := new(ChatStore)
	cs.db = db
	log.Service("init chat store")
	return cs
}

// CreateChat ...
func (cs *ChatStore) CreateChat(chat *ChatModel) error {
	var cid int64
	qc := `INSERT INTO chats (name, creator_id, users_ids, individual) VALUES($1, $2, $3, $4) RETURNING id`
	if err := cs.db.Conn().QueryRow(
		qc,
		chat.Name,
		chat.CreatorID,
		pq.Array(chat.UsersIDs),
		chat.Individual).Scan(&cid); err != nil {
		return err
	}
	chat.ID = cid

	qp := `INSERT INTO participants (chat_id, users_ids) VALUES($1, $2)`
	if err := cs.db.Conn().QueryRow(qp, cid, pq.Array(chat.UsersIDs)).Err(); err != nil {
		return err
	}
	return nil
}

// DeleteChat ...
func (cs *ChatStore) DeleteChat(chatID int64) error {
	q := `DELETE FROM chats WHERE id = $1`
	if err := cs.db.Conn().QueryRow(q, chatID).Err(); err != nil {
		return err
	}
	return nil
}

// FindChatByID ...
func (cs *ChatStore) FindChatByID(chatID int64) (*ChatModel, error) {
	chat := &ChatModel{}
	q := `SELECT id, name, creator_id, users_ids, individual WHERE id=$1 VALUES($1)`
	if err := cs.db.Conn().QueryRow(q, chatID).Scan(
		&chat.ID,
		&chat.Name,
		&chat.CreatorID,
		&chat.UsersIDs,
		&chat.Individual,
	); err != nil {
		return nil, err
	}
	return chat, nil
}

// UpdateChat ..
func (cs *ChatStore) UpdateChat(chat *ChatModel) (*ChatModel, error) {
	newChat := new(ChatModel)
	q := `UPDATE chats SET name=$1, users_ids=$2, individual=$3	
			WHERE id=$4	
			RETURNING id, name, creator_id, users_ids, private`
	if err := cs.db.Conn().QueryRow(
		q,
		chat.Name,
		chat.UsersIDs,
		chat.Individual,
		chat.ID).Err(); err != nil {
		return nil, err
	}
	return newChat, nil
}

// GetAllChats through messages table
func (cs *ChatStore) GetAllChats(userID int64) ([]*ChatModel, error) {
	chats := []*ChatModel{}
	q := `SELECT * FROM chats WHERE $1 <@ users_ids`
	exist := make([]int64, 0)
	exist = append(exist, userID)
	rows, err := cs.db.Conn().Query(q, pq.Array(exist))
	if err != nil {
		fmt.Println("[error query] :", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := new(ChatModel)
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.CreatorID,
			pq.Array(&item.UsersIDs),
			&item.Individual)
		if err != nil {
			fmt.Println("[error scan] :", err)
			return nil, err
		}
		chats = append(chats, item)
	}
	return chats, nil
}
