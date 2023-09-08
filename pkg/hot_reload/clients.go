package hot_reload

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var connectedClients = make(map[uuid.UUID]*websocket.Conn)


func BroadcastFileUpdateToClients() {
	for _, ws := range connectedClients {
		err := ws.WriteMessage(1, []byte("reload"))
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				// remove client
				for k, v := range connectedClients {
					if v == ws {
						delete(connectedClients, k)
						break
					}
				}
			}else {
				fmt.Println(err)
			}
		}
	}
}
