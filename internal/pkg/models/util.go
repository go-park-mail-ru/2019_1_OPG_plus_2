package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SuccessOrErrorMessage struct {
	Status  int    `json:"status, int"`
	Message string `json:"message, string"`
}

func (message *SuccessOrErrorMessage) Send(w http.ResponseWriter, status int, text string) {
	w.WriteHeader(status)
	message.Status = status
	message.Message = text
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
}
