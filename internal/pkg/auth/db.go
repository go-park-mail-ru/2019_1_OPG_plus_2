package auth

import (
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"sync"
)

var users = []*models.DbUserData{
	{1, "test@mail.ru", "test", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
	{2, "user@mail.ru", "user", "04f8996da763b7a969b1028ee3007569eaf3a635486ddab211d512c85b9df8fb"},
}
var usersMutex sync.Mutex

func create(data models.DbUserData) int {
	if findByNickname(data.Username) != nil {
		return 0
	}
	if findByEmail(data.Email) != nil {
		return 0
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()
	data.Id = users[len(users)-1].Id + 1
	users = append(users, &data)
	return data.Id
}

func findById(id int) *models.DbUserData {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	for _, v := range users {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func findByEmail(email string) *models.DbUserData {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	for _, v := range users {
		if v.Email == email {
			return v
		}
	}
	return nil
}

func findByNickname(nickname string) *models.DbUserData {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	for _, v := range users {
		if v.Username == nickname {
			return v
		}
	}
	return nil
}

func updateById(data models.DbUserData) bool {
	oldData := findById(data.Id)
	if oldData == nil {
		return false
	}
	*oldData = data
	return true
}

func removeByEmail(email string, passHash string) bool {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	for i, v := range users {
		if v.Email == email {
			if v.PassHash == passHash {
				copy(users[i:], users[i+1:])
				users[len(users)-1] = nil
				users = users[:len(users)-1]
				usersMutex.Unlock()
				return true
			}
			return false
		}
	}
	return false
}

func removeByNickname(nickname string, passHash string) bool {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	for i, v := range users {
		if v.Username == nickname {
			if v.PassHash == passHash {
				copy(users[i:], users[i+1:])
				users[len(users)-1] = nil
				users = users[:len(users)-1]
				return true
			}
			return false
		}
	}
	return false
}
