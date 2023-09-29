package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
	"log"
)

// App struct
type App struct {
	ctx    context.Context
	hotkey *hotkey.Hotkey
	conn   *websocket.Conn
}

// NewApp creates a new App application struct
func NewApp() *App {

	return &App{}
}

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

func (app *App) startup(ctx context.Context) {
	app.ctx = ctx

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/api/v1/ws?room=development&username=isitcorrect", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	app.conn = conn

	go func() {
		var wsMessage WsMessage
		for {
			if err := conn.ReadJSON(&wsMessage); err != nil {
				log.Println("read:", err)
				return
			}
			runtime.LogError(ctx, fmt.Sprintf("wsMessage: %v", wsMessage))

			switch wsMessage.Message.Type {
			case "generateAnswer.stream":
				runtime.EventsEmit(ctx, "generateAnswer.stream", wsMessage.Message.Message)
			case "generateAnswer.done":
				runtime.EventsEmit(ctx, "generateAnswer.stream", wsMessage.Message.Message)
			}
		}
	}()

	app.hotkey = hotkey.New([]hotkey.Modifier{hotkey.ModCmd, hotkey.ModShift}, hotkey.KeyC)
	if err := app.hotkey.Register(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			<-app.hotkey.Keyup()
			runtime.WindowShow(ctx) // Show the window
			question, err := app.GetSelectionText()
			if err != nil {
				fmt.Println("Error getting selection text: ", err.Error())
				continue
			}

			if question == "" {
				continue
			}

			app.SendMessage(question)
			runtime.EventsEmit(ctx, "question", question)
			runtime.EventsEmit(ctx, "isLoading", true)
		}
	}()

}

func (app *App) GetSelectionText() (string, error) {
	text, err := runtime.ClipboardGetText(app.ctx)
	if err != nil {
		return "", fmt.Errorf("error getting clipboard text: %v", err)
	}

	return text, nil
}

func (app *App) Shutdown(ctx context.Context) {

}

func (app *App) OnBeforeClose(ctx context.Context) bool {
	if app.hotkey != nil {
		app.hotkey.Unregister()
	}

	app.conn.Close()

	return false
}

func (app *App) SendMessage(message string) error {
	if app.conn == nil {
		return fmt.Errorf("no websocket connection")
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"type":      "dm",
		"recipient": "chatgpt",
		"message": map[string]interface{}{
			"type":           "generateAnswer",
			"conversationID": "abc",
			"message":        "check the grammar and spelling of the following sentence, Rewrite it if necessary, and explain any changes you make. Response in format: \n\nRewrite sentence: \n\nChanges made and why:\n\n The sentence to check is: " + message,
		},
	})
	return app.conn.WriteMessage(websocket.TextMessage, payload)
}

// IsConnected GetConnection returns the websocket connection
func (app *App) IsConnected() bool {
	return app.conn != nil
}

func (app *App) OnDomReady(ctx context.Context) {
	runtime.EventsEmit(ctx, "ext-connected", app.conn != nil)
}
