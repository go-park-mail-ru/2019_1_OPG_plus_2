package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

var baseUrl = "localhost:8001/api"

func TestCreateUser(t *testing.T) {

	newUser := models.SingUpData{
		Email:    "new_user@mail.ru",
		Password: "qwerty",
		Username: "new_user",
	}

	reqBody, _ := json.Marshal(newUser)
	url := baseUrl + "/user"
	req := httptest.NewRequest("POST", url, bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	ctx := req.Context()
	data := models.JwtData{
		Username: "test_user",
		Email:    "test_email",
		Id:       1,
	}
	ctx = context.WithValue(ctx, "isAuth", true)
	ctx = context.WithValue(ctx, "jwtData", data)

	CreateUser(w, req.WithContext(ctx))
	if w.Code != http.StatusOK {
		t.Errorf("Wrong StatusCode: got %d, expected %d", w.Code, http.StatusOK)
	}
}
