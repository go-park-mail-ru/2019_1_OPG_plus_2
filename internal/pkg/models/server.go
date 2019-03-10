package models

type IndexMessage struct {
	Message string `json:"message, string"`
}

func NewIndexMessage(message string) *IndexMessage {
	return &IndexMessage{Message: message}
}
