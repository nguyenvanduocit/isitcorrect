package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
	"io/ioutil"
	"log"
	"os"
)

type Settings struct {
	RoomID     string `json:"room_id"`
	WsEndpoint string `json:"ws_endpoint"`
}

// App struct
type App struct {
	ctx      context.Context
	hotkey   *hotkey.Hotkey
	conn     *websocket.Conn
	settings Settings
}

// NewApp creates a new App application struct
func NewApp() *App {

	return &App{
		settings: Settings{},
	}
}

// StartWsHandler
func (app *App) StartWsHandler() {
	var buf []byte
	var err error
	var gMessage gjson.Result
	for {
		if _, buf, err = app.conn.ReadMessage(); err != nil {
			log.Println("read:", err)
			panic(err)
		}

		gMessage = gjson.ParseBytes(buf)
		messageType := gMessage.Get("type").String()

		switch messageType {
		case "generateAnswer.stream":
			runtime.EventsEmit(app.ctx, "generateAnswer.stream", gMessage.Get("message").String())
		case "generateAnswer.done":
			runtime.EventsEmit(app.ctx, "generateAnswer.done", "")
		}
	}
}

func (app *App) startup(ctx context.Context) {
	app.ctx = ctx

	if err := app.LoadSettings(); err != nil {
		runtime.LogError(ctx, "error loading settings: "+err.Error())
		log.Fatal(err)
	}

	if roomID := app.GetRoomID(); roomID == "" {
		runtime.EventsEmit(ctx, "message", "no room ID saved")
	} else if err := app.StartWsConnection(app.settings.WsEndpoint, roomID); err != nil {
		runtime.LogError(ctx, "error starting websocket connection: "+err.Error())
		log.Fatal(err)
	} else {
		app.StartWsHandler()
	}

	runtime.EventsOn(ctx, "WebsocketDisconnected", func(args ...interface{}) {
		app.conn = nil
	})

	runtime.EventsOn(ctx, "roomIdSaved", func(args ...interface{}) {
		if app.conn != nil {
			app.conn.Close()
		}

		if err := app.StartWsConnection(app.settings.WsEndpoint, app.GetRoomID()); err != nil {
		}
	})

	app.hotkey = hotkey.New([]hotkey.Modifier{hotkey.ModCmd, hotkey.ModShift}, hotkey.KeyC)
	if err := app.hotkey.Register(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			<-app.hotkey.Keyup()
			runtime.WindowShow(ctx)
			question, err := app.GetSelectionText()
			if err != nil {
				runtime.EventsEmit(ctx, "message", "error getting selection text: "+err.Error())
				continue
			}

			if question == "" {
				runtime.EventsEmit(ctx, "message", "no text selected")
				continue
			}

			runtime.EventsEmit(ctx, "question", question)

			if err := app.GenerateAnswer(question); err != nil {
				runtime.EventsEmit(ctx, "message", "error generating answer: "+err.Error())
				continue
			}
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
	app.SaveSettings()
}

func (app *App) OnBeforeClose(ctx context.Context) bool {
	if app.hotkey != nil {
		app.hotkey.Unregister()
	}

	app.conn.Close()

	return false
}

func (app *App) GenerateAnswer(message string) error {
	if app.conn == nil {
		return fmt.Errorf("no websocket connection")
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"type":      "dm",
		"recipient": "chatgpt",
		"message": map[string]interface{}{
			"type":    "generateAnswer",
			"message": "check the grammar and spelling of the following sentence, Rewrite it if necessary, and explain any changes you make. Response in format: \n\nRewrite sentence: \n\nChanges made and why:\n\n The sentence to check is: " + message,
		},
	})
	return app.conn.WriteMessage(websocket.TextMessage, payload)
}

// SetSystemMessage SetSystemMessage
func (app *App) SetSystemMessage() error {
	if app.conn == nil {
		return fmt.Errorf("no websocket connection")
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"type":      "dm",
		"recipient": "chatgpt",
		"message": map[string]interface{}{
			"type": "setSystemMessage",
			"message": map[string]string{
				"about_user_message":  "I am learning English and would like to improve my writing skills. Please help me by correcting my grammar and spelling. Thank you!",
				"about_model_message": "I am a GPT-3 model trained on the English language. I can help you improve your writing skills by correcting your grammar and spelling.",
				"enabled":             "true",
			},
		},
	})
	return app.conn.WriteMessage(websocket.TextMessage, payload)
}

// IsConnected GetConnection returns the websocket connection
func (app *App) IsConnected() bool {
	return app.conn != nil
}

func (app *App) OnDomReady(ctx context.Context) {

}

func (app *App) StartWsConnection(wsEndpoint, roomID string) error {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint+"?room="+roomID+"&username=isitcorrect", nil)
	if err != nil {
		return err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		runtime.EventsEmit(app.ctx, "message", "connection closed")
		app.conn = nil
		return nil
	})

	app.conn = conn

	return nil
}

// GetRoomID returns the room ID
func (app *App) GetRoomID() string {
	return app.settings.RoomID
}

// SetRoomID sets the room ID
func (app *App) SetRoomID(roomID string) {
	app.settings.RoomID = roomID
	if err := app.SaveSettings(); err != nil {
		runtime.EventsEmit(app.ctx, "message", "error saving settings: "+err.Error())
		return
	}
	runtime.EventsEmit(app.ctx, "message", roomID+" room ID saved")
	runtime.EventsEmit(app.ctx, "roomIdSaved")
}

// LoadSettings loads the settings
func (app *App) LoadSettings() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Create the user settings file path.
	settingsFilePath := homeDir + "/.isitcorrect-settings.json"

	// Check if the settings file exists.
	if _, err := os.Stat(settingsFilePath); err != nil {
		// If the file doesn't exist, create it.
		if os.IsNotExist(err) {
			if err := app.SaveSettings(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Read the settings file.
	content, err := ioutil.ReadFile(settingsFilePath)
	if err != nil {
		return err
	}

	// Unmarshal the settings file.
	if err := json.Unmarshal(content, &app.settings); err != nil {
		return err
	}

	if app.settings.WsEndpoint == "" {
		app.settings.WsEndpoint = "wss://aibridge.fly.dev/api/v1/ws"
	}

	return nil
}

// Init settings
func (app *App) SaveSettings() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Create the user settings file path.
	settingsFilePath := homeDir + "/.isitcorrect-settings.json"

	content, _ := json.Marshal(app.settings)
	err = ioutil.WriteFile(settingsFilePath, content, 0600)
	if err != nil {
		return err
	}

	return nil
}
