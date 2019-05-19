package models

type RoomData struct {
	Id         string   `json:"id, string"`
	PlayersNum int      `json:"players_num"`
	Players    []string `json:"players"`
}

type RoomsOnlineMessage struct {
	RoomsOnline []RoomData `json:"rooms_online"`
}
