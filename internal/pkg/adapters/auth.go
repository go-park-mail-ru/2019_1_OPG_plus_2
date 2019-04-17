package adapters

import (
	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/models"
)

type mockAuthStorage struct {
	users map[int64]*db.AuthData
}

func NewMockAuthStorage() (storage *mockAuthStorage) {
	storage = new(mockAuthStorage)
	storage.users = make(map[int64]*db.AuthData)
	storage.users[1] = &db.AuthData{
		Email:    "mail1@mail.ru",
		Username: "username1",
		Id:       1,
		Password: "pass1",
	}
	storage.users[2] = &db.AuthData{
		Email:    "mail2@mail.ru",
		Username: "username2",
		Id:       2,
		Password: "pass2",
	}
	storage.users[3] = &db.AuthData{
		Email:    "mail3@mail.ru",
		Username: "username3",
		Id:       3,
		Password: "pass3",
	}
	return storage
}

func (storage mockAuthStorage) SignUp(signUpData models.SignUpData) (models.JwtData, error, []string) {
	incorrectFields := signUpData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}

	id := int64(1)
	for _, user := range storage.users {
		if user.Id >= id {
			id = user.Id + 1
		}
		if user.Email == signUpData.Email || user.Username == signUpData.Username {
			return models.JwtData{}, models.AlreadyExists, nil
		}
	}

	storage.users[id] = &db.AuthData{
		Email:    signUpData.Email,
		Username: signUpData.Username,
		Password: signUpData.Password,
	}

	return models.JwtData{
		Id:       id,
		Email:    signUpData.Email,
		Username: signUpData.Username,
	}, nil, nil
}

func (storage mockAuthStorage) SignIn(signInData models.SignInData) (data models.JwtData, err error, incorrectFields []string) {
	var userData db.AuthData
	passHash := signInData.Password

	isEmail := models.CheckEmail(signInData.Login)
	if isEmail {
		isExists := false
		for _, user := range storage.users {
			if user.Email == signInData.Login && user.Password == passHash {
				userData = *user
				isExists = true
				break
			}
		}
		if !isExists {
			return data, models.FieldsError, append(incorrectFields, "password")
		}
	}

	isUsername := !isEmail && models.CheckUsername(signInData.Login)
	if isUsername {
		isExists := false
		for _, user := range storage.users {
			if user.Username == signInData.Login && user.Password == passHash {
				userData = *user
				isExists = true
				break
			}
		}
		if !isExists {
			return data, models.FieldsError, append(incorrectFields, "password")
		}
	}

	if !isEmail && !isUsername {
		return data, models.FieldsError, append(incorrectFields, "login")
	}

	return models.JwtData{
		Id:       userData.Id,
		Email:    userData.Email,
		Username: userData.Username,
	}, nil, nil
}

func (storage mockAuthStorage) UpdatePassword(id int64, passwordData models.UpdatePasswordData) (error, []string) {
	incorrectFields := passwordData.Check()
	if len(incorrectFields) > 0 {
		return models.FieldsError, incorrectFields
	}

	if storage.users[id] == nil {
		return models.NotFound, nil
	}

	storage.users[id].Password = passwordData.NewPassword
	return nil, nil
}

func (storage mockAuthStorage) UpdateAuth(id int64, updateData models.UpdateUserData) (models.JwtData, error, []string) {
	incorrectFields := updateData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}

	if storage.users[id] == nil {
		return models.JwtData{}, models.NotFound, nil
	}

	user := storage.users[id]
	user.Username = updateData.Username
	user.Email = updateData.Email

	newJwtData := models.JwtData{
		Id:       id,
		Username: user.Username,
		Email:    user.Email,
	}
	return newJwtData, nil, nil
}

func (storage mockAuthStorage) RemoveAuth(id int64, removeData models.RemoveUserData) (error, []string) {
	incorrectFields := removeData.Check()
	if len(incorrectFields) > 0 {
		return models.FieldsError, incorrectFields
	}

	if storage.users[id] == nil {
		return models.NotFound, nil
	}

	if storage.users[id].Password != removeData.Password {
		return models.FieldsError, append(incorrectFields, "password")
	}

	delete(storage.users, id)
	return nil, nil
}
