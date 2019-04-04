package controllers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	a "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/adapters"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

func init() {
	a.SetStorages(NewMockStorage(), auth.NewStorage())
	a.SetHandlers(NewUserHandlers(), NewAuthHandlers())
}

/**********************
 *  IS_AUTH CONTROLLER*
 **********************/

func TestIsAuth(t *testing.T) {
	tCases := []TestCase{
		{
			handler: a.GetHandlers().Auth.IsAuth,

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
			handler: a.GetHandlers().Auth.IsAuth,

			params: TestParams{
				muxVars: map[string]string{},
				method:  "GET",
				isAuth:  false,
				url:     "/session",
				jwt:     models.JwtData{},
			},

			expStatus:    http.StatusUnauthorized,
			inputMessage: nil,

			expMessage: models.SignedOutAnswer,
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
