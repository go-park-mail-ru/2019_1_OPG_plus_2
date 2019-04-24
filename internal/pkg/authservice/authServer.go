package authservice

import (
	"2019_1_OPG_plus_2/internal/pkg/auth"
	"2019_1_OPG_plus_2/internal/pkg/config"
	"2019_1_OPG_plus_2/internal/pkg/models"
	authService "2019_1_OPG_plus_2/internal/proto"
	"context"
	"fmt"
	"time"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (*Server) SignUp(ctx context.Context, request *authService.SignUpRequest) (*authService.SignUpResponse, error) {
	data := models.SignUpData{
		Email:    request.Data.GetEmail(),
		Password: request.Data.GetPassword(),
		Username: request.Data.GetUsername(),
	}
	jwtData, err, fields := auth.SignUp(data)
	token, er := jwtData.Marshal(30*24*time.Hour, []byte(config.Auth.Secret))
	if er != nil {
		panic(er)
	}

	if err == nil {
		err = fmt.Errorf("")
	}

	response := &authService.SignUpResponse{
		Error:    err.Error(),
		Fields:   fields,
		JwtToken: token,
	}
	return response, nil
}

func (*Server) SignIn(context.Context, *authService.SignInRequest) (*authService.SignInResponse, error) {
	panic("implement me")
}

func (*Server) UpdateAuth(context.Context, *authService.UpdateAuthRequest) (*authService.UpdateAuthResponse, error) {
	panic("implement me")
}

func (*Server) UpdatePassword(context.Context, *authService.UpdatePasswordRequest) (*authService.UpdatePasswordResponse, error) {
	panic("implement me")
}

func (*Server) RemoveAuth(context.Context, *authService.RemoveAuthRequest) (*authService.RemoveAuthResponse, error) {
	panic("implement me")
}
