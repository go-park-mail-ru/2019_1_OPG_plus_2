package gameservice

import (
	"2019_1_OPG_plus_2/internal/pkg/monitoring"
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
			tsLogger.LogErr("ROOM %d EXISTS", room.id)
			return fmt.Errorf("ROOM %d EXISTS", room.id)
		}
		h.rooms[room.id] = room
		go room.Run()
		tsLogger.LogTrace("CREATING ROOM %d", room.id)
	}
	return nil
}

func (h *Hub) run() {
	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()
	for _, room := range h.rooms {
		go room.Run()
	}

	activeConns := func() int {
		cnt := 0
		for _, room := range h.rooms {
			for _, v := range room.clients {
				if v {
					cnt += 1
				}
			}
		}
		return cnt
	}
	for range ticker.C {
		tsLogger.LogInfo("HUB INFO: conns: %d, rooms : %d", activeConns(), len(h.rooms))
		monitoring.ActiveConns.Set(float64(activeConns()))
		monitoring.ActiveRooms.Set(float64(len(h.rooms)))
	}
}

func (h *Hub) closeRoom(id int) {
	h.closer <- int(id)
	delete(h.rooms, int(id))
}
