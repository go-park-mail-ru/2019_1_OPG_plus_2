package controllers

import (
	"net/http"
	"testing"

	a "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/adapters"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

func init() {
	testInitial()
}

func TestIsAuth(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().Auth.IsAuth,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  true,
				url:     "/session",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  nil,
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusOK,
			expMessage: &models.SignedInAnswer,
		},
		{
			handler: a.GetHandlers().Auth.IsAuth,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/session",
				jwt:     models.JwtData{},
			},

			inputMessage:  nil,
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusUnauthorized,
			expMessage: &models.SignedOutAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestSignInCorrect(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().Auth.SignIn,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  true,
				url:     "/session",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  nil,
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusMethodNotAllowed,
			expMessage: &models.AlreadySignedInAnswer,
		},
		{
			handler: a.GetHandlers().Auth.SignIn,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/session",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`{"login": "username1","password": "pass1"}`),
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusOK,
			expMessage: &models.SignedInAnswer,
		},
		{
			handler: a.GetHandlers().Auth.SignIn,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/session",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`{"login": "mail1@mail.ru","password": "pass1"}`),
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusOK,
			expMessage: &models.SignedInAnswer,
		},
		{
			handler: a.GetHandlers().Auth.SignIn,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/session",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`{"login": "username1","password": "unknown"}`),
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  http.StatusBadRequest,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"password"}),
		},
		{
			handler: a.GetHandlers().Auth.SignIn,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/session",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`{"login": "unknown","password": "pass1"}`),
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  http.StatusBadRequest,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"password"}),
		},
		{
			handler: a.GetHandlers().Auth.SignIn,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/session",
				jwt:     models.JwtData{},
			},

			inputMessage:  []byte(`{"login": "No username or email!","password": "pass1"}`),
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  http.StatusBadRequest,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"login"}),
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestSignOutCorrect(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().Auth.SignOut,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  true,
				url:     "/session",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  nil,
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusOK,
			expMessage: &models.SignedOutAnswer,
		},
		{
			handler: a.GetHandlers().Auth.SignOut,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/session",
				jwt:     models.JwtData{},
			},

			inputMessage:  nil,
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusMethodNotAllowed,
			expMessage: &models.AlreadySignedOutAnswer,
		},
	}
	
	for _, tCase := range tCases {
		test(t, &tCase)
	}
}

func TestUpdatePasswordCorrect(t *testing.T) {
	tCases := []testCase{
		{
			handler: a.GetHandlers().Auth.UpdatePassword,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/session",
				jwt:     models.JwtData{},
			},

			inputMessage:  nil,
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusUnauthorized,
			expMessage: &models.NotSignedInAnswer,
		},
		{
			handler: a.GetHandlers().Auth.UpdatePassword,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  true,
				url:     "/session",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  []byte(`{"new_password": "new_pass","password_confirm": "new_pass"}`),
			outputMessage: &models.MessageAnswer{},

			expStatus:  http.StatusOK,
			expMessage: &models.PasswordUpdatedAnswer,
		},
		{
			handler: a.GetHandlers().Auth.UpdatePassword,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  true,
				url:     "/session",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  []byte(`{"new_password": "","password_confirm": ""}`),
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  http.StatusBadRequest,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"new_password"}),
		},
		{
			handler: a.GetHandlers().Auth.UpdatePassword,

			params: testParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  true,
				url:     "/session",
				jwt: models.JwtData{
					Id:       1,
					Username: "username1",
					Email:    "mail1@mail.ru",
				},
			},

			inputMessage:  []byte(`{"new_password": "new2","password_confirm": "new3"}`),
			outputMessage: &models.IncorrectFieldsAnswer{},

			expStatus:  http.StatusBadRequest,
			expMessage: models.GetIncorrectFieldsAnswer([]string{"password_confirm"}),
		},
	}
	
	for _, tCase := range tCases {
		test(t, &tCase)
	}
}
