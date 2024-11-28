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
		"415d2711-a692-4378-bb59-774b6ed6ae06",
		"552f9413-2efb-4320-9bc8-bb37a20da932",
		"f7a112b1-0f3e-4528-8eeb-a544a1f4ea3c",
		"d59a7929-ffda-4f1e-8c44-fabb0459868a",
		"2ad90cd8-7ef7-4ba0-ae59-c46e78ceca78",
		//"ad29d179-9773-453f-9d6d-1489fb417033",
		//"78ca936e-227a-48c5-8172-27f5411b8272",
		//"5e681a64-fb13-42a6-be6b-c999602e6f0c",
		//"10399eb9-4df0-4985-9184-c56079f27884",
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
