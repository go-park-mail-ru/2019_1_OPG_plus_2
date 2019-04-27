package chatservice

import (
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *ChatService) AddChatServicePaths(router *mux.Router) *mux.Router {
	router.HandleFunc("/chat/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			tsLogger.LogWarn("could not parse %d", id)
			return
		}
		if s.Hub.rooms[int(id)] == nil {
			_, _ = fmt.Fprint(w, "no such room with id ", id)
			return
		}
		http.ServeFile(w, r, "home.html")
	}).Methods("GET")

	router.HandleFunc("/chat/{id}/room", s.ConnectToRoom)

	router.HandleFunc("/chat/{id}", s.DeleteRoom).Methods("DELETE")

	router.HandleFunc("/chat/{id}", s.CreateRoom).Methods("CREATE")
	return router
}

func serveClientConnection(room *ChatRoom, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		tsLogger.LogErr("CONNECTION UPGRADE ERROR: %s", err)
		return err
	}
	client := NewClient(room, conn)
	client.room.register <- client

	go client.writePump()
	go client.readPump()
	return nil
}
