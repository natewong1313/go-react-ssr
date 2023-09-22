package hot_reload

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/internal/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartServer() {
	logger.L.Info().Msgf("Serving hot reload websocket at port %d", config.C.HotReloadServerPort)
	http.HandleFunc("/ws", serve)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.C.HotReloadServerPort), nil)
	if err != nil {
		logger.L.Error().Err(err).Msg("Hot reload server quit unexpectedly")
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to upgrade websocket")
		return
	}
	// Client should send routeID as first message
	_, routeID, err := ws.ReadMessage()
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to read message from websocket")
		return
	}
	err = ws.WriteMessage(1, []byte("Connected"))
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to write message to websocket")
		return
	}
	// Add client to connectedClients
	connectedClients[string(routeID)] = append(connectedClients[string(routeID)], ws)
}
