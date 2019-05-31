package gameservice

type GenericMessage struct {
    User  string `json:"user, string"`
    MType string `json:"type, string"`
}

type Point struct {
    X int `json:"x"`
    Y int `json:"y"`
}

type GameMessage struct {
    GenericMessage
    Data struct {
        Coords []Point `json:"coords"`
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

type RegisterMessage struct {
    GenericMessage
    Avatar string `json:"avatar, string"`
}

func NewBroadcastEventMessage(eType string, eData interface{}) *BroadcastEventMessage {
    return &BroadcastEventMessage{
        GenericMessage: GenericMessage{
            User:  "SERVICE",
            MType: "event",
        },
        Data: EventData{
            EventType: eType,
            EventData: eData,
        },
    }
}

func NewBroadcastErrorMessage(e string) *BroadcastEventMessage {
    return NewBroadcastEventMessage("error", e)
}
