package auth

import (
	"2019_1_OPG_plus_2/internal/pkg/authproto"
	"2019_1_OPG_plus_2/internal/pkg/cookiecheckerproto"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"net/http"
	"time"

	"2019_1_OPG_plus_2/internal/pkg/models"
)

var Manager authManager

// TODO: dial connections in runtime, not in init

// TODO: leave connections persistent, monitor connection and ping it sometimes, reconnect
func init() {

	authurl := serviceLocation + ":" + port
	authconn, err := grpc.Dial(authurl, grpc.WithInsecure())
	if err != nil {
		tsLogger.LogErr("AUTH: can not connect to service [%v]", err)
	}

	cookieurl := cookielocation + ":" + cookieport
	cookieconn, err := grpc.Dial(cookieurl, grpc.WithInsecure())
	if err != nil {
		tsLogger.LogErr("СOOKIE: can not connect to service [%v]", err)
	}

	Manager = authManager{
		AuthConn:     authconn,
		CookieConn:   cookieconn,
		AuthClient:   authproto.NewAuthServiceClient(authconn),
		CookieClient: cookiecheckerproto.NewCookieCheckerClient(cookieconn),
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

// RPC
func CheckJwt(token string) (models.JwtData, error) {
	//data := models.JwtData{}
	//err := data.UnMarshal(token, secret)
	//return data, err

	req := &cookiecheckerproto.CookieRequest{
		JwtToken: token,
	}

	res, err := Manager.CookieClient.CheckCookie(context.Background(), req)
	if err != nil {
		tsLogger.LogErr("AUTH: CheckCookie call ended in: %v", err)
		return models.JwtData{}, err
	}

	data := models.JwtData{
		Email:    res.Data.GetEmail(),
		Id:       res.Data.GetId(),
		Username: res.Data.GetUsername(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: res.Data.GetExp(),
		},
	}

	err = errors.New(res.Error)
	if res.Error == "" {
		err = nil
	}
	return data, err
}

func PasswordHash(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func NewStorage() *Storage {
	return &Storage{}
}

type Storage struct{}

//RPC
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

//RPC
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

//RPC
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

//RPC
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

//RPC
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
	AuthConn     *grpc.ClientConn
	CookieConn   *grpc.ClientConn
	AuthClient   authproto.AuthServiceClient
	CookieClient cookiecheckerproto.CookieCheckerClient
}
