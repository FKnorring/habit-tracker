package sockets

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)         // Connected clients
var clientUserMap = make(map[*websocket.Conn]string) // Map client connections to user IDs
var broadcast = make(chan []byte)                    // Broadcast channel
var mutex = &sync.Mutex{}                            // Protect clients map

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type AuthData struct {
	UserID string `json:"userId"`
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebSocket connection attempt from: %s", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer func() {
		conn.Close()
		mutex.Lock()
		delete(clients, conn)
		delete(clientUserMap, conn)
		clientCount := len(clients)
		mutex.Unlock()
		log.Printf("Client disconnected and cleaned up. Total clients: %d", clientCount)
	}()

	mutex.Lock()
	clients[conn] = true
	clientCount := len(clients)
	mutex.Unlock()

	log.Printf("Client connected. Total clients: %d", clientCount)

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		if msg.Type == "auth" {
			authDataBytes, err := json.Marshal(msg.Data)
			if err != nil {
				log.Printf("Error marshaling auth data: %v", err)
				continue
			}

			var authData AuthData
			if err := json.Unmarshal(authDataBytes, &authData); err != nil {
				log.Printf("Error parsing auth data: %v", err)
				continue
			}

			mutex.Lock()
			clientUserMap[conn] = authData.UserID
			mutex.Unlock()

			log.Printf("Client authenticated with user ID: %s", authData.UserID)
			continue
		}

		log.Printf("Received message: %s", string(messageBytes))
		broadcast <- messageBytes
	}
}

func HandleMessages() {
	log.Println("WebSocket message handler started")
	for {

		message := <-broadcast
		log.Printf("Broadcasting message to %d clients: %s", len(clients), string(message))

		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error sending message to client: %v", err)
				client.Close()
				delete(clients, client)
				delete(clientUserMap, client)
			}
		}
		mutex.Unlock()
	}
}

func MessageUser(userID string, message []byte) error {
	mutex.Lock()
	defer mutex.Unlock()

	var targetConn *websocket.Conn
	for conn, connUserID := range clientUserMap {
		if connUserID == userID {
			targetConn = conn
			break
		}
	}

	if targetConn == nil {
		log.Printf("User %s not found or not connected", userID)
		return nil
	}

	err := targetConn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Error sending message to user %s: %v", userID, err)

		targetConn.Close()
		delete(clients, targetConn)
		delete(clientUserMap, targetConn)
		return err
	}

	log.Printf("Message sent to user %s: %s", userID, string(message))
	return nil
}
