package hot_reload

import (
	"fmt"
	"gossr/pkg/react_renderer"

	"github.com/gorilla/websocket"
)

var connectedClients = make(map[string][]*websocket.Conn)

func BroadcastFileUpdateToClients(filePath string) {
	fmt.Println("Broadcasting update to file", filePath)
	routes := react_renderer.GetRoutesForFile(filePath)
	for _, route := range routes {
		for _, ws := range connectedClients[route] {
			fmt.Println("OK")
			err := ws.WriteMessage(1, []byte("reload"))
			if err != nil {
				fmt.Println(err)
				// if errors.Is(err, syscall.EPIPE) {
				// 	// remove client
				// 	for k, v := range connectedClients {
				// 		if v == ws {
				// 			delete(connectedClients, k)
				// 			break
				// 		}
				// 	}
				// }else {
				// 	fmt.Println(err)
				// }
			}
		}
	}
}
