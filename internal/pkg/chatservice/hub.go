package chatservice

import (
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"fmt"
	"time"
)

type Hub struct {
	rooms  map[int]*ChatRoom
	Log    *tsLogger.TSLogger
	closer chan int
}

func NewHub() *Hub {
	return &Hub{
		closer: make(chan int),
		rooms:  make(map[int]*ChatRoom),
	}
}

func (h *Hub) AttachRooms(rooms ...*ChatRoom) error {
	for _, room := range rooms {
		if h.rooms[room.id] != nil {
			h.Log.LogErr("CHAT: ROOM %d EXISTS", room.id)
			return fmt.Errorf("CHAT: ROOM %d EXISTS", room.id)
		}
		h.rooms[room.id] = room
		go room.Run()
		h.Log.LogTrace("CHAT: CREATING ROOM %d", room.id)
	}
	return nil
}

func (h *Hub) Run() {
	ticker := time.NewTicker(time.Second * 5)
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

		h.Log.LogInfo("CHAT: HUB INFO: conns: %d, rooms : %d", activeConns(), len(h.rooms))
	}
}

func (h *Hub) closeRoom(id int) {
	h.closer <- int(id)
	delete(h.rooms, int(id))
	h.Log.LogTrace("CHAT: HUB: Closing room %d", id)
}
