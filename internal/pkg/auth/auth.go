package auth

import (
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	authService "2019_1_OPG_plus_2/internal/proto"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net/http"
	"time"

	"2019_1_OPG_plus_2/internal/pkg/models"
)

var Manager authManager

func init() {
	conn, err := grpc.Dial("localhost:50242", grpc.WithInsecure())
	if err != nil {
		tsLogger.LogErr("AUTH: can not connect to service [%v]", err)
	}

	Manager = authManager{
		Conn:       conn,
		AuthClient: authService.NewAuthServiceClient(conn),
	}
}

func CreateAuthCookie(data models.JwtData, lifetime time.Duration) *http.Cookie {
	jwtStr, err := data.Marshal(lifetime, secret)
	if err != nil {
		return &http.Cookie{}
	}

	return &http.Cookie{
		Name:     CookieName,
		Value:    jwtStr,
		Expires:  time.Now().Add(lifetime),
		HttpOnly: true,
	}
}

func CheckJwt(token string) (models.JwtData, error) {
	data := models.JwtData{}
	err := data.UnMarshal(token, secret)
	return data, err
}

func PasswordHash(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func NewStorage() *Storage {
	return &Storage{}
}

type Storage struct{}

func (*Storage) SignUp(signUpData models.SignUpData) (models.JwtData, error, []string) {
	data := &authService.SignUpRequest{
		Data: &authService.SignUpData{
			Username: signUpData.Username,
			Password: signUpData.Password,
			Email:    signUpData.Email,
		},
	}

	response, err := Manager.AuthClient.SignUp(context.Background(), data)
	if err != nil {
		tsLogger.LogErr("AUTH: SignUp call ended in: %v", err)
		return models.JwtData{}, err, nil
	}

	var reterr error
	if response.Error != "" {
		reterr = errors.New(response.Error)
	} else {
		reterr = nil
	}
	responseData, err := CheckJwt(response.JwtToken)
	if err != nil {
		panic(err)
	}

	return responseData, reterr, response.Fields
}

func (*Storage) SignIn(signInData models.SignInData) (models.JwtData, error, []string) {
	data := &authService.SignInRequest{
		Data: &authService.SignInData{
			Password: signInData.Password,
			Login:    signInData.Login,
		},
	}

	response, err := Manager.AuthClient.SignIn(context.Background(), data)
	if err != nil {
		tsLogger.LogErr("AUTH: SignIn call ended in: %v", err)
		return models.JwtData{}, err, nil
	}

	var reterr error
	if response.Error != "" {
		reterr = errors.New(response.Error)
	} else {
		reterr = nil
	}
	responseData, err := CheckJwt(response.JwtToken)
	if err != nil {
		panic(err)
	}

	return responseData, reterr, response.Fields
	//return SignIn(signInData)
}

func (*Storage) UpdateAuth(id int64, userData models.UpdateUserData) (models.JwtData, error, []string) {
	data := &authService.UpdateAuthRequest{
		Id: id,
		UserData: &authService.UpdateUserData{
			Email:    userData.Email,
			Username: userData.Username,
		},
	}

	response, err := Manager.AuthClient.UpdateAuth(context.Background(), data)
	if err != nil {
		tsLogger.LogErr("AUTH: UpdateAuth call ended in: %v", err)
		return models.JwtData{}, err, nil
	}

	var reterr error
	if response.Error != "" {
		reterr = errors.New(response.Error)
	} else {
		reterr = nil
	}
	responseData, err := CheckJwt(response.JwtToken)
	if err != nil {
		panic(err)
	}

	return responseData, reterr, response.Fields
}

func (*Storage) UpdatePassword(id int64, passwordData models.UpdatePasswordData) (error, []string) {
	return UpdatePassword(id, passwordData)
}

func (*Storage) RemoveAuth(id int64, removeData models.RemoveUserData) (error, []string) {
	return RemoveAuth(id, removeData)
}

type authManager struct {
	Conn       *grpc.ClientConn
	AuthClient authService.AuthServiceClient
}
