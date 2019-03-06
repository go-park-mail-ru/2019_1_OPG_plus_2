package cache_dummies

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

type Cache interface {
	Get(sessionToken string) (*models.UserSessionRecord, error)
	Set(sessionToken string, dummy *models.UserSessionRecord) error
	Delete(sessionToken string) error
}

type CookieCacheDummy struct {
	Data map[string]*models.UserSessionRecord
}

func NewCookieCacheDummy() *CookieCacheDummy {
	return &CookieCacheDummy{Data: make(map[string]*models.UserSessionRecord)}
}

func (cache *CookieCacheDummy) Set(sessionToken string, dummy *models.UserSessionRecord) error {
	cache.Data[sessionToken] = dummy
	return nil
}

func (cache *CookieCacheDummy) Get(sessionToken string) (*models.UserSessionRecord, error) {
	result := cache.Data[sessionToken]
	//if result == nil {
	//	return nil, fmt.Errorf("NO COOKIE IN CACHE")
	//}
	return result, nil
}

func (cache *CookieCacheDummy) Delete(sessionToken string) error {
	delete(cache.Data, sessionToken)
	return nil
}

type UserStorage struct {
	Data map[string]*models.UserData
}

func NewUserStorage() *UserStorage {
	return &UserStorage{Data: make(map[string]*models.UserData)}
}

func (storage *UserStorage) Get(username string) (*models.UserData, error) {
	result := storage.Data[username]
	if result == nil {
		return nil, fmt.Errorf("NO USER IN CACHE")
	}
	return result, nil
}

func (storage *UserStorage) Set(username string, user *models.UserData) error {
	if storage.Data[username] != nil {
		return fmt.Errorf("USER ALREADY EXISTS")
	}
	storage.Data[username] = user
	return nil
}

func (storage *UserStorage) Delete(username string) error {
	delete(storage.Data, username)
	return nil
}
