package models

import (
	"encoding/base64"
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	FromID    string    `json:"from_id"`
	ToID      string    `json:"to_id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`   // "text", "image", etc.
	Status    string    `json:"status"` // "sent", "delivered", "read"
	Timestamp time.Time `json:"timestamp"`
}

func NewMessage(ID string, fromID string, toID string, content string, Type string, status string) *Message {
	return &Message{ID: ID, FromID: fromID, ToID: toID, Content: content, Type: Type, Status: status, Timestamp: time.Now()}
}
func (m *Message) EncodeContent() error {
	switch m.Type {
	case "image", "audio":
		// Assuming content is a base64-encoded string of the file
		_, err := base64.StdEncoding.DecodeString(m.Content)
		if err != nil {
			return err
		}
		// If content is not base64-encoded, encode it
		// encodedContent := base64.StdEncoding.EncodeToString([]byte(m.Content))
		// m.Content = encodedContent
		// return nil
	default:
		// No encoding needed for text
		return nil
	}
	return nil
}

func (m *Message) DecodeContent() ([]byte, error) {
	switch m.Type {
	case "image", "audio":
		// Decode the base64-encoded content
		decodedBytes, err := base64.StdEncoding.DecodeString(m.Content)
		if err != nil {
			return nil, err
		}
		return decodedBytes, nil
	default:
		// No decoding needed for text
		return nil, nil
	}
}

//type Client struct {
//	ID   string
//	User User
//	Hub  *Hub
//	Conn *websocket.Conn
//	send chan []byte
//}

//type Hub struct {
//	Clients    map[string]*Client
//	Register   chan *Client
//	UnRegister chan *Client
//	Redis      *database.RedisClient
//	mu         sync.RWMutex
//	ctx        context.Context
//	cancel     context.CancelFunc
//	Send       chan []byte
//	// For direct messages
//	Private chan *Message
//	// For public posts
//	Public chan *Message
//	// Database for persistent storage
//	DB database.DatabaseRepo
//}
