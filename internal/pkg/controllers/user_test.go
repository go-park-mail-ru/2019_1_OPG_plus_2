package controllers

import (
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

var baseUrl = "localhost:8001/api"

type authData struct {
	Id       int64
	Username string
	Email    string
	Password string
}

type MockStorageAdapter struct {
	ProfileData map[int64]models.UserData
	AuthData    map[int64]authData
}

func NewMockStorageAdapter() *MockStorageAdapter {
	data := make(map[int64]models.UserData)
	data[1] = models.UserData{
		Id:       1,
		Email:    "mail1",
		Username: "username1",
		Score:    1000,
		Avatar:   "avatar1",
		Games:    100,
		Lose:     50,
		Win:      50,
	}
	data[2] = models.UserData{
		Id:       2,
		Email:    "mail2",
		Username: "username2",
		Score:    2000,
		Avatar:   "avatar2",
		Games:    200,
		Lose:     100,
		Win:      100,
	}
	data[3] = models.UserData{
		Id:       3,
		Email:    "mail3",
		Username: "username3",
		Score:    3000,
		Avatar:   "avatar3",
		Games:    300,
		Lose:     150,
		Win:      150,
	}

	return &MockStorageAdapter{ProfileData: data}
}

func (storage *MockStorageAdapter) CreateUser(signUpData models.SingUpData) (jwtData models.JwtData, err error) {
	newUser := models.UserData{
		Id:       int64(len(storage.ProfileData)),
		Username: signUpData.Username,
		Email:    signUpData.Email,
		Avatar:   signUpData.Avatar,
	}

	storage.ProfileData[newUser.Id] = newUser
	storage.AuthData[newUser.Id] = authData{
		Id:       newUser.Id,
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: signUpData.Password,
	}

	newJwtData := models.JwtData{
		Id:       newUser.Id,
		Email:    newUser.Email,
		Username: newUser.Username,
	}
	return newJwtData, nil
}

func (storage *MockStorageAdapter) GetUser(id int64) (userData models.UserData, err error) {
	return storage.ProfileData[id], nil
}

func (storage *MockStorageAdapter) UpdateUser(id int64, updateData models.UpdateUserData) (jwtData models.JwtData, err error) {
	user := storage.ProfileData[id]
	user.Username = updateData.Username
	user.Email = updateData.Email

	newJwtData := models.JwtData{
		Id:       id,
		Username: user.Username,
		Email:    user.Email,
	}
	return newJwtData, nil
}

func (*MockStorageAdapter) RemoveUser(id int64, removeData models.RemoveUserData) error {
	return nil
}

var mockedStorageAdapter = NewMockStorageAdapter()
var mockedUserHandlers = NewUserHandlers(mockedStorageAdapter)

func TestGetUserSelf(t *testing.T) {
	url := baseUrl + "/user"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	ctx := req.Context()

	data := models.JwtData{
		Username: "username1",
		Email:    "mail1",
		Id:       1,
	}
	ctx = context.WithValue(ctx, "isAuth", true)
	ctx = context.WithValue(ctx, "jwtData", data)

	mockedUserHandlers.GetUser(w, req.WithContext(ctx))
	if w.Code != http.StatusOK {
		t.Errorf("Wrong StatusCode: got %d, expected %d\n Body: %v", w.Code, http.StatusOK, w.Body)
	}
}

func TestGetUserId(t *testing.T) {
	url := baseUrl + "/user/" + "3"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	ctx := req.Context()

	data := models.JwtData{
		Username: "username1",
		Email:    "mail1",
		Id:       1,
	}
	ctx = context.WithValue(ctx, "isAuth", true)
	ctx = context.WithValue(ctx, "jwtData", data)
	req = mux.SetURLVars(req.WithContext(ctx), map[string]string{
		"id": "3",
	})

	mockedUserHandlers.GetUser(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong StatusCode: got %d, expected %d\n Body: %v", w.Code, http.StatusOK, w.Body)
	}
	var gotMessage models.AnswerMessageWithData
	err := json.NewDecoder(w.Body).Decode(&gotMessage)
	profileMap := gotMessage.Data
	profileJson, err := json.Marshal(profileMap)
	var profile models.UserData
	err = json.Unmarshal([]byte(profileJson), &profile)
	t.Log(profile)
	if err != nil {
		t.Errorf("Json parsing bad: %v", err)
	}
}
