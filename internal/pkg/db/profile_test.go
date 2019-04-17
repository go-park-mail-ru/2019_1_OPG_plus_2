package db

import (
	"log"
	"math/rand"
	"reflect"
	"testing"

	"2019_1_OPG_plus_2/internal/pkg/models"
)

var profiles = []ProfileData{
	{
		Id:     1,
		Avatar: "avatar_1",
		// Score:  100,
		// Games:  10,
		// Win:    6,
		// Lose:   4,
	},
	{
		Id:     2,
		Avatar: "avatar_2",
		// Score:  200,
		// Games:  20,
		// Win:    12,
		// Lose:   8,
	},
	{
		Id:     3,
		Avatar: "avatar_3",
		// Score:  300,
		// Games:  30,
		// Win:    18,
		// Lose:   12,
	},
}

func init() {
	testsInitial()
	if err := ProfileTruncate(); err != nil {
		log.Fatal(err)
	}
}

func TestProfileCreate(t *testing.T) {
	for _, profile := range profiles {
		err := ProfileCreate(profile)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
		}
	}
}

func TestProfileCreateAlreadyExists(t *testing.T) {
	for _, profile := range profiles {
		err := ProfileCreate(profile)
		if err != models.AlreadyExists {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.AlreadyExists)
		}
	}
}

func TestProfileFindById(t *testing.T) {
	for _, profile := range profiles {
		data, err := ProfileFindById(profile.Id)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		if !reflect.DeepEqual(profile, data) {
			t.Errorf("Wrong Data:\n\tGot: %v\n\tExpected: %v\n", data, profile)
		}
	}
}

func TestProfileFindByIdIncorrectId(t *testing.T) {
	_, err := ProfileFindById(0)
	if err != models.NotFound {
		t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
	}
}

func TestProfileUpdateData(t *testing.T) {
	for i, profile := range profiles {
		profile.Games += int64(rand.Intn(100))
		profile.Win = int64(rand.Intn(int(profile.Games)))
		profile.Lose = profile.Games - profile.Win
		profile.Score = profile.Win * 100
		err := ProfileUpdateData(profile)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		profiles[i] = profile
	}

	// Test after updating
	TestProfileFindById(t)
}

func TestProfileUpdateIncorrectId(t *testing.T) {
	profile := profiles[0]
	profile.Id = 0
	err := ProfileUpdateData(profile)
	if err != models.NotFound {
		t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
	}

	// Test after updating
	TestProfileFindById(t)
}

func TestProfileUpdateAvatar(t *testing.T) {
	for i, profile := range profiles {
		profile.Avatar += "_new"
		err := ProfileUpdateAvatar(profile.Id, profile.Avatar)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
			continue
		}
		profiles[i] = profile
	}

	// Test after updating
	TestProfileFindById(t)
}

func TestProfileUpdateAvatarIncorrectId(t *testing.T) {
	err := ProfileUpdateAvatar(0, "new_avatar")
	if err != models.NotFound {
		t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
	}
	
	// Test after updating
	TestProfileFindById(t)
}

func TestProfileRemoveIncorrectId(t *testing.T) {
	err := ProfileRemove(0)
	if err != models.NotFound {
		t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
	}

	// Test after updating
	TestProfileFindById(t)
}

func TestProfileRemove(t *testing.T) {
	for _, profile := range profiles {
		err := ProfileRemove(profile.Id)
		if err != nil {
			t.Errorf("Unknown Error: %v", err)
		}
	}
}

func TestProfileRemoveAlreadyRemoved(t *testing.T) {
	for _, profile := range profiles {
		err := ProfileRemove(profile.Id)
		if err != models.NotFound {
			t.Errorf("Wrong Error:\n\tGot: %v\n\tExpected: %v\n", err, models.NotFound)
		}
	}
}
