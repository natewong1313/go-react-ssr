package hot_reload

import (
	"errors"
	"fmt"
	"gossr/pkg/react_renderer"
	"syscall"

	"github.com/gorilla/websocket"
)

var connectedClients = make(map[string][]*websocket.Conn)

func BroadcastFileUpdateToClients(filePath string) {
	fmt.Println("Broadcasting update to file", filePath)
	routes := react_renderer.GetRoutesForFile(filePath)
	for _, route := range routes {
		for k, ws := range connectedClients[route] {
			err := ws.WriteMessage(1, []byte("reload"))
			if err != nil {
				if errors.Is(err, syscall.EPIPE) {
					// remove client
					connectedClients[route] = append(connectedClients[route][:k], connectedClients[route][k+1:]...)
				}
			}
		}
	}
}
