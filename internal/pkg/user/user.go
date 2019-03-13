package user

import (
    "database/sql"
    "fmt"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

func CreateUser(signUpData models.SingUpData) (jwtData models.JwtData, err error) {
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

func GetUser(id int64) (userData models.UserData, err error) {
    userData, err = db.GetUser(id)
    if err == sql.ErrNoRows {
        return userData, fmt.Errorf("user not found")
    }
    return
}

func UpdateUser(id int64, userData models.UpdateUserData) (jwtData models.JwtData, err error) {
    jwtData, err = auth.UpdateAuth(id, userData)
    //if err != nil {
    //   return
    //}

    //err = db.ProfileUpdateData(db.ProfileData{})
    return
}

func RemoveUser(id int64, removeData models.RemoveUserData) error {
    err := auth.RemoveAuth(id, removeData)
    if err != nil {
        return err
    }
    err = db.ProfileRemove(id)
    return err
}
