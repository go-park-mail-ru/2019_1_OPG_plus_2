package cookieservice

import (
	"2019_1_OPG_plus_2/internal/pkg/config"
	"2019_1_OPG_plus_2/internal/pkg/cookiecheckerproto"
	"2019_1_OPG_plus_2/internal/pkg/models"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"context"
)

type Server struct {
	Log *tsLogger.TSLogger
}

func NewServer(log *tsLogger.TSLogger) *Server {
	return &Server{Log: log}
}

func (s Server) CheckCookie(ctx context.Context, req *cookiecheckerproto.CookieRequest) (*cookiecheckerproto.CookieResponse, error) {
	s.Log.LogTrace("COOKIE: Call to CheckCookie RPC")
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
