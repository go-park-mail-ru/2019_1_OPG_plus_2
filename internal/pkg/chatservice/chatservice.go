package chatservice

import "2019_1_OPG_plus_2/internal/pkg/tsLogger"

type ChatService struct {
	Hub *Hub
	Log *tsLogger.TSLogger
}

func NewChatService(hub *Hub, log *tsLogger.TSLogger) *ChatService {
	s := &ChatService{Hub: hub, Log: log}
	s.Hub.Log = log
	err := s.Hub.AttachRooms(newRoom(hub, 0))
	s.Log.Run()
	if err != nil {
		s.Log.LogErr("CHAT: ROOM ATTACHMENT ERROR: %v", s.Hub.rooms)
		panic("WTF")
	}
	s.Log.LogTrace("CHAT: INITIAL ROOM 0 CREATED")
	return s
}
