package controllers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	a "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/adapters"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

func init() {
	a.SetStorages(newMockStorage(), auth.NewStorage())
	a.SetHandlers(NewUserHandlers(), NewAuthHandlers())
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
	userToFind, _ := a.GetStorages().User.GetUser(1)
	tCases := []TestCase{
		{
			handler: a.GetHandlers().User.GetUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  true,
				url:     "/User",
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
	userToFind, _ := a.GetStorages().User.GetUser(3)
	tCases := []TestCase{
		{
			handler: a.GetHandlers().User.GetUser,

			params: TestParams{
				muxVars: map[string]string{"id": "3"},
				method:  "GET",
				isAuth:  true,
				url:     "/User",
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
		// // testLog(t, tCase)
		// CheckStatusAndAnswer(t, w, retMessage, tCase)
	}
}

func TestGetUserIdNotExists(t *testing.T) {
	tCases := []TestCase{
		{
			handler: a.GetHandlers().User.GetUser,

			params: TestParams{
				muxVars: map[string]string{"id": "1278"},
				method:  "GET",
				isAuth:  true,
				url:     "/User",
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
			handler: a.GetHandlers().User.GetUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/User",
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
			handler: a.GetHandlers().User.GetUser,

			params: TestParams{
				muxVars: map[string]string{"id": "qwerty"},
				method:  "GET",
				isAuth:  true,
				url:     "/User",
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
			handler: a.GetHandlers().User.UpdateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/User",
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

		newUserData, err := a.GetStorages().User.GetUser(tCase.params.jwt.Id)
		if err != nil {
			t.Errorf("Test failed while getting User from storage: %e", err)
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
			handler: a.GetHandlers().User.UpdateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  false,
				url:     "/User",
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
			handler: a.GetHandlers().User.UpdateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/User",
				jwt: models.JwtData{
					Id:       1,
					Username: "qwerty",
					Email:    "qwerty@mail.com",
				},
			},

			expStatus:    400,
			inputMessage: []byte(`{"email": "qwerty@mail.com","User": "qwerty"}`),

			expMessage: models.GetIncorrectFieldsAnswer([]string{"username"}),
		},
		{
			handler: a.GetHandlers().User.UpdateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/User",
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
			handler: a.GetHandlers().User.UpdateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/User",
				jwt: models.JwtData{
					Id:       1,
					Username: "qwerty",
					Email:    "qwerty@mail.com",
				},
			},

			expStatus:    400,
			inputMessage: []byte(`{"email": "qwerty","User-name": "qwerty"}`),

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
			handler: a.GetHandlers().User.UpdateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/User",
				jwt: models.JwtData{
					Id:       1,
					Username: "qwerty",
					Email:    "qwerty@mail.com",
				},
			},

			expStatus:    500,
			inputMessage: []byte(`{"email": "qwerty@mail.com","User": "qwerty"`), // no closing parentheses in JSON

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
			handler: a.GetHandlers().User.RemoveUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  true,
				url:     "/User",
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

		_, err := a.GetStorages().User.GetUser(tCase.params.jwt.Id)
		if err == nil {
			t.Errorf("Data did not actually delete")
		}
	}
}

// since this moment i would use another User as a client for testing
// because in previous test User 1 has been deleted successfully
func TestRemoveUserNoAuth(t *testing.T) {
	tCases := []TestCase{
		{
			handler: a.GetHandlers().User.RemoveUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  false,
				url:     "/User",
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
			handler: a.GetHandlers().User.RemoveUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  true,
				url:     "/User",
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
			handler: a.GetHandlers().User.RemoveUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  true,
				url:     "/User",
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

		_, err := a.GetStorages().User.GetUser(tCase.params.jwt.Id)
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
			handler: a.GetHandlers().User.CreateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/User",
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

		// checking if User is actually created

		// getting data field from cookie set to User
		rawJwtData := strings.Split(w.Result().Cookies()[0].Value, ".")[1]
		parsedCookie, _ := base64.StdEncoding.DecodeString(rawJwtData)

		// parsing data of User stored in cache into struct
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
			handler: a.GetHandlers().User.CreateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  true,
				url:     "/User",
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
			handler: a.GetHandlers().User.CreateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/User",
				jwt:     models.JwtData{},
			},

			expStatus:    http.StatusInternalServerError,
			inputMessage: []byte(`"email": "new_user@mail.com","username": "new_user", "password": "pass"}`), // no opening bracket

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
			handler: a.GetHandlers().User.CreateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/User",
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
			handler: a.GetHandlers().User.CreateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/User",
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
