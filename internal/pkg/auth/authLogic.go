package auth

import (
	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/models"
)

func SignUp(signUpData models.SignUpData) (models.JwtData, error, []string) {
	incorrectFields := signUpData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}

	id, err := db.AuthGetIdByEmail(signUpData.Email)
	if err != nil {
		return models.JwtData{}, err, nil
	}
	if id != 0 {
		incorrectFields = append(incorrectFields, "email")
	}
	id, err = db.AuthGetIdByUsername(signUpData.Username)
	if err != nil {
		return models.JwtData{}, err, nil
	}
	if id != 0 {
		incorrectFields = append(incorrectFields, "username")
	}
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.AlreadyExists, incorrectFields
	}

	id, err = db.AuthCgitreate(db.AuthData{
		Email:    signUpData.Email,
		Username: signUpData.Username,
		Password: PasswordHash(signUpData.Password),
	})
	if err != nil {
		return models.JwtData{}, err, nil
	}

	return models.JwtData{
		Id:       id,
		Email:    signUpData.Email,
		Username: signUpData.Username,
	}, nil, nil
}

func SignIn(signInData models.SignInData) (data models.JwtData, err error, incorrectFields []string) {
	var userData db.AuthData
	passHash := PasswordHash(signInData.Password)

	isEmail := models.CheckEmail(signInData.Login)
	if isEmail {
		userData, err = db.AuthFindByEmailAndPassHash(signInData.Login, passHash)
		if err != nil {
			if err == models.NotFound {
				return data, models.FieldsError, append(incorrectFields, "password")
			}
			return
		}
	}

	isUsername := !isEmail && models.CheckUsername(signInData.Login)
	if isUsername {
		userData, err = db.AuthFindByUsernameAndPassHash(signInData.Login, passHash)
		if err != nil {
			if err == models.NotFound {
				return data, models.FieldsError, append(incorrectFields, "password")
			}
			return
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

func UpdateAuth(id int64, userData models.UpdateUserData) (models.JwtData, error, []string) {
	incorrectFields := userData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}

	err := db.AuthUpdateData(db.AuthData{
		Id:       id,
		Email:    userData.Email,
		Username: userData.Username,
	})
	if err != nil {
		return models.JwtData{}, err, nil
	}

	return models.JwtData{
		Id:       id,
		Email:    userData.Email,
		Username: userData.Username,
	}, nil, nil
}

func UpdatePassword(id int64, passwordData models.UpdatePasswordData) (error, []string) {
	incorrectFields := passwordData.Check()
	if len(incorrectFields) > 0 {
		return models.FieldsError, incorrectFields
	}

	return db.AuthUpdatePassword(id, PasswordHash(passwordData.NewPassword)), nil
}

func RemoveAuth(id int64, removeData models.RemoveUserData) (error, []string) {
	incorrectFields := removeData.Check()
	if len(incorrectFields) > 0 {
		return models.FieldsError, incorrectFields
	}

	err := db.AuthRemove(id, PasswordHash(removeData.Password))
	if err != nil {
		if err.Error() == "incorrect password" {
			return models.FieldsError, append(incorrectFields, "password")
		}
		return err, nil
	}
	return nil, nil
}
