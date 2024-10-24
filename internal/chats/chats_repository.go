package chats

import (
	"github.com/bertoxic/graphqlChat/internal/database"
)

type ChatService struct {
	DB  *database.Database
	Rdb *database.Database
}

func newChatService(DB *database.Database, rdb *database.Database) *ChatService {
	return &ChatService{DB: DB, Rdb: rdb}
}

type chatRepo interface {
}
