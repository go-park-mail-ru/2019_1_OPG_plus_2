package user

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

type IUser interface {
	CreateUser(signUpData models.SingUpData) (jwtData models.JwtData, err error)
	GetUser(id int64) (userData models.GetUserAnswer, err error)
	UpdateUser(id int64, updateData models.UpdateUserData) (jwtData models.JwtData, err error)
	RemoveUser(id int64, removeData models.RemoveUserData) error
}

type StorageAdapter struct{}

func NewStorageAdapter() *StorageAdapter {
	return &StorageAdapter{}
}
func (*StorageAdapter) CreateUser(signUpData models.SingUpData) (jwtData models.JwtData, err error) {
	jwtData, err = auth.SignUp(signUpData)
	if err != nil {
		return
	}

	err = db.ProfileCreate(db.ProfileData{
		Id:     jwtData.Id,
		Avatar: signUpData.Avatar,
	})
	return
}
func (*StorageAdapter) GetUser(id int64) (userData models.GetUserAnswer, err error) {
	userData, err = db.GetUser(id)
	if err == sql.ErrNoRows {
		return userData, fmt.Errorf("user not found")
	}
	return
}
func (*StorageAdapter) UpdateUser(id int64, updateData models.UpdateUserData) (jwtData models.JwtData, err error) {
	jwtData, err = auth.UpdateAuth(id, updateData)
	return
}
func (*StorageAdapter) RemoveUser(id int64, removeData models.RemoveUserData) error {
	err := auth.RemoveAuth(id, removeData)
	if err != nil {
		return err
	}
	err = db.ProfileRemove(id)
	return err
}
