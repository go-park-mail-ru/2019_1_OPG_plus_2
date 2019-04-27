package chatservice

import (
	"fmt"
	"time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
	return []byte(stamp), nil
}

type GenericMessage struct {
	Username string `json:"username, string"`
	MType    string `json:"type, string"`
}

type ChatMessage struct {
	GenericMessage
	Content  string `json:"content, string"`
	RandomId int64  `json:"random_id, number"`
}

type EventData struct {
	EventType string      `json:"event_type, string"`
	EventData interface{} `json:"event_data"`
}

type BroadcastEventMessage struct {
	GenericMessage
	Data EventData `json:"data"`
}

func NewBroadcastEventMessage(eType string, eData interface{}) *BroadcastEventMessage {
	return &BroadcastEventMessage{
		GenericMessage: GenericMessage{
			Username: "SERVICE",
			MType:    "event",
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

type BroadcastMessage struct {
	Id       int      `json:"id"`
	Avatar   string   `json:"avatar, string"`
	Username string   `json:"username, string"`
	Content  string   `json:"content, string"`
	Type     string   `json:"type, string"`
	Date     JSONTime `json:"date"`
}
