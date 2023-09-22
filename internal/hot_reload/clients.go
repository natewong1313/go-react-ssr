package hot_reload

import (
	"github.com/gorilla/websocket"
	"github.com/natewong1313/go-react-ssr/react_renderer"
)

// Each "client" is a websocket connection that is listening for file updates for the given routeID
var connectedClients = make(map[string][]*websocket.Conn)

// Tell all clients listening for a specific routeID to reload
func BroadcastFileUpdateToClients(filePath string) {
	// Find which routeIDS have this file as a dependency or parent file
	routeIDS := react_renderer.GetRouteIDSForFile(filePath)
	for _, routeID := range routeIDS {
		for k, ws := range connectedClients[routeID] {
			// Send reload message to client
			err := ws.WriteMessage(1, []byte("reload"))
			if err != nil {
				// remove client if browser is closed or page changed
				connectedClients[routeID] = append(connectedClients[routeID][:k], connectedClients[routeID][k+1:]...)
			}
		}
	}
}
