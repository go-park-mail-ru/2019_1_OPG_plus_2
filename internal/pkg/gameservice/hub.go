package gameservice

import (
	"2019_1_OPG_plus_2/internal/pkg/monitoring"
	"fmt"
	"sync"
	"time"
)

type Hub struct {
	rooms    map[string]*Room
	attacher chan *Room
	closer   chan string
	service  *Service
	mutex    *sync.Mutex
}

func NewHub(service *Service) *Hub {
	return &Hub{
		closer:   make(chan string, 1024),
		attacher: make(chan *Room),
		rooms:    make(map[string]*Room),
		service:  service,
		mutex:    &sync.Mutex{},
	}
}

func (h *Hub) AttachRooms(rooms ...*Room) error {
	for _, room := range rooms {
		h.mutex.Lock()
		if h.rooms[room.id] != nil {
			h.mutex.Unlock()
			h.service.Log.LogErr("ROOM %q EXISTS", room.id)
			return fmt.Errorf("ROOM %q EXISTS", room.id)
		}
		h.mutex.Unlock()
		h.attacher <- room
	}
	return nil
}

func (h *Hub) Run() {
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

	for {
		select {
		case room := <-h.attacher:
			h.service.Log.LogTrace("CREATING ROOM %q", room.id)
			h.mutex.Lock()
			h.rooms[room.id] = room
			go room.Run()
			h.mutex.Unlock()
		case <-ticker.C:
			h.service.Log.LogInfo("HUB INFO: conns: %d, rooms : %d", activeConns(), len(h.rooms))
			monitoring.ActiveConns.Set(float64(activeConns()))
			monitoring.ActiveRooms.Set(float64(len(h.rooms)))
		}
	}
}

func (h *Hub) closeRoom(id string) {
	h.closer <- id
	delete(h.rooms, id)
}
