package auth

import (
	authproto "2019_1_OPG_plus_2/internal/pkg/proto"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
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
	url := serviceLocation + ":" + port
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		tsLogger.LogErr("AUTH: can not connect to service [%v]", err)
	}

	Manager = authManager{
		Conn:       conn,
		AuthClient: authproto.NewAuthServiceClient(conn),
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
	data := &authproto.SignUpRequest{
		Data: &authproto.SignUpData{
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
	data := &authproto.SignInRequest{
		Data: &authproto.SignInData{
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
	data := &authproto.UpdateAuthRequest{
		Id: id,
		UserData: &authproto.UpdateUserData{
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
	data := &authproto.UpdatePasswordRequest{
		Id: id,
		PasswordData: &authproto.UpdatePasswordData{
			NewPassword:     passwordData.NewPassword,
			PasswordConfirm: passwordData.PasswordConfirm,
		},
	}
	response, err := Manager.AuthClient.UpdatePassword(context.Background(), data)
	if err != nil {
		tsLogger.LogErr("AUTH: UpdatePassword call ended in: %v", err)
		return err, nil
	}

	var reterr error
	if response.Error != "" {
		reterr = errors.New(response.Error)
	} else {
		reterr = nil
	}

	return reterr, response.Fields
}

func (*Storage) RemoveAuth(id int64, removeData models.RemoveUserData) (error, []string) {
	data := &authproto.RemoveAuthRequest{
		Id: id,
		RemoveData: &authproto.RemoveUserData{
			Password: removeData.Password,
		},
	}

	response, err := Manager.AuthClient.RemoveAuth(context.Background(), data)
	if err != nil {
		tsLogger.LogErr("AUTH: RemoveAuth call ended in: %v", err)
		return err, nil
	}

	var reterr error
	if response.Error != "" {
		reterr = errors.New(response.Error)
	} else {
		reterr = nil
	}

	return reterr, response.Fields
	//return RemoveAuth(id, removeData)
}

type authManager struct {
	Conn         *grpc.ClientConn
	AuthClient   authproto.AuthServiceClient
	CookieClient authproto.CookieCheckerClient
}
