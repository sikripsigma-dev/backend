package handler

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"

	"Skripsigma-BE/internal/config"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
)

type Client struct {
	Conn   *websocket.Conn
	RoomID string
}

var (
	clients = make(map[string][]*Client) // room_id -> clients
	mu      sync.Mutex
)

func ChatWebSocketHandler(c *websocket.Conn) {
	roomID := c.Params("room_id")

	client := &Client{
		Conn:   c,
		RoomID: roomID,
	}

	mu.Lock()
	clients[roomID] = append(clients[roomID], client)
	log.Printf("Client joined room %s (total: %d)", roomID, len(clients[roomID]))
	mu.Unlock()

	defer func() {
		mu.Lock()
		for i, cl := range clients[roomID] {
			if cl.Conn == c {
				clients[roomID] = append(clients[roomID][:i], clients[roomID][i+1:]...)
				break
			}
		}
		log.Printf("Client left room %s (remaining: %d)", roomID, len(clients[roomID]))
		mu.Unlock()

		c.Close()
	}()

	for {
		var msg struct {
			SenderID   string `json:"sender_id"`
			SenderType string `json:"sender_type"`
			Message    string `json:"message,omitempty"`
			Action     string `json:"action,omitempty"`  // NEW
		}
	
		if err := c.ReadJSON(&msg); err != nil {
			log.Println("read error:", err)
			break
		}
	
		if msg.Action == "mark_read" {
			// Saat client buka room
			oppositeType := "student"
			if msg.SenderType == "student" {
				oppositeType = "company"
			}
			if err := repository.MarkMessagesAsRead(config.DB, roomID, oppositeType); err != nil {
				log.Println("failed to mark messages as read:", err)
			}
	
			// Broadcast read update
			mu.Lock()
			for _, cl := range clients[roomID] {
				if cl.Conn != nil {
					cl.Conn.WriteJSON(map[string]interface{}{
						"type":         "read_update",
						"room_id":      roomID,
						"by":           msg.SenderType,
						"unread_count": 0,
					})
				}
			}
			mu.Unlock()
			continue
		}
	
		// Handle pesan baru
		chatMsg := &models.ChatMessage{
			ID:         uuid.New().String(),
			RoomID:     roomID,
			SenderID:   &msg.SenderID,
			SenderType: msg.SenderType,
			Message:    msg.Message,
			IsRead:     false,
			CreatedAt:  time.Now(),
		}
	
		if err := repository.SaveChatMessage(config.DB, chatMsg); err != nil {
			log.Println("failed to save message:", err)
			continue
		}
	
		// Mark lawan read (optional atau hapus kalau sudah pakai mark_read di FE)
		oppositeType := "student"
		if msg.SenderType == "student" {
			oppositeType = "company"
		}
		if err := repository.MarkMessagesAsRead(config.DB, roomID, oppositeType); err != nil {
			log.Println("failed to mark messages as read:", err)
		}
	
		// Broadcast message
		mu.Lock()
		for _, cl := range clients[roomID] {
			if cl.Conn != nil {
				cl.Conn.WriteJSON(map[string]interface{}{
					"type":         "new_message",
					"id":           chatMsg.ID,
					"room_id":      chatMsg.RoomID,
					"sender_id":    chatMsg.SenderID,
					"sender_type":  chatMsg.SenderType,
					"message":      chatMsg.Message,
					"is_read":      chatMsg.IsRead,
					"created_at":   chatMsg.CreatedAt,
				})
			}
		}
		mu.Unlock()
	}
	
}
