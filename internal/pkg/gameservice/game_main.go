package gameservice

import (
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type Service struct {
	Hub      *Hub
	Log      *tsLogger.TSLogger
	upgrader websocket.Upgrader
}

func NewService(hub *Hub, log *tsLogger.TSLogger) *Service {
	return &Service{
		Hub: hub,
		Log: log,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// TODO: update users' score mechanics
// TODO: delimit game as separate service
func (s *Service) AddGameServicePaths(router *mux.Router) *mux.Router {
	router.HandleFunc("/{id}", s.CreateRoom).Methods("POST")
	router.HandleFunc("/{id}", s.GetRoom).Methods("GET")
	router.HandleFunc("/{id}", s.DeleteRoom).Methods("DELETE")
	router.HandleFunc("/{id}/room", s.ConnectionEndpoint)

	go s.Hub.run()
	return router
}

func (s *Service) serveClientConnection(room *Room, w http.ResponseWriter, r *http.Request) error {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Log.LogErr("CONNECTION UPGRADE ERROR: %s", err)
		return err
	}
	client := NewClient(room, conn)
	client.room.register <- client

	go client.writePump()
	go client.readPump()
	return nil
}

func (s *Service) GetRoom(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.Log.LogWarn("could not parse %d", id)
		return
	}
	if s.Hub.rooms[int(id)] == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "no such room with id ", id)
		return
	}
	//http.ServeFile(w, r, "home.html")
	_, _ = fmt.Fprint(w, "IDI NAHUY")
}

func (s *Service) ConnectionEndpoint(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.Log.LogWarn("could not parse %d", id)
		return
	}
	if s.Hub.rooms[int(id)] == nil {
		s.upgrader.Error(w, r, http.StatusNotFound, fmt.Errorf("no room with id %v", id))
		return
	}
	err = s.serveClientConnection(s.Hub.rooms[int(id)], w, r)
	if err != nil {
		s.Log.LogErr("CONNECTION FAILED")
		s.upgrader.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.Log.LogTrace("CONNECTION TO %q", r.RequestURI)

}

func (s *Service) CreateRoom(w http.ResponseWriter, r *http.Request) {
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
	_, _ = fmt.Fprint(w, "Room ", id, " created")
}

func (s *Service) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.Log.LogWarn("could not parse %d", id)
		return
	}
	s.Log.LogTrace("CLOSING ROOM %d", id)
	s.Hub.closeRoom(int(id))

	_, _ = fmt.Fprint(w, "Room ", id, " closing")
}
