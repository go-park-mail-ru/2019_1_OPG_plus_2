package cookieservice

import (
	"2019_1_OPG_plus_2/internal/pkg/config"
	"2019_1_OPG_plus_2/internal/pkg/cookiecheckerproto"
	"2019_1_OPG_plus_2/internal/pkg/models"
	"context"
	"log"
)

type Server struct{}

func (Server) CheckCookie(ctx context.Context, req *cookiecheckerproto.CookieRequest) (*cookiecheckerproto.CookieResponse, error) {
	log.Println("call to CheckCookie")
	data := models.JwtData{}
	err := data.UnMarshal(req.JwtToken, []byte(config.Auth.Secret))

	res := &cookiecheckerproto.CookieResponse{}

	if err == nil {
		res.Error = ""
	} else {
		res.Error = err.Error()
		return res, nil
	}

	res = &cookiecheckerproto.CookieResponse{
		Data: &cookiecheckerproto.JwtData{
			Exp:      data.ExpiresAt,
			Username: data.Username,
			Id:       data.Id,
			Email:    data.Email,
		},
	}

	return res, nil
}
