package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
)

var baseUrl = "localhost:8002/api"

type authData struct {
	Id       int64
	Username string
	Email    string
	Password string
}

type mockStorageAdapter struct {
	ProfileData map[int64]*models.UserData
	AuthData    map[int64]*authData
}

func newMockStorageAdapter() (storage *mockStorageAdapter) {
	storage = new(mockStorageAdapter)
	storage.ProfileData = make(map[int64]*models.UserData)
	storage.AuthData = make(map[int64]*authData)
	storage.ProfileData[1] = &models.UserData{
		Id:       1,
		Email:    "mail1",
		Username: "username1",
		Score:    1000,
		Avatar:   "avatar1",
		Games:    100,
		Lose:     50,
		Win:      50,
	}
	storage.AuthData[1] = &authData{
		Email:    "mail1",
		Username: "username1",
		Id:       1,
		Password: "pass1",
	}
	storage.ProfileData[2] = &models.UserData{
		Id:       2,
		Email:    "mail2",
		Username: "username2",
		Score:    2000,
		Avatar:   "avatar2",
		Games:    200,
		Lose:     100,
		Win:      100,
	}
	storage.AuthData[2] = &authData{
		Email:    "mail2",
		Username: "username2",
		Id:       2,
		Password: "pass2",
	}
	storage.ProfileData[3] = &models.UserData{
		Id:       3,
		Email:    "mail3",
		Username: "username3",
		Score:    3000,
		Avatar:   "avatar3",
		Games:    300,
		Lose:     150,
		Win:      150,
	}
	storage.AuthData[3] = &authData{
		Email:    "mail3",
		Username: "username3",
		Id:       3,
		Password: "pass3",
	}
	return storage
}

func (storage *mockStorageAdapter) CreateUser(signUpData models.SingUpData) (models.JwtData, error, []string) {
	incorrectFields := signUpData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}

	newUser := models.UserData{
		Id:       int64(storage.ProfileData[int64(len(storage.ProfileData))].Id + 1),
		Username: signUpData.Username,
		Email:    signUpData.Email,
	}

	for _, profile := range storage.ProfileData {
		if profile.Username == newUser.Username || profile.Email == newUser.Email {
			return models.JwtData{}, models.AlreadyExists, nil
		}
	}

	storage.ProfileData[newUser.Id] = &newUser
	storage.AuthData[newUser.Id] = &authData{
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
	return newJwtData, nil, nil
}

func (storage *mockStorageAdapter) GetUser(id int64) (userData models.UserData, err error) {
	if storage.ProfileData[id] == nil {
		return models.UserData{}, models.NotFound
	}
	return *storage.ProfileData[id], nil
}

func (storage *mockStorageAdapter) UpdateUser(id int64, updateData models.UpdateUserData) (models.JwtData, error, []string) {
	incorrectFields := updateData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}
	if storage.ProfileData[id] == nil {
		return models.JwtData{}, models.NotFound, nil
	}

	user := storage.ProfileData[id]
	user.Username = updateData.Username
	user.Email = updateData.Email

	newJwtData := models.JwtData{
		Id:       id,
		Username: user.Username,
		Email:    user.Email,
	}
	return newJwtData, nil, nil
}

func (storage *mockStorageAdapter) RemoveUser(id int64, removeData models.RemoveUserData) (error, []string) {
	incorrectFields := removeData.Check()
	if len(incorrectFields) > 0 {
		return models.FieldsError, incorrectFields
	}

	if removeData.Password != storage.AuthData[id].Password {
		return models.FieldsError, append(incorrectFields, "password")
	}

	delete(storage.ProfileData, id)
	delete(storage.AuthData, id)
	return nil, nil
}

type TestParams struct {
	isAuth  bool
	muxVars map[string]string
	jwt     models.JwtData
	method  string
	url     string
}

type TestCase struct {
	handler      http.HandlerFunc
	params       TestParams
	inputMessage json.RawMessage
	expStatus    int
	expMessage   interface{}
}

func testInitial(tCase TestCase) (*httptest.ResponseRecorder, *http.Request) {
	testParams := tCase.params
	url := baseUrl + testParams.url
	req := httptest.NewRequest(testParams.method, url, bytes.NewReader(tCase.inputMessage))
	w := httptest.NewRecorder()
	ctx := req.Context()

	data := testParams.jwt
	ctx = context.WithValue(ctx, "isAuth", testParams.isAuth)
	ctx = context.WithValue(ctx, "jwtData", data)

	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, tCase.params.muxVars)

	return w, req
}
