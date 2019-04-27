package chatservice

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (c *ChatService) ConnectToRoom(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		c.Log.LogWarn("could not parse %d", id)
		return
	}
	err = c.serveClientConnection(c.Hub.rooms[int(id)], w, r)
	if err != nil {
		c.Log.LogErr("CONNECTION FAILED")
		_, _ = fmt.Fprintln(w, err)
		return
	}
	c.Log.LogTrace("CONNECTION TO %q", r.RequestURI)
}

func (c *ChatService) CreateRoom(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		c.Log.LogWarn("could not parse %d", id)
		return
	}

	err = c.Hub.AttachRooms(newRoom(c.Hub, int(id)))
	if err != nil {
		_, _ = fmt.Fprint(w, err)
		return
	}
	_, _ = fmt.Fprint(w, "ChatRoom ", id, " created")
}

func (c *ChatService) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		c.Log.LogWarn("could not parse %d", id)
		return
	}
	c.Log.LogTrace("CLOSING ROOM %d", id)
	c.Hub.closeRoom(int(id))

	_, _ = fmt.Fprint(w, "ChatRoom ", id, " closing")
}
