package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
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
		return models.UserData{}, fmt.Errorf("user not found")
	}
	return *storage.ProfileData[id], nil
}

func (storage *mockStorageAdapter) UpdateUser(id int64, updateData models.UpdateUserData) (jwtData models.JwtData, err error) {
	incorrectFields := updateData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, fmt.Errorf("incorrect: " + strings.Join(incorrectFields, ", "))
	}
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
	inputMessage json.RawMessage
	expStatus    int
	expMessage   interface{}
}

func testInitial(tCase testCase) (*httptest.ResponseRecorder, *http.Request) {
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

func testLog(t *testing.T, tCase testCase) {
	if !t.Failed() {
		t.Logf("\nPASSED TEST:\n"+
			"\tURL:\t\t%v\n"+
			"\tAUTH:\t\t%v\n"+
			"\tMETHOD:\t\t%v\n"+
			"\tJWT:\t\t%v\n"+
			"\tMUXVARS:\t%v\n"+
			"\tBODY:\t\t%v\n"+
			"\n"+
			"\tEXP_STATUS:\t\t%v\n"+
			"\tEXP_BODY:\t\t%v\n",

			tCase.params.url,
			tCase.params.isAuth,
			tCase.params.method,
			tCase.params.jwt,
			tCase.params.muxVars,
			tCase.inputMessage,

			tCase.expStatus,
			tCase.expMessage,
		)
	}
}

/*************************
 *  GET_USER CONTROLLER  *
 *************************/
func TestGetUserSelf(t *testing.T) {
	tCases := []testCase{
		{
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
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.UserDataAnswerMessage
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
		}

		//testLog(t, tCase)
	}
}

func TestGetUserId(t *testing.T) {
	retData, _ := mockedStorageAdapter.GetUser(3)
	tCases := []testCase{
		{
			handler: mockedUserHandlers.GetUser,

			params: testParams{
				muxVars: map[string]string{"id": "3"},
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
				Data: retData,
				AnswerMessage: models.AnswerMessage{
					Status:  200,
					Message: "user found",
				},
			},
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.UserDataAnswerMessage
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
		}

		//testLog(t, tCase)
	}
}

func TestGetUserIdNotExists(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.GetUser,

			params: testParams{
				muxVars: map[string]string{"id": "1278"},
				method:  "GET",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1",
				},
			},

			expStatus:    500,
			inputMessage: nil,

			expMessage: models.AnswerMessage{
				Status:  500,
				Message: "user not found",
			},
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.AnswerMessage
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
		}

		//testLog(t, tCase)
	}
}

func TestGetUserNoAuth(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.GetUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			expStatus:    400,
			inputMessage: nil,

			expMessage: models.AnswerMessage{
				Status:  400,
				Message: "no id in query",
			},
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.AnswerMessage
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
		}

		//testLog(t, tCase)
	}
}

func TestGetUserIdNotNumber(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.GetUser,

			params: testParams{
				muxVars: map[string]string{"id": "qwerty"},
				method:  "GET",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1",
				},
			},

			expStatus:    400,
			inputMessage: nil,

			expMessage: models.AnswerMessage{
				Status:  400,
				Message: "incorrect id in query",
			},
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.AnswerMessage
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
		}

		//testLog(t, tCase)
	}
}

/****************************
 *  UPDATE_USER CONTROLLER  *
 ****************************/
func TestUpdateUserCorrect(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.UpdateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1",
				},
			},

			expStatus:    200,
			inputMessage: []byte(`{"email": "qwerty@mail.com","username": "qwerty"}`),

			expMessage: models.AnswerMessage{
				Status:  200,
				Message: "user updated",
			},
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.AnswerMessage
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
		}

		newUserData, err := mockedStorageAdapter.GetUser(tCase.params.jwt.Id)
		if err != nil {
			t.Errorf("Test failed while getting user from storage: %e", err)
		}

		expUserData := models.UserData{
			Id:       1,
			Username: "qwerty",
			Email:    "qwerty@mail.com",
			Avatar:   newUserData.Avatar,
			Score:    newUserData.Score,
			Games:    newUserData.Games,
			Win:      newUserData.Win,
			Lose:     newUserData.Win,
		}

		if !reflect.DeepEqual(newUserData, expUserData) {
			t.Errorf("Data did not actually update")
		}

		//testLog(t, tCase)
	}
}

func TestUpdateUserNoAuth(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.UpdateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			expStatus:    401,
			inputMessage: []byte(`{"email": "qwerty","username": "qwerty"}`),

			expMessage: models.AnswerMessage{
				Status:  401,
				Message: "not signed in",
			},
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.AnswerMessage
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
		}

		//testLog(t, tCase)
	}
}

func TestUpdateUserInvalidField(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.UpdateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "qwerty",
					Email:    "qwerty@mail.com",
				},
			},

			expStatus:    500,
			inputMessage: []byte(`{"email": "qwerty@mail.com","user": "qwerty"}`),

			expMessage: models.AnswerMessage{
				Status:  500,
				Message: "incorrect: username",
			},
		},
		{
			handler: mockedUserHandlers.UpdateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "qwerty",
					Email:    "qwerty@mail.com",
				},
			},

			expStatus:    500,
			inputMessage: []byte(`{"e-mail": "qwerty@mail.com","username": "qwerty"}`),

			expMessage: models.AnswerMessage{
				Status:  500,
				Message: "incorrect: email",
			},
		},
		{
			handler: mockedUserHandlers.UpdateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "qwerty",
					Email:    "qwerty@mail.com",
				},
			},

			expStatus:    500,
			inputMessage: []byte(`{"email": "qwerty","user-name": "qwerty"}`),

			expMessage: models.AnswerMessage{
				Status:  500,
				Message: "incorrect: email, username",
			},
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.AnswerMessage
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
		}

		//testLog(t, tCase)
	}
}

func TestUpdateUserInvalidJSON(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.UpdateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "qwerty",
					Email:    "qwerty@mail.com",
				},
			},

			expStatus:    500,
			inputMessage: []byte(`{"email": "qwerty@mail.com","user": "qwerty"`), //no closing parentheses in JSON

			expMessage: models.AnswerMessage{
				Status:  500,
				Message: "incorrect JSON",
			},
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.AnswerMessage
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
		}

		//testLog(t, tCase)
	}
}
