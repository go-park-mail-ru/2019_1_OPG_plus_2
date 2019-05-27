package gameservice

import (
	"2019_1_OPG_plus_2/internal/pkg/middleware"
	"2019_1_OPG_plus_2/internal/pkg/models"
	"2019_1_OPG_plus_2/internal/pkg/randomgenerator"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
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
func (s *Service) AddGameServicePaths(router *mux.Router) *mux.Router {
	router.HandleFunc("/new_room", s.NewRoom).Methods("GET")
	router.HandleFunc("/rooms", s.ListRooms).Methods("GET")
	router.HandleFunc("/free_room", s.GetFreeRoom)
	router.HandleFunc("/{id}", s.CreateRoom).Methods("POST")
	router.HandleFunc("/{id}", s.GetRoom).Methods("GET")
	router.HandleFunc("/{id}", s.DeleteRoom).Methods("DELETE")
	router.HandleFunc("/{id}/room", s.ConnectionEndpoint)
	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.AuthMiddleware)
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

// RoomInfo
func (s *Service) GetRoom(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		s.Log.LogWarn("could not parse %q", id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	room := s.Hub.rooms[id]
	if room == nil {
		models.Send(w, http.StatusNotFound, models.NotFound)
		return
	}

	roomData := models.RoomData{
		Id:         room.id,
		PlayersNum: room.currentPlayersNum,
		Players:    room.gameModel.players,
	}

	models.Send(w, http.StatusOK, roomData)

}

func (s *Service) ConnectionEndpoint(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	if id == "" {
		s.Log.LogWarn("could not parse %q", id)
		return
	}
	if s.Hub.rooms[id] == nil {
		s.upgrader.Error(w, r, http.StatusNotFound, fmt.Errorf("no room with id %v", id))
		return
	}
	err := s.serveClientConnection(s.Hub.rooms[id], w, r)
	if err != nil {
		s.Log.LogErr("CONNECTION FAILED")
		s.upgrader.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.Log.LogTrace("CONNECTION TO %q", r.RequestURI)

}

func (s *Service) CreateRoom(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		s.Log.LogWarn("could not parse %q", id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	room := newRoom(s.Hub, id)
	err := s.Hub.AttachRooms(room)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	var roomData = models.RoomData{
		Id:         room.id,
		PlayersNum: room.currentPlayersNum,
		Players:    room.gameModel.players,
	}
	models.Send(w, http.StatusOK, roomData)
}

func (s *Service) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		s.Log.LogWarn("could not parse %q", id)
		return
	}

	if s.Hub.rooms[id] == nil {
		models.Send(w, http.StatusNotFound, models.NotFound)
	}
	s.Log.LogTrace("CLOSING ROOM %q", id)
	s.Hub.closeRoom(id)

	models.Send(w, http.StatusOK, models.NewRoomDeletedMessage(id))
}

func (s *Service) ListRooms(w http.ResponseWriter, r *http.Request) {
	roomsOnline := models.RoomsOnlineMessage{}
	for k, v := range s.Hub.rooms {
		room := models.RoomData{
			Id:         k,
			PlayersNum: v.currentPlayersNum,
			Players:    v.gameModel.players,
		}
		roomsOnline.RoomsOnline = append(roomsOnline.RoomsOnline, room)
	}

	models.Send(w, http.StatusOK, roomsOnline)
}

func (s *Service) GetFreeRoom(w http.ResponseWriter, r *http.Request) {
	var freeRoom string
	found := false
	for k, v := range s.Hub.rooms {
		if v.currentPlayersNum < v.maxPlayersNum {
			found = true
			freeRoom = k
		}
	}

	var room *Room
	if !found {
		freeRoomId, err := randomgenerator.RandomString(6)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}
		room = newRoom(s.Hub, freeRoomId)
		err = s.Hub.AttachRooms(room)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}
	} else {
		room = s.Hub.rooms[freeRoom]
	}

	var roomData = models.RoomData{
		Id:         room.id,
		PlayersNum: room.currentPlayersNum,
		Players:    room.gameModel.players,
	}
	models.Send(w, http.StatusOK, roomData)
}
func (s *Service) NewRoom(w http.ResponseWriter, r *http.Request) {
	freeRoom, err := randomgenerator.RandomString(10)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
	}
	room := newRoom(s.Hub, freeRoom)
	err = s.Hub.AttachRooms(room)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
	}
	var roomData = models.RoomData{
		Id:         room.id,
		PlayersNum: room.currentPlayersNum,
		Players:    room.gameModel.players,
	}
	models.Send(w, http.StatusOK, roomData)
}
