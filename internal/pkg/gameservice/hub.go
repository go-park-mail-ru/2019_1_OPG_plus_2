package gameservice

import (
	"fmt"
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
			return fmt.Errorf("room exists")
		}
		h.rooms[room.id] = room
		go room.Run()
	}
	return nil
}

func (h *Hub) run() {
	for _, room := range h.rooms {
		go room.Run()
	}
}

func (h *Hub) closeRoom(id int) {
	h.closer <- int(id)
	delete(h.rooms, int(id))
}
