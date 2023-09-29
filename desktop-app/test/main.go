package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type WsMessage struct {
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Type      string    `json:"type"`
	Message   AiMessage `json:"message"`
}

type AiMessage struct {
	ConversationID string `json:"conversationID"`
	Message        string `json:"message"`
	Type           string `json:"type"`
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/api/v1/ws?room=development&username=isitcorrect", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	var wsMessage WsMessage
	for {
		if err := conn.ReadJSON(&wsMessage); err != nil {
			log.Println("read:", err)
			return
		}

		log.Println(wsMessage)
	}
}
