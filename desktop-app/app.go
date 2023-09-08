package main

import (
	"context"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
)

// App struct
type App struct {
	ctx      context.Context
	fiberApp *fiber.App
	conn     *websocket.Conn
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

		for {
			messageType, message, err := app.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}
				return
			}

			if messageType == websocket.TextMessage {
				runtime.EventsEmit(ctx, "message", string(message))
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

}

func (app *App) Shutdown(ctx context.Context) {

}

func (app *App) OnBeforeClose(ctx context.Context) bool {
	if err := app.fiberApp.Shutdown(); err != nil {
		fmt.Println("Error shutting down fiber app")
		return true
	}

	return false
}

func (app *App) SendMessage(message string) error {
	if app.conn == nil {
		return fmt.Errorf("no websocket connection")
	}
	return app.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// GetConnection returns the websocket connection
func (app *App) IsConnected() bool {
	return app.conn != nil
}

func (app *App) OnDomReady(ctx context.Context) {
	runtime.EventsEmit(ctx, "ext-connected", app.conn != nil)
}
