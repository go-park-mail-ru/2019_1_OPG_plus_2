package gameservice

import (
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"fmt"
	"time"
)

type Hub struct {
	rooms  map[int]*Room
	closer chan int
}

func NewHub() *Hub {
	return &Hub{
		closer: make(chan int),
		rooms:  make(map[int]*Room),
	}
}

func (h *Hub) AttachRooms(rooms ...*Room) error {
	for _, room := range rooms {
		if h.rooms[room.id] != nil {
			tsLogger.Logger.LogErr("ROOM %d EXISTS", room.id)
			return fmt.Errorf("ROOM %d EXISTS", room.id)
		}
		h.rooms[room.id] = room
		go room.Run()
	}
	return nil
}

func (h *Hub) run() {
	ticker := time.NewTicker(time.Second * 5)
	for _, room := range h.rooms {
		go room.Run()
	}

	for range ticker.C {
		tsLogger.Logger.LogInfo("HUB INFO: %+v", h.rooms)
	}
}

func (h *Hub) closeRoom(id int) {
	h.closer <- int(id)
	delete(h.rooms, int(id))
}
