package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

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
	userToFind, _ := mockedStorageAdapter.GetUser(1)
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

			expMessage: models.UserDataAnswer{
				Status:  105,
				Message: "user found",
				Data:    userToFind,
			},
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.UserDataAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\n%v\n%v", retMessage, tCase.expMessage)
		}

		//testLog(t, tCase)
	}
}

func TestGetUserId(t *testing.T) {
	userToFind, _ := mockedStorageAdapter.GetUser(3)
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

			expMessage: models.UserDataAnswer{
				Status:  105,
				Message: "user found",
				Data:    userToFind,
			},
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.UserDataAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
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

			expStatus:    404,
			inputMessage: nil,

			expMessage: models.UserNotFoundAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
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

			expMessage: models.GetIncorrectFieldsAnswer([]string{"id"}),
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.IncorrectFieldsAnswer
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

			expMessage: models.GetIncorrectFieldsAnswer([]string{"id"}),
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.IncorrectFieldsAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
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

			expMessage: models.UserUpdatedAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
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

			expMessage: models.NotSignedInAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
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

			expStatus:    400,
			inputMessage: []byte(`{"email": "qwerty@mail.com","user": "qwerty"}`),

			expMessage: models.GetIncorrectFieldsAnswer([]string{"username"}),
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

			expStatus:    400,
			inputMessage: []byte(`{"e-mail": "qwerty@mail.com","username": "qwerty"}`),

			expMessage: models.GetIncorrectFieldsAnswer([]string{"email"}),
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

			expStatus:    400,
			inputMessage: []byte(`{"email": "qwerty","user-name": "qwerty"}`),

			expMessage: models.GetIncorrectFieldsAnswer([]string{"email", "username"}),
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.IncorrectFieldsAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
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

			expMessage: models.IncorrectJsonAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
		}

		//testLog(t, tCase)
	}
}

/****************************
 *  DELETE_USER CONTROLLER  *
 ****************************/

func TestRemoveUserCorrect(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.RemoveUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "qwerty",
					Email:    "qwerty@mail.com",
				},
			},

			expStatus:    200,
			inputMessage: []byte(`{"password": "pass1"}`),

			expMessage: models.UserRemovedAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
		}

		_, err := mockedStorageAdapter.GetUser(tCase.params.jwt.Id)
		if err == nil {
			t.Errorf("Data did not actually delete")
		}
	}
}

//since this moment i would use another user as a client for testing
//because in previous test user 1 has been deleted successfully

func TestRemoveUserNoAuth(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.RemoveUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			expStatus:    401,
			inputMessage: []byte(`{"password": "pass2"}`),

			expMessage: models.NotSignedInAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
		}
	}
}

func TestRemoveUserInvalidJSON(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.RemoveUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       2,
					Username: "username2",
					Email:    "mail2",
				},
			},

			expStatus:    500,
			inputMessage: []byte(`{"password": "pass1"`), // no closing parentheses

			expMessage: models.IncorrectJsonAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
		}
	}
}

func TestRemoveUserInvalidData(t *testing.T) {
	tCases := []testCase{
		{
			handler: mockedUserHandlers.RemoveUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       2,
					Username: "username2",
					Email:    "mail2",
				},
			},

			inputMessage: []byte(`{"passw": "pass1"}`),

			expStatus:  400,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"password"}),
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot %d\n\tExpected %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.IncorrectFieldsAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong body\nGOT:\t%v\nEXP:\t%v", retMessage, tCase.expMessage)
		}

		_, err := mockedStorageAdapter.GetUser(tCase.params.jwt.Id)
		if err != nil {
			t.Errorf("Data did actually delete!")
		}
	}
}

/****************************
 *  CREATE_USER CONTROLLER  *
 ****************************/
