package ws

import (
	"encoding/json"
	"fmt"
	quizDto "github.com/quiz_be/services/core/application/quiz/dto"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type QuizRoomID string

type ClientConn struct {
	Conn *websocket.Conn
}

type WsHub struct {
	mu         sync.RWMutex
	rooms      map[QuizRoomID][]*ClientConn
	register   chan registerRequest
	unregister chan unregisterRequest
	broadcast  chan broadcastRequest
	//
	leaderboardCache map[string]*quizDto.GetLeaderboardRes
}

type registerRequest struct {
	RoomID QuizRoomID
	Client *ClientConn
}

type unregisterRequest struct {
	RoomID QuizRoomID
	Client *ClientConn
}

type broadcastRequest struct {
	RoomID QuizRoomID
	Data   []byte
}

// Create a new hub and start managing connections
func NewWsHub() *WsHub {
	hub := &WsHub{
		rooms:            make(map[QuizRoomID][]*ClientConn),
		register:         make(chan registerRequest),
		unregister:       make(chan unregisterRequest),
		broadcast:        make(chan broadcastRequest),
		leaderboardCache: make(map[string]*quizDto.GetLeaderboardRes),
	}
	go hub.run()
	return hub
}

// Run the central loop
func (h *WsHub) run() {
	for {
		select {
		case req := <-h.register:
			h.mu.Lock()
			h.rooms[req.RoomID] = append(h.rooms[req.RoomID], req.Client)
			h.mu.Unlock()

		case req := <-h.unregister:
			h.mu.Lock()
			clients := h.rooms[req.RoomID]
			for i, c := range clients {
				if c == req.Client {
					h.rooms[req.RoomID] = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			h.mu.Unlock()

		case req := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.rooms[req.RoomID] {
				err := client.Conn.WriteMessage(websocket.TextMessage, req.Data)
				if err != nil {
					log.Printf("Write error: %v", err)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Handle incoming WebSocket connection from client
func (h *WsHub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println(time.Now(), "\tHandleWebSocket")
	quizID := r.URL.Query().Get("quiz_id")
	if quizID == "" {
		log.Printf("Missing quiz_id")
		http.Error(w, "Missing quiz_id", http.StatusBadRequest)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		//http.Error(w, "Failed to upgrade to websocket", http.StatusInternalServerError)
		return
	}

	client := &ClientConn{Conn: conn}
	roomID := QuizRoomID(quizID)

	h.register <- registerRequest{RoomID: roomID, Client: client}
	defer func() {
		h.unregister <- unregisterRequest{RoomID: roomID, Client: client}
		conn.Close()
	}()

	// 👉 Gửi ping định kỳ từ server
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println("Ping error:", err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// 👉 Server sẽ reset deadline khi nhận pong
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			close(done)
			break
		}
		// Optional: xử lý message type
		var incoming map[string]interface{}
		if err := json.Unmarshal(message, &incoming); err == nil {
			switch incoming["type"] {
			case "ping":
				log.Println("Received ping: ", incoming)
				continue
			case "pong":
				log.Println("Received pong: ", incoming)
				continue
			case "leaderboard":
				log.Println("Received leaderboard: ", incoming, " quizID: ", quizID)
				log.Println("leaderboardCache: ", h.leaderboardCache)
				leaderboard, ok := h.leaderboardCache[quizID]
				log.Println("leaderboard from Cache: ", leaderboard, ok)
				if !ok {
					log.Println("No cached leaderboard found for quiz:", quizID)
					continue
				}

				respBytes, err := json.Marshal(map[string]interface{}{
					"leaderboard": leaderboard.Leaderboard,
					"page":        leaderboard.Page,
					"total":       leaderboard.Total,
				})
				if err != nil {
					log.Println("Marshal leaderboard error: ", err)
					continue
				}

				if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
					log.Println("Error sending leaderboard to WS: ", err)
				}
				continue
			}
		}
	}
}

// Call this when leaderboard updates (e.g. from gRPC handler)
func (h *WsHub) PushLeaderboardUpdate(quizID string, leaderboard *quizDto.GetLeaderboardRes) {
	fmt.Println(time.Now(), "\tPushLeaderboardUpdate: ", quizID, " leaderboard: ", leaderboard, " leaderboardCache: ", h.leaderboardCache)
	h.mu.Lock()
	h.leaderboardCache[quizID] = leaderboard
	h.mu.Unlock()
	data, err := json.Marshal(map[string]interface{}{
		"leaderboard": leaderboard.Leaderboard,
		"page":        leaderboard.Page,
		"total":       leaderboard.Total,
	})
	if err != nil {
		log.Printf("Marshal leaderboard error: %v", err)
		return
	}
	h.broadcast <- broadcastRequest{
		RoomID: QuizRoomID(quizID),
		Data:   data,
	}
}
