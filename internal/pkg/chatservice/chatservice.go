package chatservice

import "2019_1_OPG_plus_2/internal/pkg/tsLogger"

type ChatService struct {
	Hub *Hub
	Log *tsLogger.TSLogger
}

func NewChatService(hub *Hub, log *tsLogger.TSLogger) *ChatService {
	s := &ChatService{Hub: hub, Log: log}
	err := s.Hub.AttachRooms(newRoom(hub, 0))
	if err != nil {
		tsLogger.LogErr("ROOM ATTACHMENT ERROR: %v", s.Hub.rooms)
		panic("WTF")
	}
	s.Log.LogTrace("INITIAL ROOM CREATED")
	s.Hub.Log = log
	return s
}
