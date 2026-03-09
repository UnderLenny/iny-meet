package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func server() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, HTTP!")
	})

	e.GET("/ws", wsHandler)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Fatal("failed to start server", err)
	}
}

var upgrader = websocket.Upgrader{}

func wsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error("failed to write WS message", "error", err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error("failed to read WS message", "error", err)
		}
		fmt.Printf("%s\n", msg)
	}
}

func main() {
	server()
}
