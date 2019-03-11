package models

import "fmt"

type UserProfile struct {
	ID        int    `json:"id, string" example:"1"`
	Username  string `json:"username, string" example:"user_test"`
	Email     string `json:"email, string" example:"user_test@test.com"`
	AvatarUrl string `json:"avatar_url, string" example:"<some avatar url>"`

	Score       int `json:"score, number"`
	GamesPlayed int `json:"games_played, number"`
	Win         int `json:"win, number"`
	Lose        int `json:"lose, number"`
}

type UserProfileStorage struct {
	Data map[int]*UserProfile
}

func (storage *UserProfileStorage) Get(key int) (object *UserProfile, err error) {
	return storage.Data[key], nil
}

func (storage *UserProfileStorage) Set(key int, object *UserProfile) (err error) {
	storage.Data[key] = object
	return nil
}

func (storage *UserProfileStorage) Delete(key int) (err error) {
	if storage.Data[key] == nil {
		return fmt.Errorf("NO USER IN DB")
	}
	delete(storage.Data, key)
	return nil
}

func NewUserProfileStorage() *UserProfileStorage {
	return &UserProfileStorage{Data: make(map[int]*UserProfile)}
}
