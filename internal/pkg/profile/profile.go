package profile

import (
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

func CreateProfile(signUpData models.SingUpData) (jwtData models.JwtData, err error) {
    jwtData, err = auth.SignUp(signUpData)
    if err != nil {
        return
    }

    err = db.ProfileCreate(db.ProfileData{
        Id: jwtData.Id,
        Avatar: signUpData.Avatar,
    })
    return
}

func GetProfile(id int64) (userData models.UserData, err error) {
    // TODO: Exchange to join query
    authData, err := db.AuthFindById(id)
    if err != nil {
        return
    }
    profileData, err := db.ProfileFindById(id)
    if err != nil {
        return
    }
    return models.UserData{
        Id: authData.Id,
        Email: authData.Email,
        Username: authData.Username,
        Avatar: profileData.Avatar,
        Score: profileData.Score,
        Games: profileData.Games,
        Win: profileData.Win,
        Lose: profileData.Lose,
    }, nil
}

func GetProfilesOrderedByScore() (userData []models.UserData, err error) {
    return
}

func UpdateProfile(userData models.UserData) error {
    return nil
}

func UpdateAvatar(id int64, avatar string) error {
    return nil
}

func RemoveProfile(id int64, removeData models.RemoveUserData) error {
    return nil
}
