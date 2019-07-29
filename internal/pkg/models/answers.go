package models

import "fmt"

// TODO
//go generate: easyjson -all

// 100-199 - success answer
// 200-299 - incorrect answer because of the user
// 300-399 - incorrect answer because developers

type MessageAnswer struct {
	Status  int    `json:"status, int" example:"100"`
	Message string `json:"message, string" example:"ok"`
}

type IncorrectFieldsAnswer struct {
	Status  int      `json:"status, int" example:"204"`
	Message string   `json:"message, string" example:"incorrect fields"`
	Data    []string `json:"data" example:"[email, username]"`
}

type UserDataAnswer struct {
	Status  int      `json:"status, int" example:"105"`
	Message string   `json:"message, string" example:"user found"`
	Data    UserData `json:"data"`
}

type ScoreboardAnswer struct {
	Status  int            `json:"status, int" example:"106"`
	Message string         `json:"message, string" example:"users found"`
	Data    ScoreboardData `json:"data"`
}

type UploadAvatarAnswer struct {
	Status  int    `json:"status, int" example:"108"`
	Message string `json:"message, string" example:"avatar uploaded"`
	Data    string `json:"data" example:"/upload/1.jpg"`
}

/* SUCCESS ANSWERS */

func GetSuccessAnswer(message string) *MessageAnswer {
	return &MessageAnswer{
		Status:  100,
		Message: message,
	}
}

var SignedInAnswer = MessageAnswer{
	Status:  101,
	Message: "signed in",
}

var SignedOutAnswer = MessageAnswer{
	Status:  102,
	Message: "signed out",
}

var SignedUpAnswer = MessageAnswer{
	Status:  103,
	Message: "signed up",
}

var PasswordUpdatedAnswer = MessageAnswer{
	Status:  104,
	Message: "password updated",
}

func GetUserDataAnswer(data UserData) *UserDataAnswer {
	return &UserDataAnswer{
		Status:  105,
		Message: "user found",
		Data:    data,
	}
}

func GetScoreboardAnswer(data ScoreboardData) *ScoreboardAnswer {
	return &ScoreboardAnswer{
		Status:  106,
		Message: "users found",
		Data:    data,
	}
}

var UserUpdatedAnswer = MessageAnswer{
	Status:  107,
	Message: "user updated",
}

func GetUploadAvatarAnswer(url string) *UploadAvatarAnswer {
	return &UploadAvatarAnswer{
		Status:  108,
		Message: "avatar uploaded",
		Data:    url,
	}
}

var UserRemovedAnswer = MessageAnswer{
	Status:  109,
	Message: "user removed",
}

/* USERS ERRORS */

// For future use
//
// func GetUserErrorAnswer(error string) *MessageAnswer {
// 	return &MessageAnswer{
// 		Status:  200,
// 		ChatMessage: error,
// 	}
// }

var NotSignedInAnswer = MessageAnswer{
	Status:  201,
	Message: "not signed in",
}

var AlreadySignedInAnswer = MessageAnswer{
	Status:  202,
	Message: "already signed in",
}

var AlreadySignedOutAnswer = MessageAnswer{
	Status:  203,
	Message: "already signed out",
}

var UserNotFoundAnswer = MessageAnswer{
	Status:  205,
	Message: "user not found",
}

func GetIncorrectFieldsAnswer(data []string) *IncorrectFieldsAnswer {
	return &IncorrectFieldsAnswer{
		Status:  204,
		Message: "incorrect fields",
		Data:    data,
	}
}

var FileTooBigAnswer = MessageAnswer{
	Status:  206,
	Message: "file is too big",
}

var NotImageAnswer = MessageAnswer{
	Status:  207,
	Message: "not image",
}

/* DEVELOPERS ERRORS */

func GetDeveloperErrorAnswer(error string) *MessageAnswer {
	return &MessageAnswer{
		Status:  300,
		Message: error,
	}
}

var IncorrectJsonAnswer = MessageAnswer{
	Status:  301,
	Message: "incorrect JSON",
}

var ImpossibleReadFileAnswer = MessageAnswer{
	Status:  302,
	Message: "impossible to read file",
}

var ImpossibleSaveFileAnswer = MessageAnswer{
	Status:  303,
	Message: "impossible save file",
}

var IncorrectQueryParams = MessageAnswer{
	Status:  401,
	Message: "incorrect query params",
}

func GetNotFoundRoomAnswer(id string) MessageAnswer {
	return MessageAnswer{
		Status:  402,
		Message: fmt.Sprintf("Room %v not found", id),
	}
}
