package user

import (
	"2019_1_OPG_plus_2/internal/pkg/adapters"

	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/models"
)

func NewStorage() *Storage {
	return &Storage{}
}

type Storage struct{}

func (*Storage) CreateUser(signUpData models.SignUpData) (models.JwtData, error, []string) {
	incorrectFields := signUpData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}

	jwtData, err, fields := adapters.GetStorages().Auth.SignUp(signUpData)
	if err != nil {
		return jwtData, err, fields
	}

	err = db.ProfileCreate(db.ProfileData{
		Id: jwtData.Id,
	})
	return jwtData, err, fields
}
func (*Storage) GetUser(id int64) (models.UserData, error) {
	return db.GetUser(id)
}
func (*Storage) UpdateUser(id int64, updateData models.UpdateUserData) (models.JwtData, error, []string) {
	return adapters.GetStorages().Auth.UpdateAuth(id, updateData)
}
func (*Storage) RemoveUser(id int64, removeData models.RemoveUserData) (error, []string) {
	incorrectFields := removeData.Check()
	if len(incorrectFields) > 0 {
		return models.FieldsError, incorrectFields
	}

	err, fields := adapters.GetStorages().Auth.RemoveAuth(id, removeData)
	if err != nil {
		return err, fields
	}

	err = db.ProfileRemove(id)
	return err, nil
}
