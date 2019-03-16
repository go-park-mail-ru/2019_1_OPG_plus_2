package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var baseUrl = "localhost:8001/api"

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

func (storage *mockStorageAdapter) CreateUser(signUpData models.SingUpData) (jwtData models.JwtData, err error) {
	newUser := models.UserData{
		Id:       int64(len(storage.ProfileData)),
		Username: signUpData.Username,
		Email:    signUpData.Email,
		Avatar:   signUpData.Avatar,
	}

	for _, profile := range storage.ProfileData {
		if profile.Username == newUser.Username || profile.Email == newUser.Email {
			return models.JwtData{}, fmt.Errorf("USERNAME OR EMAIL IS ALREADY USED")
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
	return newJwtData, nil
}

func (storage *mockStorageAdapter) GetUser(id int64) (userData models.UserData, err error) {
	if storage.ProfileData[id] == nil {
		return models.UserData{}, fmt.Errorf("NO USER IN STORAGE")
	}
	return *storage.ProfileData[id], nil
}

func (storage *mockStorageAdapter) UpdateUser(id int64, updateData models.UpdateUserData) (jwtData models.JwtData, err error) {
	if storage.ProfileData[id] == nil {
		return models.JwtData{}, fmt.Errorf("NO USER IN STORAGE")
	}

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

func (storage *mockStorageAdapter) RemoveUser(id int64, removeData models.RemoveUserData) error {
	if removeData.Password != storage.AuthData[id].Password {
		return fmt.Errorf("PASSWORD DOES'T MATCH WITH ONE IN STORAGE")
	}

	delete(storage.ProfileData, id)
	delete(storage.AuthData, id)
	return nil
}

var mockedStorageAdapter = newMockStorageAdapter()
var mockedUserHandlers = NewUserHandlers(mockedStorageAdapter)

type testParams struct {
	isAuth  bool
	muxVars map[string]string
	jwt     models.JwtData
	method  string
	url     string
}

type testCase struct {
	handler      http.HandlerFunc
	params       testParams
	expStatus    int
	inputMessage interface{}
	expMessage   interface{}
}

var tCase = testCase{
	handler: mockedUserHandlers.GetUser,

	params: testParams{
		muxVars: map[string]string{},
		method:  "GET",
		isAuth:  true,
		url:     "/user",
		jwt: models.JwtData{
			Id:       1,
			Username: "username1",
			Email:    "mail1",
		},
	},

	expStatus:    200,
	inputMessage: nil,

	expMessage: models.UserDataAnswerMessage{
		Data: models.UserData{
			Id:       1,
			Email:    "mail1",
			Username: "username1",
			Score:    1000,
			Avatar:   "avatar1",
			Games:    100,
			Lose:     50,
			Win:      50,
		},
		AnswerMessage: models.AnswerMessage{
			Status:  200,
			Message: "user found",
		},
	},
}

func TestGetUserController(t *testing.T) {

	testParams := tCase.params
	url := baseUrl + testParams.url
	req := httptest.NewRequest(testParams.method, url, nil)
	w := httptest.NewRecorder()
	ctx := req.Context()

	data := testParams.jwt
	ctx = context.WithValue(ctx, "isAuth", testParams.isAuth)
	ctx = context.WithValue(ctx, "jwtData", data)

	tCase.handler(w, req.WithContext(ctx))

	if w.Code != tCase.expStatus {
		t.Errorf("Wrong StatusCode:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
	}
	var retMessage models.UserDataAnswerMessage
	_ = json.NewDecoder(w.Body).Decode(&retMessage)
	if !reflect.DeepEqual(retMessage, tCase.expMessage) {
		t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
	}

	if !t.Failed() {
		t.Logf("\nPASSED TEST:\n"+
			"\tURL:\t%v\n"+
			"\tAUTH:\t%v\n"+
			"\tMETHOD:\t%v\n"+
			"\tJWT:\t%v\n"+
			"\tMUXVARS:\t%v\n"+
			"\tBODY:\t%v\n"+
			"\n"+
			"\tEXP_STATUS:\t%v\n"+
			"\tEXP_BODY:\t%v\n",

			testParams.url,
			testParams.isAuth,
			testParams.method,
			testParams.jwt,
			testParams.muxVars,
			tCase.inputMessage,

			tCase.expStatus,
			tCase.expMessage,
		)
	}
}
