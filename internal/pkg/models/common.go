package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	FieldsError   = errors.New("incorrect fields")
	NotFound      = errors.New("not found")
	AlreadyExists = errors.New("already exists")
)

type InputModel interface {
	Check() []string
}

type OutputModel interface {
	Send(w http.ResponseWriter)
}

func Send(w http.ResponseWriter, statusCode int, outModel interface{}) {
	w.WriteHeader(statusCode)
	jsonModel, _ := json.Marshal(outModel)
	_, _ = fmt.Fprintln(w, string(jsonModel))
}
