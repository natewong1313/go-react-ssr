package hot_reload

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Serve(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, clientRoute, err := ws.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ws.WriteMessage(1, []byte("Connected"))
	if err != nil {
		fmt.Println(err)
		return
	}
	connectedClients[string(clientRoute)] = append(connectedClients[string(clientRoute)], ws)
}
