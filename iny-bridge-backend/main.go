package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


func server() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Сервер запущен")
	})

	e.Start(":8000")
}

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка обновлений", err)
		return
	}
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Ошибка чтения сообщения:", err)
			break
		}
  		fmt.Printf("Получено: %s\\n", message)
    if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
    	fmt.Println("Ошибка записи сообщения:", err)
     	break
    }

	}
}


func main() {

	server()
}
