package models

import (
    "fmt"
    "time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
    stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
    return []byte(stamp), nil
}

type ChatMessage struct {
    Id       int64    `json:"id, string" example:"1"`
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
