package hot_reload

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/natewong1313/go-react-ssr/pkg/config"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartServer() {
	http.HandleFunc("/ws", serve)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", config.C.HotReloadServerPort), nil))
}

func serve(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
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
