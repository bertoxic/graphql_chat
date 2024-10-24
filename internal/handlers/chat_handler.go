package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/bertoxic/graphqlChat/internal/chats"
	"github.com/bertoxic/graphqlChat/internal/database"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatRepository struct {
	app *config.AppConfig
	DB  database.DatabaseRepo
	RDB database.RedisClient
	hub *chats.Hub
}

func NewChatRepository(app *config.AppConfig, DB database.DatabaseRepo, RDB database.RedisClient, hub *chats.Hub) *ChatRepository {
	return &ChatRepository{app: app, DB: DB, RDB: RDB, hub: hub}
}

var ChatRepo *ChatRepository

func NewChatRepoInit(repo *ChatRepository) {
	ChatRepo = repo
}

func (ch *ChatRepository) HandleChatWs(w http.ResponseWriter, r *http.Request) {
	serveWs(ch.hub, w, r)
}

func serveWs(hub *chats.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	var ListID = []string{
		"7c15d670-3e50-4246-ab60-9a2e834e9da6",
		"f205ca45-4906-4c77-b936-17b4452f3d1f",
		"dedff378-43f9-4b6a-85b7-e7310d5b8a83",
		"8ca881fc-0fc3-4892-a3b6-0d6c754c58c9",
		"af3024a9-a492-4313-aea1-240695a2fba0",
		"ef259b7a-8085-4c1c-bed2-3cf3234ee45a",
		"3712411e-2855-4c9b-b01e-98f0bd8ca83e",
		"f16e7add-614a-42f6-ba0b-e9e9cb751a28",
		"18ca2dda-d86d-4824-9f9c-6fcb1b2a7636",
	}

	rand.Shuffle(len(ListID), func(i, j int) {
		ListID[i], ListID[j] = ListID[j], ListID[i]
	})

	randomIndex := rand.Intn(len(ListID))
	randomID := ListID[randomIndex]

	client := &chats.Client{
		Config: config.AppConfig{},
		ID:     randomID,
		User:   models.User{ID: randomID},
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
	fmt.Printf("userID is :", client.User.ID)

	//client := &chats.Client{
	//	Config: config.AppConfig{},
	//	ID:     "78ca936e-227a-48c5-8172-27f5411b8272",
	//	User:   models.User{},
	//	Hub:    hub,
	//	Conn:   conn,
	//	Send:   make(chan []byte, 256),
	//}

	client.Hub.Register <- client
	go client.HandleConnection()

}
