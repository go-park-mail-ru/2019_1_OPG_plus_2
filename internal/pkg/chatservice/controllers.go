package chatservice

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *ChatService) ConnectToRoom(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.Log.LogWarn("could not parse %d", id)
		return
	}
	err = s.serveClientConnection(s.Hub.rooms[int(id)], w, r)
	if err != nil {
		s.Log.LogErr("CONNECTION FAILED")
		_, _ = fmt.Fprintln(w, err)
		return
	}
	s.Log.LogTrace("CONNECTION TO %q", r.RequestURI)
}

func (s *ChatService) CreateRoom(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.Log.LogWarn("could not parse %d", id)
		return
	}

	err = s.Hub.AttachRooms(newRoom(s.Hub, int(id)))
	if err != nil {
		_, _ = fmt.Fprint(w, err)
		return
	}
	_, _ = fmt.Fprint(w, "ChatRoom ", id, " created")
}

func (s *ChatService) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.Log.LogWarn("could not parse %d", id)
		return
	}
	s.Log.LogTrace("CLOSING ROOM %d", id)
	s.Hub.closeRoom(int(id))

	_, _ = fmt.Fprint(w, "ChatRoom ", id, " closing")
}
