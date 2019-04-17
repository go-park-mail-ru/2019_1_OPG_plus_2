package gameservice

import (
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

func AddGameServicePaths(router *mux.Router) *mux.Router {
	hub := NewHub()
	err := hub.AttachRooms(newRoom(hub, 0))
	if err != nil {
		panic("WTF")
	}

	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			fmt.Println("could not parse ", id)
			return
		}
		if hub.rooms[int(id)] == nil {
			_, _ = fmt.Fprint(w, "no such room with id ", id)
			return
		}
		http.ServeFile(w, r, "home.html")
	}).Methods("GET")

	router.HandleFunc("/{id}/room", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			fmt.Println("could not parse ", id)
			return
		}
		fmt.Println("CONN", r.URL)
		serveClientConnection(hub.rooms[int(id)], w, r)
	}).Methods("GET")

	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			fmt.Println("could not parse ", id)
			return
		}
		hub.closeRoom(int(id))

		_, _ = fmt.Fprint(w, "Room ", id, " closing")
	}).Methods("DELETE")

	router.HandleFunc("/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			fmt.Println("could not parse ", id)
			return
		}

		err = hub.AttachRooms(newRoom(hub, int(id)))
		if err != nil {
			_, _ = fmt.Fprint(w, err)
			return
		}
		_, _ = fmt.Fprint(w, "Room ", id, " created")
	}).Methods("CREATE")

	go hub.run()

	return router

}

// serveClientConnection handles websocket requests from the peer.
func serveClientConnection(room *Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {

		_, _ = fmt.Fprintln(w, err)
		return
	}
	client := NewClient(room, conn)
	client.room.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
