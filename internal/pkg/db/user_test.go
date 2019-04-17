package db

import (
	"log"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"2019_1_OPG_plus_2/internal/pkg/models"
)

var users = []models.UserData{
	{
		Username: "username_1",
		Email:    "mail_1@mail.ru",
		Avatar:   "avatar_1",
	},
	{
		Username: "username_2",
		Email:    "mail_2@mail.ru",
		Avatar:   "avatar_2",
	},
	{
		Username: "username_3",
		Email:    "mail_3@mail.ru",
		Avatar:   "avatar_3",
	},
}

func init() {
	testsInitial()
	if err := AuthTruncate(); err != nil {
		log.Fatal(err)
	}
	if err := ProfileTruncate(); err != nil {
		log.Fatal(err)
	}
}

func TestUsersCreate(t *testing.T) {
	for i, user := range users {
		id, err := AuthCreate(AuthData{
			Username: user.Username,
			Email:    user.Email,
		})
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		user.Id = id

		err = ProfileCreate(ProfileData{
			Id:     user.Id,
			Avatar: user.Avatar,
		})
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}

		users[i] = user
	}
}

func TestUsersUpdate(t *testing.T) {
	for i, user := range users {
		user.Games += int64(rand.Intn(100))
		user.Win = int64(rand.Intn(int(user.Games)))
		user.Lose = user.Games - user.Win
		user.Score = user.Win * 100
		err := ProfileUpdateData(ProfileData{
			Id:    user.Id,
			Score: user.Score,
			Games: user.Games,
			Win:   user.Win,
			Lose:  user.Lose,
		})
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		users[i] = user
	}
}

func TestGetUser(t *testing.T) {
	for _, user := range users {
		data, err := GetUser(user.Id)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if !reflect.DeepEqual(data, user) {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", data, user)
		}
	}
}

func TestTestGetUserIncorrectId(t *testing.T) {
	_, err := GetUser(0)
	if err != models.NotFound {
		t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
	}
}

func TestGetScoreboard(t *testing.T) {
	limits := []int64{1, 1, 1, 1, 2, 2, 2, 3, 3, 10, 10, 100}
	pages := []int64{0, 1, 2, 3, 0, 1, 2, 0, 1, 0, 1, 100}

	if len(limits) != len(pages) {
		t.Fatal("Length of limits and pages are different!")
	}

	// Sort users
	sort.Slice(users, func(i, j int) bool {
		return users[i].Score > users[j].Score ||
			users[i].Score == users[j].Score && users[i].Win > users[j].Win ||
			users[i].Score == users[j].Score && users[i].Win == users[j].Win && users[i].Id < users[j].Id
	})

	for i := range limits {
		data, count, err := GetScoreboard(limits[i], limits[i]*pages[i])
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}

		if count != uint64(len(users)) {
			t.Errorf("Wrong Count:\n\tGot: %v\n\tExpected: %v\n", count, len(users))
		}

		var scoreboard []models.ScoreboardUserData
		for u, user := range users {
			if int64(u) >= limits[i]*pages[i] && int64(u) < limits[i]*(pages[i]+1) {
				scoreboard = append(scoreboard, models.ScoreboardUserData{
					Id:       user.Id,
					Username: user.Username,
					Avatar:   user.Avatar,
					Score:    user.Score,
				})
			}
		}
		sort.Slice(scoreboard, func(i, j int) bool { return scoreboard[i].Score > scoreboard[j].Score })

		if !reflect.DeepEqual(data, scoreboard) {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n\tLimit: %v, Offset: %v\n", data, scoreboard, limits[i], pages[i])
		}
	}
}
