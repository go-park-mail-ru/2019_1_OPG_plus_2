package adapters

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

var (
	handlers Handlers
	storages Storages
)

func SetHandlers(userHandlers IUserHandlers, authHandlers IAuthHandlers) {
	handlers = Handlers{User: userHandlers, Auth: authHandlers}
}

func GetHandlers() *Handlers {
	return &handlers
}

func SetStorages(userStorage IUserStorage, authStorage IAuthStorage) {
	storages = Storages{User: userStorage, Auth: authStorage}
}

func GetStorages() *Storages {
	return &storages
}

type Handlers struct {
	User IUserHandlers
	Auth IAuthHandlers
}

type IUserHandlers interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	RemoveUser(w http.ResponseWriter, r *http.Request)
}

type IAuthHandlers interface {
	IsAuth(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
	UpdatePassword(w http.ResponseWriter, r *http.Request)
}

type IOAuthHandlers interface {
	Login1stStageRetrieveCode(w http.ResponseWriter, r *http.Request)
	Login2ndStageRetrieveTokenGetData(w http.ResponseWriter, r *http.Request)
}

type Storages struct {
	User IUserStorage
	Auth IAuthStorage
}

type IUserStorage interface {
	CreateUser(signUpData models.SignUpData) (models.JwtData, error, []string)
	GetUser(id int64) (models.UserData, error)
	UpdateUser(id int64, updateData models.UpdateUserData) (models.JwtData, error, []string)
	RemoveUser(id int64, removeData models.RemoveUserData) (error, []string)
}

type IAuthStorage interface {
	SignUp(signUpData models.SignUpData) (models.JwtData, error, []string)
	SignIn(signInData models.SignInData) (data models.JwtData, err error, incorrectFields []string)
	UpdatePassword(id int64, passwordData models.UpdatePasswordData) (error, []string)
	UpdateAuth(id int64, userData models.UpdateUserData) (models.JwtData, error, []string)
	RemoveAuth(id int64, removeData models.RemoveUserData) (error, []string)
}
