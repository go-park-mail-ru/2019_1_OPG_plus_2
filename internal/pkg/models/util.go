package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AnswerMessage struct {
	Status  int    `json:"status, int" example:"200"`
	Message string `json:"message, string" example:"Query processed successfully"`
}

func (message AnswerMessage) Send(w http.ResponseWriter) {
	w.WriteHeader(message.Status)
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
}

func SendMessage(w http.ResponseWriter, status int, message string) {
	AnswerMessage{
		Status:  status,
		Message: message,
	}.Send(w)
}

type AnswerMessageWithData struct {
	Status  int         `json:"status, int" example:"200"`
	Message string      `json:"message, string" example:"Query processed successfully"`
	Data    interface{} `json:"data" example:"[{id: 1}, {id: 2}]"`
}

func (message AnswerMessageWithData) Send(w http.ResponseWriter) {
	w.WriteHeader(message.Status)
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
}

func SendMessageWithData(w http.ResponseWriter, status int, message string, data interface{}) {
	AnswerMessageWithData{
		Status:  status,
		Message: message,
		Data:    data,
	}.Send(w)
}
