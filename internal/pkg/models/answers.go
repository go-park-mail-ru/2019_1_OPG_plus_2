package models

/* ANSWER WITH MESSAGE */

type MessageAnswer struct {
	Status  int    `json:"status, int" example:"200"`
	Message string `json:"message, string" example:"ok"`
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

var UserUpdatedAnswer = MessageAnswer{
	Status:  105,
	Message: "user updated",
}

var UserRemovedAnswer = MessageAnswer{
	Status:  106,
	Message: "user removed",
}

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

var IncorrectJsonAnswer = MessageAnswer{
	Status:  301,
	Message: "incorrect JSON",
}

var UserNotFoundAnswer = MessageAnswer{
	Status:  302,
	Message: "user not found",
}

var FileTooBigAnswer = MessageAnswer{
	Status:  303,
	Message: "file is too big",
}

var ImpossibleReadFileAnswer = MessageAnswer{
	Status:  304,
	Message: "impossible to read file",
}

var NotImageAnswer = MessageAnswer{
	Status:  305,
	Message: "not image",
}

var ImpossibleSaveFileAnswer = MessageAnswer{
	Status:  306,
	Message: "impossible save file",
}

/* ANSWER WITH INCORRECT FIELDS */

type IncorrectFieldsAnswer struct {
	Status  int      `json:"status, int" example:"200"`
	Message string   `json:"message, string" example:"Query processed successfully"`
	Data    []string `json:"data" example:"[email, username]"`
}

func NewIncorrectFieldsAnswer(data []string) IncorrectFieldsAnswer {
	return IncorrectFieldsAnswer{
		Status:  401,
		Message: "incorrect fields",
		Data:    data,
	}
}
