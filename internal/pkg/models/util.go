package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SuccessOrErrorMessage struct {
	Status  int    `json:"status, int" example:"200"`
	Message string `json:"message, string" example:"Query processed successfully"`
}

func (message SuccessOrErrorMessage) Send(w http.ResponseWriter) {
	w.WriteHeader(message.Status)
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
}

func SendMessage(w http.ResponseWriter, status int, text string) {
	SuccessOrErrorMessage{
		Status:  status,
		Message: text,
	}.Send(w)
}
