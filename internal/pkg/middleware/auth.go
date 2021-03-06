package middleware

import (
	"context"
	"net/http"

	"2019_1_OPG_plus_2/internal/pkg/auth"
	"2019_1_OPG_plus_2/internal/pkg/models"
)

func AuthMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		jwtCookie, errNoCookie := req.Cookie(auth.CookieName)
		if errNoCookie != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			ctx = context.WithValue(ctx, "jwtData", models.JwtData{})
		} else {
			data, err := auth.CheckJwt(jwtCookie.Value)
			if err != nil {
				ctx = context.WithValue(ctx, "isAuth", false)
				ctx = context.WithValue(ctx, "jwtData", models.JwtData{})
			} else {
				ctx = context.WithValue(ctx, "isAuth", true)
				ctx = context.WithValue(ctx, "jwtData", data)
			}
		}
		h.ServeHTTP(res, req.WithContext(ctx))
	}

	return mw
}
