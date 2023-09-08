package hot_reload

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Serve(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ws.WriteMessage(1, []byte("Connected"))
	if err != nil {
		fmt.Println(err)
		return
	}
	connectionID := uuid.New()
	connectedClients[connectionID] = ws
}