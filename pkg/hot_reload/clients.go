package hot_reload

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/natewong1313/go-react-ssr/pkg/react_renderer"
)

var connectedClients = make(map[string][]*websocket.Conn)

func BroadcastFileUpdateToClients(filePath string) {
	fmt.Println("Broadcasting update to file", filePath)
	routes := react_renderer.GetRoutesForFile(filePath)
	for _, route := range routes {
		for k, ws := range connectedClients[route] {
			err := ws.WriteMessage(1, []byte("reload"))
			if err != nil {
				// remove client
				connectedClients[route] = append(connectedClients[route][:k], connectedClients[route][k+1:]...)
			}
		}
	}
}
