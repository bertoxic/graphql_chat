package chats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/database"
	"github.com/bertoxic/graphqlChat/internal/database/postgres"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const (
	maxMessageSize = 1024 * 1024
	pongWait       = 60 * time.Second
	writeWait      = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
)

type Client struct {
	Config config.AppConfig
	ID     string
	User   models.User
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	UnRegister chan *Client
	Redis      database.RedisClient
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	Send       chan []byte
	Private    chan *models.Message
	Public     chan *models.Message
	DB         database.DatabaseRepo
}

type HubInterface interface {
	StoreUnreadMessage(ctx context.Context, msg *models.Message) error
	GetUnreadMessages(ctx context.Context, userID string) ([]models.Message, error)
	UpdateUserPresence(ctx context.Context, userID string) error
	IsUserOnline(ctx context.Context, userID string) bool
	GetFollowers(ctx context.Context, userID string) ([]models.User, error)
	UpdateFollowerCounts(ctx context.Context, userID, followerID string, isFollow bool) error
	RemoveUserPresence(ctx context.Context, userID string) error
}

func NewHub(dB database.DatabaseRepo, rdb database.RedisClient) HubInterface {
	ctx, cancel := context.WithCancel(context.Background())
	hub := &Hub{
		Clients:    make(map[string]*Client),
		Private:    make(chan *models.Message),
		Public:     make(chan *models.Message),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Redis:      rdb,
		DB:         dB,
		ctx:        ctx,
		cancel:     cancel,
	}

	// Start the hub's main loop
	go hub.run()
	return hub
}

func (h *Hub) run() {
	//defer h.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return

		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.UnRegister:
			h.unregisterClient(client)

		case message := <-h.Public:
			h.broadcastMessage(encodeMessage(message))

		case message := <-h.Private:
			h.handlePrivateMessage(message)

		case message := <-h.Public:
			h.handlePublicMessage(message)
		}
	}

}
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	h.Clients[client.ID] = client
	h.mu.Unlock()
	if err := h.UpdateUserPresence(h.ctx, client.User.ID); err != nil {
		fmt.Printf("an error occured %w", err)
	}
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	if _, exists := h.Clients[client.ID]; exists {
		delete(h.Clients, client.ID)
	}
	h.mu.Unlock()
	if err := h.RemoveUserPresence(h.ctx, client.User.ID); err != nil {
	}
}
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			go h.unregisterClient(client)
		}
	}
}
func (h *Hub) handlePrivateMessage(message *models.Message) {
	if message == nil {
		return
	}

	h.mu.RLock()
	recipient, exists := h.Clients[message.ToID]
	h.mu.RUnlock()

	if exists {
		msgBytes, err := json.Marshal(message)
		if err != nil {
			return
		}

		select {
		case recipient.Send <- msgBytes:
		default:
			go h.StoreUnreadMessage(h.ctx, message)
		}
	} else {
		go h.StoreUnreadMessage(h.ctx, message)
	}
}

func (h *Hub) handlePublicMessage(message *models.Message) {

	if message == nil {
		log.Printf("there is no message availablw")
		return
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("cannot marshall message")
		return
	}

	h.broadcastMessage(msgBytes)
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.UnRegister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &models.Message{
			ID:        uuid.New().String(),
			FromID:    c.User.ID,
			Timestamp: time.Now(),
			Status:    "sent",
		}

		var incoming struct {
			Type        string `json:"type"`
			Content     string `json:"content"`
			RecipientID string `json:"to_id,omitempty"`
			MessageType string `json:"messageType"`
		}

		if err := json.Unmarshal(message, &incoming); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		msg.Content = incoming.Content
		msg.Type = incoming.MessageType
		switch incoming.Type {
		case "private":
			msg.ToID = incoming.RecipientID
			if msg.ToID == "" {
				log.Printf("private message missing recipient ID")
				continue
			}
			c.Hub.Private <- msg
		case "public":
			c.Hub.Public <- msg
		default:
			log.Printf("unknown message type: %s", incoming.Type)
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.UnRegister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, payload, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		msg := decodeMessage(payload)
		switch msg.Type {
		case "text":
			c.Hub.Private <- msg
		case "post":
			c.Hub.Public <- msg
		}
	}
}
func encodePost(post models.Post) []byte {
	data, err := json.Marshal(post)
	if err != nil {
		log.Printf("Error encoding post: %v", err)
		return nil
	}
	return data
}

func encodeMessage(msg *models.Message) []byte {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error encoding message: %v", err)
		return nil
	}
	return data
}

const (
	userPresenceKey   = "presence:%s"
	unreadMsgKey      = "unread:%s"
	messageHistoryKey = "history:%s:%s"
	userFollowersKey  = "followers:%s"
	messageExpiry     = 24 * time.Hour
	presenceExpiry    = 5 * time.Minute
)

func (h *Hub) StoreUnreadMessage(ctx context.Context, msg *models.Message) error {
	log.Printf("about to store message in database: %v", msg)
	unreadKey := fmt.Sprintf(unreadMsgKey, msg.ToID)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	err = h.Redis.Client.RPush(h.ctx, unreadKey, msgBytes).Err()
	if err != nil {
		log.Printf("Failed to store message in Redis: %v", err)
	}
	h.Redis.Client.Expire(h.ctx, unreadKey, messageExpiry)

	db, ok := h.DB.(*postgres.PostgresDBRepo)
	if !ok {
		log.Printf("Database type assertion failed") // Add logging
		return fmt.Errorf("pr.Repo does not implement database.Database")
	}

	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO messages (id, from_user_id, to_user_id, content, created_at, message_type) 
             VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = db.DB.Exec(ctx, query,
		msg.ID,
		msg.FromID,
		msg.ToID,
		msg.Content,
		msg.Timestamp,
		msg.Type,
	)

	if err != nil {
		log.Printf("Failed to store message in database: %v", err)
		return err
	}

	return nil
}

func (h *Hub) GetUnreadMessages(ctx context.Context, id string) ([]models.Message, error) {
	var messages []models.Message
	unreadKey := fmt.Sprintf(unreadMsgKey, id)

	msgList, err := h.Redis.Client.LRange(h.ctx, unreadKey, 0, -1).Result()
	if err == nil && len(msgList) > 0 {
		for _, msgStr := range msgList {
			var msg models.Message
			if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
				continue
			}
			messages = append(messages, msg)
		}
		return messages, nil
	}

	db, ok := h.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.Repo does not implement database.Database")
	}

	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, from_user_id, to_user_id, content, created_at, message_type 
             FROM messages 
             WHERE to_user_id = $1 AND is_read = false
             ORDER BY created_at ASC`

	rows, err := db.DB.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread messages: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.ID, &msg.FromID, &msg.ToID, &msg.Content, &msg.Timestamp, &msg.Type)
		if err != nil {
			continue
		}
		messages = append(messages, msg)

		// Cache in Redis
		if msgBytes, err := json.Marshal(msg); err == nil {
			h.Redis.Client.RPush(h.ctx, unreadKey, msgBytes)
		}
	}

	h.Redis.Client.Expire(h.ctx, unreadKey, messageExpiry)
	return messages, nil
}

func (h *Hub) UpdateUserPresence(ctx context.Context, userID string) error {
	presenceKey := fmt.Sprintf(userPresenceKey, userID)
	return h.Redis.Client.Set(h.ctx, presenceKey, "online", presenceExpiry).Err()
}

func (h *Hub) IsUserOnline(ctx context.Context, userID string) bool {
	presenceKey := fmt.Sprintf(userPresenceKey, userID)
	result, err := h.Redis.Client.Get(h.ctx, presenceKey).Result()
	return err == nil && result == "online"
}

func (h *Hub) GetFollowers(ctx context.Context, id string) ([]models.User, error) {
	var followers []models.User
	followersKey := fmt.Sprintf(userFollowersKey, id)

	// Try to get from Redis first
	followersData, err := h.Redis.Client.Get(h.ctx, followersKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(followersData), &followers)
		if err == nil {
			return followers, nil
		}
	}

	// If not in Redis, get from DB
	db, ok := h.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.Repo does not implement database.Database")
	}

	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()

	query := `
        SELECT u.id, u.username, u.email, u.full_name, u.profile_picture_url, 
               u.is_verified, u.follower_count, u.following_count
        FROM follows f 
        JOIN users u ON f.follower_id = u.id 
        WHERE f.follower_id = $1
        ORDER BY u.username`

	rows, err := db.DB.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.UserName,
			&user.Email,
			&user.FullName,
			&user.ProfilePictureURL,
			&user.IsVerified,
			&user.Followers,
			&user.Followings,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan follower: %v", err)
		}
		followers = append(followers, user)
	}

	// Cache in Redis
	if followersBytes, err := json.Marshal(followers); err == nil {
		h.Redis.Client.Set(h.ctx, followersKey, followersBytes, time.Hour)
	}

	return followers, nil
}

func (h *Hub) RemoveUserPresence(ctx context.Context, userID string) error {
	presenceKey := fmt.Sprintf(userPresenceKey, userID)
	return h.Redis.Client.Del(h.ctx, presenceKey).Err()
}
func (h *Hub) UpdateFollowerCounts(ctx context.Context, userID, followerID string, isFollow bool) error {
	db, ok := h.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return fmt.Errorf("pr.Repo does not implement database.Database")
	}

	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()

	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	increment := 1
	if !isFollow {
		increment = -1
	}

	// Update follower and following counts
	queries := []string{
		`UPDATE users 
         SET follower_count = follower_count + $1 
         WHERE id = $2`,
		`UPDATE users 
         SET following_count = following_count + $1 
         WHERE id = $2`,
	}

	for _, query := range queries {
		_, err = tx.Exec(ctx, query, increment, userID)
		if err != nil {
			return fmt.Errorf("failed to update counts: %w", err)
		}
	}

	// Update follows table based on isFollow flag
	if isFollow {
		// Add a new follow relationship
		_, err = tx.Exec(ctx, `
			INSERT INTO follows (followed_id, follower_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, userID, followerID)
		if err != nil {
			return fmt.Errorf("failed to insert follow relationship: %w", err)
		}
	} else {
		// Remove the follow relationship
		_, err = tx.Exec(ctx, `
			DELETE FROM follows
			WHERE followed_id = $1 AND follower_id = $2
		`, userID, followerID)
		if err != nil {
			return fmt.Errorf("failed to delete follow relationship: %w", err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Invalidate the followers cache in Redis
	followersKey1 := fmt.Sprintf(userFollowersKey, userID)
	followersKey2 := fmt.Sprintf(userFollowersKey, followerID)
	h.Redis.Client.Del(h.ctx, followersKey1, followersKey2)

	return nil
}

func (c *Client) HandleConnection() {

	go c.ReadPump()
	go c.writePump()
}

func decodeMessage(payload []byte) *models.Message {
	var msg models.Message
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		log.Printf("Error decoding message: %v", err)
		return nil
	}
	return &msg
}
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			c.Conn.WriteMessage(websocket.PingMessage, nil)

		}
	}
}
