package adapters

import (
	"errors"

	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

type mockCoreStorage struct {
	users map[int64]*db.ProfileData
}

func NewMockCoreStorage() (storage *mockCoreStorage) {
	storage = new(mockCoreStorage)
	storage.users = make(map[int64]*db.ProfileData)
	storage.users[1] = &db.ProfileData{
		Id:     1,
		Score:  1000,
		Avatar: "avatar1",
		Games:  100,
		Lose:   50,
		Win:    50,
	}
	storage.users[2] = &db.ProfileData{
		Id:     2,
		Score:  2000,
		Avatar: "avatar2",
		Games:  200,
		Lose:   100,
		Win:    100,
	}
	storage.users[3] = &db.ProfileData{
		Id:     3,
		Score:  3000,
		Avatar: "avatar3",
		Games:  300,
		Lose:   150,
		Win:    150,
	}
	return storage
}

func (storage *mockCoreStorage) CreateUser(signUpData models.SignUpData) (jwtData models.JwtData, err error, fields []string) {
	incorrectFields := signUpData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}

	jwtData, err, fields = GetStorages().Auth.SignUp(signUpData)
	if err != nil {
		return
	}

	if storage.users[jwtData.Id] != nil {
		return models.JwtData{}, models.AlreadyExists, nil
	}

	newUser := db.ProfileData{
		Id: jwtData.Id,
	}
	storage.users[newUser.Id] = &newUser

	return
}

func (storage *mockCoreStorage) GetUser(id int64) (userData models.UserData, err error) {
	authStorage, ok := GetStorages().Auth.(*mockAuthStorage)
	if !ok {
		return userData, errors.New("can not convert storage to mock storage")
	}

	if authStorage.users[id] == nil || storage.users[id] == nil {
		return models.UserData{}, models.NotFound
	}
	return models.UserData{
		Id:       storage.users[id].Id,
		Username: authStorage.users[id].Username,
		Email:    authStorage.users[id].Email,
		Avatar:   storage.users[id].Avatar,
		Score:    storage.users[id].Score,
		Games:    storage.users[id].Games,
		Win:      storage.users[id].Win,
		Lose:     storage.users[id].Lose,
	}, nil
}

func (storage *mockCoreStorage) UpdateUser(id int64, updateData models.UpdateUserData) (models.JwtData, error, []string) {
	return GetStorages().Auth.UpdateAuth(id, updateData)
}

func (storage *mockCoreStorage) RemoveUser(id int64, removeData models.RemoveUserData) (error, []string) {
	incorrectFields := removeData.Check()
	if len(incorrectFields) > 0 {
		return models.FieldsError, incorrectFields
	}

	err, fields := GetStorages().Auth.RemoveAuth(id, removeData)
	if err != nil {
		return err, fields
	}

	if storage.users[id] == nil {
		return models.NotFound, nil
	}
	delete(storage.users, id)

	return nil, nil
}
