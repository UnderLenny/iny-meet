package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Client struct {
	conn   *websocket.Conn
	roomId string
}

type Hub struct {
	rooms map[string][]*Client
}

func server(hub *Hub) {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, HTTP!")
	})

	e.GET("/ws", wsHandler(hub))

	if err := e.Start(":1323"); err != nil {
		e.Logger.Fatal("failed to start server", err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(hub *Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		client := &Client{
			conn:   ws,
			roomId: "1",
		}
		
		hub.rooms[client.roomId] = append(hub.rooms[client.roomId], client)
		
		defer ws.Close()

		for {
			// Write
			err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
			if err != nil {
				c.Logger().Error("failed to write WS message", "error", err)
				break
			}

			// Read
			_, msg, err := ws.ReadMessage()
			if err != nil {
				c.Logger().Error("failed to read WS message", "error", err)
				break
			}
			fmt.Printf("%s\n", msg)
		}

		return nil
	}

}

func main() {
	hub := &Hub{
		rooms: make(map[string][]*Client),
	}
	server(hub)
}
