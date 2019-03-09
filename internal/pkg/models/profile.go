package models

type UserProfile struct {
	Username  string `json:"username, string"`
	Email     string `json:"email, string"`
	AvatarUrl string `json:"avatar_url, string"`

	Score       int `json:"score, number"`
	GamesPlayed int `json:"games_played, number"`
	Win         int `json:"win, number"`
	Lose        int `json:"lose, number"`
}

type UserUpdateInfo struct {
	Username string `json:"username, number"`
	Email    string `json:"email, number"`
	Password string `json:"password, number"` //поле под вопросом, скорее всего надо будет скармливать системе аутентификации
}

type UserProfileStorage struct {
	Data map[string]*UserProfile
}

func (storage *UserProfileStorage) Get(key string) (object *UserProfile, err error) {
	return storage.Data[key], nil
}

func (storage *UserProfileStorage) Set(key string, object *UserProfile) (err error) {
	storage.Data[key] = object
	return nil
}

func (storage *UserProfileStorage) Delete(key string) (err error) {
	delete(storage.Data, key)
	return nil
}

func NewUserProfileStorage() *UserProfileStorage {
	return &UserProfileStorage{Data: make(map[string]*UserProfile)}
}
