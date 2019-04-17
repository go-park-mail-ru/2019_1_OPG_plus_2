package db

import (
	"log"
	"reflect"
	"testing"

	"2019_1_OPG_plus_2/internal/pkg/models"
)

var auths = []AuthData{
	{
		Username: "username_1",
		Email:    "mail_1@mail.ru",
		PassHash: "pass_1",
	},
	{
		Username: "username_2",
		Email:    "mail_2@mail.ru",
		PassHash: "pass_2",
	},
	{
		Username: "username_3",
		Email:    "mail_3@mail.ru",
		PassHash: "pass_3",
	},
}

func init() {
	testsInitial()
	if err := AuthTruncate(); err != nil {
		log.Fatal(err)
	}
}

func TestAuthCreate(t *testing.T) {
	for i, auth := range auths {
		id, err := AuthCreate(auth)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		auths[i].Id = id
	}
}

func TestAuthCreateAlreadyExists(t *testing.T) {
	for _, auth := range auths {
		_, err := AuthCreate(auth)
		if err != models.AlreadyExists {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.AlreadyExists)
		}
	}
}

func TestAuthFindByEmailAndPassHash(t *testing.T) {
	for _, auth := range auths {
		data, err := AuthFindByEmailAndPassHash(auth.Email, auth.PassHash)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if !reflect.DeepEqual(auth, data) {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", data, auth)
		}
	}
}

func TestAuthFindByEmailAndPassHashIncorrectEmail(t *testing.T) {
	for _, auth := range auths {
		_, err := AuthFindByEmailAndPassHash(auth.Email+"salt", auth.PassHash)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}

func TestAuthFindByEmailAndPassHashIncorrectPassHash(t *testing.T) {
	for _, auth := range auths {
		_, err := AuthFindByEmailAndPassHash(auth.Email, auth.PassHash+"salt")
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}

func TestAuthFindByNicknameAndPassHash(t *testing.T) {
	for _, auth := range auths {
		data, err := AuthFindByUsernameAndPassHash(auth.Username, auth.PassHash)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if !reflect.DeepEqual(auth, data) {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", data, auth)
		}
	}
}

func TestAuthFindByUsernameAndPassHashIncorrectUsername(t *testing.T) {
	for _, auth := range auths {
		_, err := AuthFindByUsernameAndPassHash(auth.Username+"salt", auth.PassHash)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}

func TestAuthFindByUsernameAndPassHashIncorrectPassHash(t *testing.T) {
	for _, auth := range auths {
		_, err := AuthFindByUsernameAndPassHash(auth.Username, auth.PassHash+"salt")
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}

func TestAuthUpdateData(t *testing.T) {
	for i, auth := range auths {
		auth.Username += "_new"
		auth.Email += "_new"
		err := AuthUpdateData(auth)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		auths[i] = auth
	}

	// Test after updating
	TestAuthFindByEmailAndPassHash(t)
	TestAuthFindByNicknameAndPassHash(t)
}

func TestAuthUpdateAlreadyExists(t *testing.T) {
	firstId := auths[0].Id
	for i, auth := range auths {
		if i == 0 {
			continue
		}

		auth.Id = firstId
		err := AuthUpdateData(auth)
		if err != models.AlreadyExists {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.AlreadyExists)
		}
	}

	// Test after updating
	TestAuthFindByEmailAndPassHash(t)
	TestAuthFindByNicknameAndPassHash(t)
}

func TestAuthUpdateIncorrectId(t *testing.T) {
	auth := auths[0]
	auth.Id = 0
	err := AuthUpdateData(auth)
	if err != models.NotFound {
		t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
	}

	// Test after updating
	TestAuthFindByEmailAndPassHash(t)
	TestAuthFindByNicknameAndPassHash(t)
}

func TestAuthUpdatePassword(t *testing.T) {
	for i, auth := range auths {
		auth.PassHash += "_new"
		err := AuthUpdatePassword(auth.Id, auth.PassHash)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		auths[i] = auth
	}

	// Test after updating
	TestAuthFindByEmailAndPassHash(t)
	TestAuthFindByNicknameAndPassHash(t)
}

func TestAuthUpdatePasswordIncorrectId(t *testing.T) {
	err := AuthUpdatePassword(0, "new_pass_hash")
	if err != models.NotFound {
		t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
	}

	// Test after updating
	TestAuthFindByEmailAndPassHash(t)
	TestAuthFindByNicknameAndPassHash(t)
}

func TestAuthRemoveIncorrectId(t *testing.T) {
	for _, auth := range auths {
		err := AuthRemove(0, auth.PassHash)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}

	// Test after updating
	TestAuthFindByEmailAndPassHash(t)
	TestAuthFindByNicknameAndPassHash(t)
}

func TestAuthRemove(t *testing.T) {
	for _, auth := range auths {
		err := AuthRemove(auth.Id, auth.PassHash)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
		}
	}
}

func TestAuthRemoveAlreadyRemoved(t *testing.T) {
	for _, auth := range auths {
		err := AuthRemove(auth.Id, auth.PassHash)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}
