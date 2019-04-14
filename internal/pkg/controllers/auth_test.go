package controllers

import (
	"net/http"
	"testing"

	a "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/adapters"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

func init()  {
	testInitial()
}

/**********************
 *  IS_AUTH CONTROLLER*
 **********************/

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
					Username: "qwerty",
					Email:    "qwerty@mail.com",
				},
			},

			inputMessage: nil,
			outputMessage: &models.MessageAnswer{},

			expStatus:    http.StatusOK,
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

			inputMessage: nil,
			outputMessage: &models.MessageAnswer{},

			expStatus:    http.StatusUnauthorized,
			expMessage: &models.SignedOutAnswer,
		},
	}

	for _, tCase := range tCases {
		test(t, &tCase)
	}
}
