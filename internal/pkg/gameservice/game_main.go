package gameservice

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

func AddGameServicePaths(router *mux.Router) *mux.Router {
	hub := NewHub()
	err := hub.AttachRooms(newRoom(hub, 0))
	if err != nil {
		tsLogger.LogErr("ROOM ATTACHMENT ERROR: %v", hub.rooms)
		panic("WTF")
	}
	tsLogger.LogTrace("INITIAL ROOM CREATED")

	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			tsLogger.LogWarn("could not parse %d", id)
			return
		}
		if hub.rooms[int(id)] == nil {
			_, _ = fmt.Fprint(w, "no such room with id ", id)
			return
		}
		//http.ServeFile(w, r, "home.html")
		_, _ = fmt.Fprint(w, "IDI NAHUY")
	}).Methods("GET")

	router.HandleFunc("/{id}/room", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			tsLogger.LogWarn("could not parse %d", id)
			return
		}
		err = serveClientConnection(hub.rooms[int(id)], w, r)
		if err != nil {
			tsLogger.LogErr("CONNECTION FAILED")
			_, _ = fmt.Fprintln(w, err)
			return
		}
		tsLogger.LogTrace("CONNECTION TO %q", r.RequestURI)
	})

	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			tsLogger.LogWarn("could not parse %d", id)
			return
		}
		tsLogger.LogTrace("CLOSING ROOM %d", id)
		hub.closeRoom(int(id))

		_, _ = fmt.Fprint(w, "Room ", id, " closing")
	}).Methods("DELETE")

	router.HandleFunc("/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			tsLogger.LogWarn("could not parse %d", id)
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

func serveClientConnection(room *Room, w http.ResponseWriter, r *http.Request) error {
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
