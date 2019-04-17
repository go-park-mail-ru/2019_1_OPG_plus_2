package gameservice

type GenericMessage struct {
	User  string `json:"user, string"`
	MType string `json:"type, string"`
}

type GameMessage struct {
	GenericMessage
	Data struct {
		Coords []int `json:"coords"`
	} `json:"data"`
}

type ChatMessage struct {
	GenericMessage
	Data struct {
		Text string `json:"text, string"`
	} `json:"data"`
}

type EventData struct {
	EventType string      `json:"event_type, string"`
	EventData interface{} `json:"event_data"`
}

type BroadcastEventMessage struct {
	GenericMessage
	Data EventData `json:"data"`
}

type RegisterMessage GenericMessage

func NewBroadcastEventMessage(e_type string, e_data interface{}) *BroadcastEventMessage {
	return &BroadcastEventMessage{
		GenericMessage: GenericMessage{
			User:  "SERVICE",
			MType: "event",
		},
		Data: EventData{
			EventType: e_type,
			EventData: e_data,
		},
	}
}

func NewBroadcastErrorMessage(e string) *BroadcastEventMessage {
	return NewBroadcastEventMessage("error", e)
}
