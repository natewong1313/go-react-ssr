package hot_reload

import (
	"github.com/gorilla/websocket"
	// "github.com/natewong1313/go-react-ssr/react_renderer"
)

// Each "client" is a websocket connection that is listening for file updates for the given routeID
// Maps routeID's to a slice of websocket connections
var connectedClients = make(map[string][]*websocket.Conn)

// Tell all clients listening for a specific routeID to reload
func BroadcastFileUpdateToClients(routeIDS []string) {
	// Iterate over each route ID
	for _, routeID := range routeIDS {
		// Find all clients listening for that route ID
		for i, ws := range connectedClients[routeID] {
			// Send reload message to client
			err := ws.WriteMessage(1, []byte("reload"))
			if err != nil {
				// remove client if browser is closed or page changed
				connectedClients[routeID] = append(connectedClients[routeID][:i], connectedClients[routeID][i+1:]...)
			}
		}
	}
}
