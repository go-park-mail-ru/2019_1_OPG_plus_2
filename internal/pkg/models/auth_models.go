package models

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserSessionRecord struct {
	Username     string `json:"username, string"`
	SessionToken string `json:"session_token, string"`
	Timeout      int    `json:"timeout, int"` // seconds
}

type UserData struct {
	Username string `json:"username, string"`
	Password string `json:"password, string"`
	EMail    string `json:"email, string, omitempty"`
}

type SuccessOrErrorMessage struct {
	Status  int    `json:"status, int"`
	Message string `json:"message, string"`
}
