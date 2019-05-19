package models

type RoomData struct {
	Id         string   `json:"id, string"`
	PlayersNum int      `json:"players_num"`
	Players    []string `json:"players"`
}

type RoomsOnlineMessage struct {
	RoomsOnline []RoomData `json:"rooms_online"`
}

type RoomDeletedMessage struct {
	MessageAnswer
	RoomId string `json:"room_id"`
}

func NewRoomDeletedMessage(id string) RoomDeletedMessage {
	return RoomDeletedMessage{
		MessageAnswer: MessageAnswer{
			Status:  400,
			Message: "Room " + id + " deleted successfully",
		},
		RoomId: id,
	}
}
