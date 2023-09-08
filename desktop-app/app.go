package main

import (
	"context"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
	"log"
	"os"
	"os/exec"
)

// App struct
type App struct {
	ctx      context.Context
	fiberApp *fiber.App
	conn     *websocket.Conn
	hotkey   *hotkey.Hotkey
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (app *App) startup(ctx context.Context) {
	app.ctx = ctx

	app.fiberApp = fiber.New()

	app.fiberApp.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	app.fiberApp.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		app.conn = conn
		defer func() {
			if err := conn.Close(); err != nil {
				runtime.EventsEmit(ctx, "websocket", "error", err.Error())
			}
			runtime.EventsEmit(ctx, "ext-connected", false)
			app.conn = nil
		}()

		runtime.EventsEmit(ctx, "ext-connected", true)
		app.conn.WriteMessage(websocket.TextMessage, []byte("setSystemMessage"))

		for {
			messageType, message, err := app.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}
				return
			}

			if messageType == websocket.TextMessage {
				textMessage := string(message)
				if textMessage == ":done:" {
					runtime.EventsEmit(ctx, "message-end", false)
				} else {
					runtime.EventsEmit(ctx, "message", textMessage)
				}
			} else {
				runtime.EventsEmit(ctx, "message", "websocket message received of type "+string(rune(messageType)))
			}
		}
	}))

	go func() {
		if err := app.fiberApp.Listen("localhost:8991"); err != nil {
			fmt.Println("Error starting fiber app")
		}
	}()

	app.hotkey = hotkey.New([]hotkey.Modifier{hotkey.ModCmd, hotkey.ModShift}, hotkey.KeyC)
	if err := app.hotkey.Register(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			<-app.hotkey.Keyup()
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
			runtime.WindowShow(ctx) // Show the
		}
	}()

}

func (app *App) GetSelectionText() (string, error) {
	// copy to clipboard and print to stdout
	script := `
		 tell application "System Events" to keystroke "c" using command down
		 delay 0.1
		 set the clipboard to (the clipboard as text)
		 get the clipboard
	`
	c := exec.Command("/usr/bin/osascript", "-e", script)
	c.Stdin = os.Stdin
	output, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (app *App) Shutdown(ctx context.Context) {

}

func (app *App) OnBeforeClose(ctx context.Context) bool {
	if err := app.fiberApp.Shutdown(); err != nil {
		fmt.Println("Error shutting down fiber app")
		return true
	}

	if app.hotkey != nil {
		app.hotkey.Unregister()
	}

	return false
}

func (app *App) SendMessage(message string) error {
	if app.conn == nil {
		return fmt.Errorf("no websocket connection")
	}
	return app.conn.WriteMessage(websocket.TextMessage, []byte("check the grammar, rewrite to address any issues, provide a brief explanation of its structure:\n\n"+message))
}

// GetConnection returns the websocket connection
func (app *App) IsConnected() bool {
	return app.conn != nil
}

func (app *App) OnDomReady(ctx context.Context) {
	runtime.EventsEmit(ctx, "ext-connected", app.conn != nil)
}
