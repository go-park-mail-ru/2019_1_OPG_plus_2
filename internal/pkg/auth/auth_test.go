package auth

import (
	"2019_1_OPG_plus_2/internal/pkg/config"
	"log"
	"testing"

	a "2019_1_OPG_plus_2/internal/pkg/adapters"
	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/models"
	"2019_1_OPG_plus_2/internal/pkg/user"
)

var authData = []db.AuthData{
	{
		Username: "username_1",
		Email:    "mail_1@mail.ru",
		Password: "pass_1",
	},
	{
		Username: "username_2",
		Email:    "mail_2@mail.ru",
		Password: "pass_2",
	},
	{
		Username: "username_3",
		Email:    "mail_3@mail.ru",
		Password: "pass_3",
	},
}

func init() {
	// Базы для тестов
	db.AuthDbName = config.Db.AuthTestDb
	db.CoreDbName = config.Db.CoreTestDb
	a.SetStorages(user.NewStorage(), NewStorage())

	if err := db.Open(); err != nil && err != db.AlreadyInit {
		log.Fatal(err.Error())
	}

	if err := db.AuthTruncate(); err != nil {
		log.Fatal(err)
	}
}

func TestSignUp(t *testing.T) {
	for i, data := range authData {
		jwt, err, fields := a.GetStorages().Auth.SignUp(models.SignUpData{
			Username: data.Username,
			Email:    data.Email,
			Password: data.Password,
		})
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if len(fields) != 0 {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", fields, []string{})
			continue
		}
		if jwt.Username != data.Username || jwt.Email != data.Email {
			t.Errorf("Wrong Data:\n\tGot: %v, %v\n\tExpected: %v, %v\n", jwt.Username, jwt.Email, data.Username, data.Email)
			continue
		}
		authData[i].Id = jwt.Id
	}
}

func TestSignInByUsername(t *testing.T) {
	for _, data := range authData {
		jwt, err, fields := a.GetStorages().Auth.SignIn(models.SignInData{
			Login:    data.Username,
			Password: data.Password,
		})
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if len(fields) != 0 {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", fields, []string{})
			continue
		}
		if jwt.Username != data.Username || jwt.Email != data.Email {
			t.Errorf("Wrong Data:\n\tGot: %v, %v\n\tExpected: %v, %v\n", jwt.Username, jwt.Email, data.Username, data.Email)
			continue
		}
	}
}

func TestSignInByEmail(t *testing.T) {
	for _, data := range authData {
		jwt, err, fields := a.GetStorages().Auth.SignIn(models.SignInData{
			Login:    data.Email,
			Password: data.Password,
		})
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if len(fields) != 0 {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", fields, []string{})
			continue
		}
		if jwt.Username != data.Username || jwt.Email != data.Email {
			t.Errorf("Wrong Data:\n\tGot: %v, %v\n\tExpected: %v, %v\n", jwt.Username, jwt.Email, data.Username, data.Email)
			continue
		}
	}
}

func TestUpdateAuth(t *testing.T) {
	for i, data := range authData {
		data.Username += "_name"
		data.Email += "s"
		jwt, err, fields := a.GetStorages().Auth.UpdateAuth(data.Id, models.UpdateUserData{
			Username: data.Username,
			Email:    data.Email,
		})

		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if len(fields) != 0 {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", fields, []string{})
			continue
		}
		if jwt.Username != data.Username || jwt.Email != data.Email {
			t.Errorf("Wrong Data:\n\tGot: %v, %v\n\tExpected: %v, %v\n", jwt.Username, jwt.Email, data.Username, data.Email)
			continue
		}

		authData[i] = data
	}
}

func TestUpdatePassword(t *testing.T) {
	for i, data := range authData {
		data.Password += "_pass"
		err, fields := a.GetStorages().Auth.UpdatePassword(data.Id, models.UpdatePasswordData{
			NewPassword:     data.Password,
			PasswordConfirm: data.Password,
		})

		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if len(fields) != 0 {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", fields, []string{})
			continue
		}
		authData[i] = data
	}

	// Test login after password change
	TestSignInByEmail(t)
	TestSignInByUsername(t)
}

func TestRemoveAuth(t *testing.T) {
	for _, data := range authData {
		err, fields := a.GetStorages().Auth.RemoveAuth(data.Id, models.RemoveUserData{
			Password: data.Password,
		})

		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if len(fields) != 0 {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", fields, []string{})
			continue
		}
	}
}

func TestRemoveAuthAlreadyRemoved(t *testing.T) {
	for _, data := range authData {
		err, fields := a.GetStorages().Auth.RemoveAuth(data.Id, models.RemoveUserData{
			Password: data.Password,
		})

		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
			continue
		}
		if len(fields) != 0 {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", fields, []string{})
			continue
		}
	}
}
