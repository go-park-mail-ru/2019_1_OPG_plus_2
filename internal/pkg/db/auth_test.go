package db

import (
	"log"
	"reflect"
	"testing"

	"2019_1_OPG_plus_2/internal/pkg/models"
)

var users = []AuthData{
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

func TestCreate(t *testing.T) {
	for i, user := range users {
		id, err := AuthCreate(user)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		users[i].Id = id
	}
}

func TestCreateAlreadyExists(t *testing.T) {
	for _, user := range users {
		_, err := AuthCreate(user)
		if err != models.AlreadyExists {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.AlreadyExists)
		}
	}
}

func TestFindByEmailAndPassHash(t *testing.T) {
	for _, user := range users {
		data, err := AuthFindByEmailAndPassHash(user.Email, user.PassHash)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if !reflect.DeepEqual(user, data) {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", data, user)
		}
	}
}

func TestFindByEmailAndPassHashIncorrectEmail(t *testing.T) {
	for _, user := range users {
		_, err := AuthFindByEmailAndPassHash(user.Email+"salt", user.PassHash)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}

func TestFindByEmailAndPassHashIncorrectPassHash(t *testing.T) {
	for _, user := range users {
		_, err := AuthFindByEmailAndPassHash(user.Email, user.PassHash+"salt")
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}

func TestFindByNicknameAndPassHash(t *testing.T) {
	for _, user := range users {
		data, err := AuthFindByUsernameAndPassHash(user.Username, user.PassHash)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if !reflect.DeepEqual(user, data) {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", data, user)
		}
	}
}

func TestFindByUsernameAndPassHashIncorrectUsername(t *testing.T) {
	for _, user := range users {
		_, err := AuthFindByUsernameAndPassHash(user.Username+"salt", user.PassHash)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}

func TestFindByUsernameAndPassHashIncorrectPassHash(t *testing.T) {
	for _, user := range users {
		_, err := AuthFindByUsernameAndPassHash(user.Username, user.PassHash+"salt")
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}

func TestUpdateData(t *testing.T) {
	for i, user := range users {
		user.Username += "_new"
		user.Email += "_new"
		err := AuthUpdateData(user)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		users[i] = user
	}

	// Test after updating
	TestFindByEmailAndPassHash(t)
	TestFindByNicknameAndPassHash(t)
}

func TestUpdateAlreadyExists(t *testing.T) {
	firstId := users[0].Id
	for i, user := range users {
		if i == 0 {
			continue
		}

		user.Id = firstId
		err := AuthUpdateData(user)
		if err != models.AlreadyExists {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.AlreadyExists)
		}
	}

	// Test after updating
	TestFindByEmailAndPassHash(t)
	TestFindByNicknameAndPassHash(t)
}

func TestUpdateIncorrectId(t *testing.T) {
	for _, user := range users {
		user.Id = 0
		err := AuthUpdateData(user)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}

	// Test after updating
	TestFindByEmailAndPassHash(t)
	TestFindByNicknameAndPassHash(t)
}

func TestUpdatePassword(t *testing.T) {
	for i, user := range users {
		user.PassHash += "_new"
		err := AuthUpdatePassword(user.Id, user.PassHash)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		users[i] = user
	}

	// Test after updating
	TestFindByEmailAndPassHash(t)
	TestFindByNicknameAndPassHash(t)
}

func TestUpdatePasswordIncorrectId(t *testing.T) {
	for _, user := range users {
		err := AuthUpdatePassword(0, user.PassHash)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}

	// Test after updating
	TestFindByEmailAndPassHash(t)
	TestFindByNicknameAndPassHash(t)
}

func TestRemoveIncorrectId(t *testing.T) {
	for _, user := range users {
		err := AuthRemove(0, user.PassHash)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}

	// Test after updating
	TestFindByEmailAndPassHash(t)
	TestFindByNicknameAndPassHash(t)
}

func TestRemove(t *testing.T) {
	for _, user := range users {
		err := AuthRemove(user.Id, user.PassHash)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
		}
	}
}

func TestRemoveAlreadyRemoved(t *testing.T) {
	for _, user := range users {
		err := AuthRemove(user.Id, user.PassHash)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}
