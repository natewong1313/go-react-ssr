package hot_reload

import (
	"github.com/gorilla/websocket"
	"github.com/natewong1313/go-react-ssr/pkg/react_renderer"
)

var connectedClients = make(map[string][]*websocket.Conn)

func BroadcastFileUpdateToClients(filePath string) {
	routeIDS := react_renderer.GetRouteIDSForFile(filePath)
	for _, routeID := range routeIDS {
		for k, ws := range connectedClients[routeID] {
			err := ws.WriteMessage(1, []byte("reload"))
			if err != nil {
				// remove client
				connectedClients[routeID] = append(connectedClients[routeID][:k], connectedClients[routeID][k+1:]...)
			}
		}
	}
}
