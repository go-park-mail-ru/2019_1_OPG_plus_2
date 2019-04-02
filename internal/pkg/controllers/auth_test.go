package controllers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
	"reflect"
	"testing"
)

//var mockedStorageAdapter = newMockStorageAdapter()
//var mockedUserHandlers = NewUserHandlers(mockedStorageAdapter)

/**********************
 *  IS_AUTH CONTROLLER*
 **********************/

func TestIsAuth(t *testing.T) {
	tCases := []TestCase{
		{
			handler: IsAuth,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  true,
				url:     "/session",
				jwt: models.JwtData{
					Id:       1,
					Username: "qwerty",
					Email:    "qwerty@mail.com",
				},
			},

			expStatus:    http.StatusOK,
			inputMessage: nil,

			expMessage: models.SignedInAnswer,
		},
		{
			handler: mockedUserHandlers.UpdateUser,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/session",
				jwt:     models.JwtData{},
			},

			expStatus:    http.StatusUnauthorized,
			inputMessage: nil,

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
