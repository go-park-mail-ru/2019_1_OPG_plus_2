package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var FieldsError = errors.New("incorrect fields")
var NotFound = errors.New("not found")

type InputModel interface {
	Check() []string
}

type OutputModel interface {
	Send(w http.ResponseWriter)
}

/* COMMON MODELS */

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

type AnswerMessageWithFields struct {
	AnswerMessage
	Data []string `json:"data" example:"[email, username]"`
}

func (message AnswerMessageWithFields) Send(w http.ResponseWriter) {
	w.WriteHeader(message.Status)
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
}

func SendMessageWithFields(w http.ResponseWriter, status int, message string, data []string) {
	AnswerMessageWithFields{
		AnswerMessage: AnswerMessage{
			Status:  status,
			Message: message,
		},
		Data: data,
	}.Send(w)
}
