package user

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

type IUser interface {
	CreateUser(signUpData models.SingUpData) (models.JwtData, error, []string)
	GetUser(id int64) (models.UserData, error)
	UpdateUser(id int64, updateData models.UpdateUserData) (models.JwtData, error, []string)
	RemoveUser(id int64, removeData models.RemoveUserData) (error, []string)
}

type StorageAdapter struct{}

func NewStorageAdapter() *StorageAdapter {
	return &StorageAdapter{}
}
func (*StorageAdapter) CreateUser(signUpData models.SingUpData) (jwtData models.JwtData, err error, fields []string) {
	incorrectFields := signUpData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}

	jwtData, err, fields = auth.SignUp(signUpData)
	if err != nil {
		return
	}

	err = db.ProfileCreate(db.ProfileData{
		Id: jwtData.Id,
	})
	return
}
func (*StorageAdapter) GetUser(id int64) (userData models.UserData, err error) {
	userData, err = db.GetUser(id)
	if err == sql.ErrNoRows {
		return userData, models.NotFound
	}
	return
}
func (*StorageAdapter) UpdateUser(id int64, updateData models.UpdateUserData) (jwtData models.JwtData, err error, fields []string) {
	return auth.UpdateAuth(id, updateData)
}
func (*StorageAdapter) RemoveUser(id int64, removeData models.RemoveUserData) (error, []string) {
	incorrectFields := removeData.Check()
	if len(incorrectFields) > 0 {
		return models.FieldsError, incorrectFields
	}

	err, fields := auth.RemoveAuth(id, removeData)
	if err != nil {
		return err, fields
	}

	err = db.ProfileRemove(id)
	return err, nil
}
