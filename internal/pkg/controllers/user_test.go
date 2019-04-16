package controllers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	a "2019_1_OPG_plus_2/internal/pkg/adapters"
	"2019_1_OPG_plus_2/internal/pkg/models"
)

func init() {
	testsInitial()
}

/*************************
 *  GET_USER CONTROLLER  *
 *************************/

func TestGetUserSelf(t *testing.T) {
	userToFind, _ := a.GetStorages().User.GetUser(1)
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.GetUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  nil,
			outputMessage: &models.UserDataAnswer{},

			expStatus:  200,
			expMessage: models.GetUserDataAnswer(userToFind),
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestGetUserId(t *testing.T) {
	userToFind, _ := a.GetStorages().User.GetUser(3)
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.GetUser,

			params: testParams{
				muxVars: map[string]string{"id": "3"},
				method:  "GET",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  nil,
			outputMessage: &models.UserDataAnswer{},

			expStatus:  200,
			expMessage: models.GetUserDataAnswer(userToFind),
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestGetUserIdNotExists(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.GetUser,

			params: testParams{
				muxVars: map[string]string{"id": "1278"},
				method:  "GET",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  nil,
			outputMessage: &models.MessageAnswer{},

			expStatus:  404,
			expMessage: &models.UserNotFoundAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestGetUserNoAuth(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.GetUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			inputMessage:  nil,
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  400,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"id"}),
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestGetUserIdNotNumber(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.GetUser,

			params: testParams{
				muxVars: map[string]string{"id": "qwerty"},
				method:  "GET",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  nil,
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  400,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"id"}),
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

/****************************
 *  UPDATE_USER CONTROLLER  *
 ****************************/

func TestUpdateUserCorrect(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.UpdateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  []byte(`{"email": "qwerty@mail.com","username": "qwerty"}`),
			outputMessage: &models.MessageAnswer{},

			expStatus:  200,
			expMessage: &models.UserUpdatedAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)

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
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.UpdateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "PUT",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`{"email": "qwerty","username": "qwerty"}`),
			outputMessage: &models.MessageAnswer{},

			expStatus:  401,
			expMessage: &models.NotSignedInAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestUpdateUserInvalidField(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.UpdateUser,

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

			inputMessage:  []byte(`{"email": "qwerty@mail.com","User": "qwerty"}`),
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  400,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"username"}),
		},
		{
			handler: a.GetHandlers().User.UpdateUser,

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

			inputMessage:  []byte(`{"e-mail": "qwerty@mail.com","username": "qwerty"}`),
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  400,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"email"}),
		},
		{
			handler: a.GetHandlers().User.UpdateUser,

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

			inputMessage:  []byte(`{"email": "qwerty","User-name": "qwerty"}`),
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  400,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"email", "username"}),
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestUpdateUserInvalidJSON(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.UpdateUser,

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

			inputMessage:  []byte(`{"email": "qwerty@mail.com","User": "qwerty"`), // no closing parentheses in JSON
			outputMessage: &models.MessageAnswer{},

			expStatus:  500,
			expMessage: &models.IncorrectJsonAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

/****************************
 *  DELETE_USER CONTROLLER  *
 ****************************/

func TestRemoveUserCorrect(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.RemoveUser,

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

			inputMessage:  []byte(`{"password": "new_pass"}`),
			outputMessage: &models.MessageAnswer{},

			expStatus:  200,
			expMessage: &models.UserRemovedAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)

		_, err := a.GetStorages().User.GetUser(tCase.params.jwt.Id)
		if err == nil {
			t.Errorf("Data did not actually delete")
		}
	}
}

func TestRemoveUserNoAuth(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.RemoveUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`{"password": "pass2"}`),
			outputMessage: &models.MessageAnswer{},

			expStatus:  401,
			expMessage: &models.NotSignedInAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestRemoveUserInvalidJSON(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.RemoveUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       2,
					Username: "username2",
					Email:    "mail2@mail.ru",
				},
			},

			inputMessage:  []byte(`{"password": "pass1"`), // no closing parentheses
			outputMessage: &models.MessageAnswer{},

			expStatus:  500,
			expMessage: &models.IncorrectJsonAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestRemoveUserInvalidData(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.RemoveUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "DELETE",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       2,
					Username: "username2",
					Email:    "mail2@mail.ru",
				},
			},

			inputMessage:  []byte(`{"passw": "pass1"}`),
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  400,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"password"}),
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)

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
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.CreateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`{"email": "new_user@mail.com","username": "new_user", "password": "pass"}`),
			outputMessage: &models.MessageAnswer{},

			expStatus:  200,
			expMessage: &models.SignedUpAnswer,
		},
	}

	for _, tCase := range tCases {
		w, _ := test(t, &tCase)

		// checking if User is actually created

		// getting data field from cookie set to User
		rawJwtData := strings.Split(w.Result().Cookies()[0].Value, ".")[1]
		parsedCookie, _ := base64.StdEncoding.DecodeString(rawJwtData)

		// parsing data of User stored in cache into struct
		var storedUserData models.JwtData
		_ = json.Unmarshal(parsedCookie, &storedUserData)

		var inputData models.SignUpData
		_ = json.Unmarshal(tCase.inputMessage, &inputData)

		if storedUserData.Email != inputData.Email || storedUserData.Username != inputData.Username {
			t.Error("User saved incorrectly")
		}

		// testLog(t, tCase)
	}
}

func TestCreateUsersAuth(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.CreateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  true,
				url:     "/user",
				jwt: models.JwtData{
					Id:       2,
					Username: "username2",
					Email:    "mail2@mail.ru",
				},
			},

			inputMessage:  []byte(`{"email": "new_user@mail.com","username": "new_user", "password": "pass"}`),
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusMethodNotAllowed,
			expMessage: &models.AlreadySignedInAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestCreateUserIncorrectJson(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.CreateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`"email": "new_user@mail.com","username": "new_user", "password": "pass"}`), // no opening bracket
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusInternalServerError,
			expMessage: &models.IncorrectJsonAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestCreateUserIncorrectFields(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.CreateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`{"email182": "new_user@mail.com","user___name": "new_user", "pass1029word": "pass"}`),
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  http.StatusBadRequest,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"email", "username", "password"}),
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestCreateUserAlreadyExists(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().User.CreateUser,

			params: testParams{
				muxVars: map[string]string{},
				method:  "POST",
				isAuth:  false,
				url:     "/user",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`{"email": "new_user@mail.com","username": "new_user", "password": "pass"}`),
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusInternalServerError,
			expMessage: models.GetDeveloperErrorAnswer("already exists"),
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}
