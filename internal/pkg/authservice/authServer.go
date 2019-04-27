package authservice

import (
	"2019_1_OPG_plus_2/internal/pkg/auth"
	"2019_1_OPG_plus_2/internal/pkg/config"
	"2019_1_OPG_plus_2/internal/pkg/models"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	authService "2019_1_OPG_plus_2/internal/proto"
	"context"
	"fmt"
	"time"
)

type Server struct {
	Log *tsLogger.TSLogger
}

func NewServer() *Server {
	return &Server{
		Log: tsLogger.NewLogger(),
	}
}

func (s *Server) SignUp(ctx context.Context, request *authService.SignUpRequest) (*authService.SignUpResponse, error) {
	s.Log.LogTrace("AUTH: call to SignUp RPC")
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

func (s *Server) SignIn(ctx context.Context, request *authService.SignInRequest) (*authService.SignInResponse, error) {
	s.Log.LogAcc("AUTH: call to SignIn RPC")
	data := models.SignInData{
		Login:    request.Data.GetLogin(),
		Password: request.Data.GetPassword(),
	}

	jwtData, err, fields := auth.SignIn(data)
	token, er := jwtData.Marshal(30*24*time.Hour, []byte(config.Auth.Secret))
	if er != nil {
		panic(er)
	}

	if err == nil {
		err = fmt.Errorf("")
	}

	response := &authService.SignInResponse{
		Error:    err.Error(),
		Fields:   fields,
		JwtToken: token,
	}
	return response, nil
}

func (s *Server) UpdateAuth(ctx context.Context, request *authService.UpdateAuthRequest) (*authService.UpdateAuthResponse, error) {
	s.Log.LogAcc("AUTH: call to UpdateAuth RPC")

	data := models.UpdateUserData{
		Username: request.UserData.GetUsername(),
		Email:    request.UserData.GetEmail(),
	}

	jwtData, err, fields := auth.UpdateAuth(request.Id, data)
	token, er := jwtData.Marshal(30*24*time.Hour, []byte(config.Auth.Secret))
	if er != nil {
		panic(er)
	}

	if err == nil {
		err = fmt.Errorf("")
	}

	response := &authService.UpdateAuthResponse{
		Error:    err.Error(),
		Fields:   fields,
		JwtToken: token,
	}
	return response, nil
}

func (s *Server) UpdatePassword(ctx context.Context, request *authService.UpdatePasswordRequest) (*authService.UpdatePasswordResponse, error) {
	s.Log.LogAcc("AUTH: call to UpdatePassword RPC")

	data := models.UpdatePasswordData{
		NewPassword:     request.PasswordData.GetNewPassword(),
		PasswordConfirm: request.PasswordData.GetPasswordConfirm(),
	}

	err, fields := auth.UpdatePassword(request.Id, data)
	if err == nil {
		err = fmt.Errorf("")
	}

	response := &authService.UpdatePasswordResponse{
		Error:  err.Error(),
		Fields: fields,
	}
	return response, nil
}

func (s *Server) RemoveAuth(ctx context.Context, request *authService.RemoveAuthRequest) (*authService.RemoveAuthResponse, error) {
	s.Log.LogAcc("AUTH: call to RemoveAuth RPC")

	data := models.RemoveUserData{
		Password: request.RemoveData.GetPassword(),
	}

	err, fields := auth.RemoveAuth(request.Id, data)
	if err == nil {
		err = fmt.Errorf("")
	}

	response := &authService.RemoveAuthResponse{
		Error:  err.Error(),
		Fields: fields,
	}
	return response, nil
}
