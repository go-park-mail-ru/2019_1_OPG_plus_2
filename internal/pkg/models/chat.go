package models

import (
	"fmt"
	"time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))
	return []byte(stamp), nil
}

type ChatMessage struct {
	Id       int64    `json:"id, number" example:"1"`
	RandomId int64    `json:"random_id, number " example:"1"`
	Username string   `json:"username, string"`
	Avatar   string   `json:"avatar, string"`
	Content  string   `json:"content, string"`
	Type     string   `json:"type, string"`
	Date     JSONTime `json:"date"`
}

type ChatMessageList struct {
	Messages []ChatMessage `json:"messages"`
	Count    uint64        `json:"count" example:"123"`
}

type ChatMessageListAnswer struct {
	Status  int             `json:"status, int" example:"106"`
	Message string          `json:"message, string" example:"users found"`
	Data    ChatMessageList `json:"data"`
}

func GetMessageListAnswer(data ChatMessageList) *ChatMessageListAnswer {
	return &ChatMessageListAnswer{
		Status:  110,
		Message: "messages found",
		Data:    data,
	}
}
