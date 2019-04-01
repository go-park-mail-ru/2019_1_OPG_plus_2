package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

var mockedStorageAdapter = newMockStorageAdapter()
var mockedUserHandlers = NewUserHandlers(mockedStorageAdapter)

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

// func testLog(t *testing.T, tCase TestCase) {
// 	if !t.Failed() {
// 		t.Logf("\nPASSED TEST:\n"+
// 			"\tURL:\t\t%v\n"+
// 			"\tAUTH:\t\t%v\n"+
// 			"\tMETHOD:\t\t%v\n"+
// 			"\tJWT:\t\t%v\n"+
// 			"\tMUXVARS:\t%v\n"+
// 			"\tBODY:\t\t%v\n"+
// 			"\n"+
// 			"\tEXP_STATUS:\t\t%v\n"+
// 			"\tEXP_BODY:\t\t%v\n",
//
// 			tCase.params.url,
// 			tCase.params.isAuth,
// 			tCase.params.method,
// 			tCase.params.jwt,
// 			tCase.params.muxVars,
// 			tCase.inputMessage,
//
// 			tCase.expStatus,
// 			tCase.expMessage,
// 		)
// 	}
// }

/*************************
 *  GET_USER CONTROLLER  *
 *************************/
func TestGetUserSelf(t *testing.T) {
	userToFind, _ := mockedStorageAdapter.GetUser(1)
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.GetUser,

			params: TestParams{
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

			expMessage: models.GetUserDataAnswer(userToFind),
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}

		var retMessage models.UserDataAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}

		// testLog(t, tCase)
	}
}

func TestGetUserId(t *testing.T) {
	userToFind, _ := mockedStorageAdapter.GetUser(3)
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.GetUser,

			params: TestParams{
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

			expMessage: models.GetUserDataAnswer(userToFind),
		},
	}

	for _, tCase := range tCases {
		var retMessage models.UserDataAnswer
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}
		//
		//// testLog(t, tCase)
		//CheckStatusAndAnswer(t, w, retMessage, tCase)
	}
}

func TestGetUserIdNotExists(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.GetUser,

			params: TestParams{
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
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}

		// testLog(t, tCase)
	}
}

func TestGetUserNoAuth(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.GetUser,

			params: TestParams{
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
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.IncorrectFieldsAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}

		// testLog(t, tCase)
	}
}

func TestGetUserIdNotNumber(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.GetUser,

			params: TestParams{
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
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.IncorrectFieldsAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}

		// testLog(t, tCase)
	}
}

/****************************
 *  UPDATE_USER CONTROLLER  *
 ****************************/
func TestUpdateUserCorrect(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.UpdateUser,

			params: TestParams{
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
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
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

		// testLog(t, tCase)
	}
}

func TestUpdateUserNoAuth(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.UpdateUser,

			params: TestParams{
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
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}

		// testLog(t, tCase)
	}
}

func TestUpdateUserInvalidField(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.UpdateUser,

			params: TestParams{
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

			params: TestParams{
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

			params: TestParams{
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
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.IncorrectFieldsAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}

		// testLog(t, tCase)
	}
}

func TestUpdateUserInvalidJSON(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.UpdateUser,

			params: TestParams{
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
			inputMessage: []byte(`{"email": "qwerty@mail.com","user": "qwerty"`), // no closing parentheses in JSON

			expMessage: models.IncorrectJsonAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}

		// testLog(t, tCase)
	}
}

/****************************
 *  DELETE_USER CONTROLLER  *
 ****************************/

func TestRemoveUserCorrect(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.RemoveUser,

			params: TestParams{
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
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}

		_, err := mockedStorageAdapter.GetUser(tCase.params.jwt.Id)
		if err == nil {
			t.Errorf("Data did not actually delete")
		}
	}
}

// since this moment i would use another user as a client for testing
// because in previous test user 1 has been deleted successfully
func TestRemoveUserNoAuth(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.RemoveUser,

			params: TestParams{
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
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}
	}
}

func TestRemoveUserInvalidJSON(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.RemoveUser,

			params: TestParams{
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
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}
	}
}

func TestRemoveUserInvalidData(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.RemoveUser,

			params: TestParams{
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
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.IncorrectFieldsAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
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

func TestCreateUserCorrect(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.CreateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			expStatus:    200,
			inputMessage: []byte(`{"email": "new_user@mail.com","username": "new_user", "password": "pass"}`),

			expMessage: models.SignedUpAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}

		//checking if user is actually created

		//getting data field from cookie set to user
		rawJwtData := strings.Split(w.Result().Cookies()[0].Value, ".")[1]
		parsedCookie, _ := base64.StdEncoding.DecodeString(rawJwtData)

		//parsing data of user stored in cache into struct
		var storedUserData models.JwtData
		_ = json.Unmarshal(parsedCookie, &storedUserData)

		var inputData models.SingUpData
		_ = json.Unmarshal(tCase.inputMessage, &inputData)

		if storedUserData.Email != inputData.Email || storedUserData.Username != inputData.Username {
			t.Error("User saved incorrectly")
		}

		// testLog(t, tCase)
	}
}

func TestCreateUsersAuth(t *testing.T) {

	tCases := []TestCase{
		{
			handler: mockedUserHandlers.CreateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       2,
					Username: "username2",
					Email:    "mail2",
				},
			},

			expStatus:    http.StatusMethodNotAllowed,
			inputMessage: []byte(`{"email": "new_user@mail.com","username": "new_user", "password": "pass"}`),

			expMessage: models.AlreadySignedInAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}

		// testLog(t, tCase)
	}
}

func TestCreateUserIncorrectJson(t *testing.T) {

	tCases := []TestCase{
		{
			handler: mockedUserHandlers.CreateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			expStatus:    http.StatusInternalServerError,
			inputMessage: []byte(`"email": "new_user@mail.com","username": "new_user", "password": "pass"}`), //no opening bracket

			expMessage: models.IncorrectJsonAnswer,
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}
	}
}

func TestCreateUserIncorrectFields(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.CreateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			expStatus:    http.StatusBadRequest,
			inputMessage: []byte(`{"email182": "new_user@mail.com","user___name": "new_user", "pass1029word": "pass"}`),

			expMessage: models.GetIncorrectFieldsAnswer([]string{"email", "username", "password"}),
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.IncorrectFieldsAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}
	}
}

func TestCreateUserAlreadyExists(t *testing.T) {
	tCases := []TestCase{
		{
			handler: mockedUserHandlers.CreateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			expStatus:    http.StatusInternalServerError,
			inputMessage: []byte(`{"email": "new_user@mail.com","username": "new_user", "password": "pass"}`),

			expMessage: models.GetDeveloperErrorAnswer("already exists"),
		},
	}

	for _, tCase := range tCases {
		w, req := testInitial(tCase)
		tCase.handler(w, req)

		if w.Code != tCase.expStatus {
			t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
		}
		var retMessage models.MessageAnswer
		_ = json.NewDecoder(w.Body).Decode(&retMessage)
		if !reflect.DeepEqual(retMessage, tCase.expMessage) {
			t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", retMessage, tCase.expMessage)
		}
	}
}
