package gameservice

import (
	"2019_1_OPG_plus_2/internal/pkg/monitoring"
	"fmt"
	"sync"
	"time"
)

// I have both mutex and chan because Hub is the most stressed component of all system
// The mutex+chan setup has shown the best RPS result, so I stopped here
type Hub struct {
	rooms    map[string]*Room
	attacher chan *Room
	closer   chan string
	service  *Service
	mutex    *sync.RWMutex
}

func NewHub(service *Service) *Hub {
	return &Hub{
		closer:   make(chan string, 1024),
		attacher: make(chan *Room, 1024),
		rooms:    make(map[string]*Room),
		service:  service,
		mutex:    &sync.RWMutex{},
	}
}

func (h *Hub) AttachRooms(rooms ...*Room) error {
	for _, room := range rooms {
		h.mutex.RLock()
		if h.rooms[room.id] != nil {
			h.mutex.RUnlock()
			h.service.Log.LogErr("ROOM %q EXISTS", room.id)
			return fmt.Errorf("ROOM %q EXISTS", room.id)
		}
		h.mutex.RUnlock()
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
			h.mutex.Unlock()
			go room.Run()
		case <-ticker.C:
			h.service.Log.LogInfo("HUB INFO: conns: %d, rooms : %d", activeConns(), len(h.rooms))
			monitoring.ActiveConns.Set(float64(activeConns()))
			monitoring.ActiveRooms.Set(float64(len(h.rooms)))
		}
	}
}

func (h *Hub) closeRoom(id string) {
	h.closer <- id
	h.mutex.Lock()
	delete(h.rooms, id)
	h.mutex.Unlock()
}
