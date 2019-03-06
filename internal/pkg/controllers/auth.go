package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/util/cache_dummies"
	"github.com/google/uuid"
	"net/http"
	"time"
)

//Authentification controllers
var defaultTimeout = 300

var sessionCache = cache_dummies.NewCookieCacheDummy()
var userCache = cache_dummies.NewUserStorage()

func init() {
	err := userCache.Set("user1", &models.UserData{
		Username: "user1",
		Password: "password1",
		EMail:    "user1@example.com",
	})
	if err != nil {
		panic(err)
	}

	err = userCache.Set("user2", &models.UserData{
		Username: "user2",
		Password: "password2",
		EMail:    "user2@example.com",
	})
	if err != nil {
		panic(err)
	}
}

//TODO: проверять таймаут куки

// SignIn godoc
// @title Sign-in
// @description Sign-in method
// @accept json
// @produce json
// @param credentials body models.Credentials true "User credentials"
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @failure 500 {object} models.SuccessOrErrorMessage
// @router /sign_in [post]
func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	var message models.SuccessOrErrorMessage
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		message.Status = http.StatusBadRequest
		message.Message = "Json parsing error"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	fmt.Println(creds)
	// Get the expected password from our in memory map
	user, err := userCache.Get(creds.Username)
	fmt.Println(user)

	expectedPassword := user.Password

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if err != nil || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		message.Status = http.StatusUnauthorized
		message.Message = "User credentials are incorrect"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}

	// Create a new random session token
	sessionToken := uuid.New().String()
	fmt.Println(sessionToken)
	// Set the token in the sessionCache, along with the user whom it represents
	// The token has an expiry time of ${defaultTimeout} seconds
	err = sessionCache.Set(sessionToken, &models.UserSessionRecord{
		Username:     creds.Username,
		SessionToken: sessionToken,
		Timeout:      defaultTimeout,
	})
	if err != nil {
		// If there is an error in setting the sessionCache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		message.Status = http.StatusInternalServerError
		message.Message = "Error while saving user session"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of ${defaultTimeout} seconds, the same as the sessionCache
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(time.Duration(defaultTimeout) * time.Second),
	})
	message.Status = 200
	message.Message = "Logged in"
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
}

// SignOut godoc
// @title Sign-out
// @description Sign-out method
// @produce json
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @router /sign_out [post]
func SignOut(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	var message models.SuccessOrErrorMessage
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			message.Status = http.StatusUnauthorized
			message.Message = "User is not logged in to log out)"
			msg, _ := json.Marshal(message)
			_, _ = fmt.Fprintln(w, string(msg))
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		message.Status = http.StatusBadRequest
		message.Message = "Bad request sent"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	err = sessionCache.Delete(c.Value)
	if err != nil {
		panic(err)
	}
	message.Status = 200
	message.Message = "Logged out"
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
	return

}

// Register godoc
// @title Registration
// @description Method to register new user
// @accept json
// @produce json
// @param user_data body models.UserData true "User profile"
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var newUser models.UserData
	var message models.SuccessOrErrorMessage
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message.Status = http.StatusBadRequest
		message.Message = "Error while parsing profile json"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	fmt.Println(newUser)
	err = userCache.Set(newUser.Username, &newUser)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		message.Status = http.StatusConflict
		message.Message = "Username already busy"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	result, err := json.Marshal(newUser)
	_, _ = fmt.Fprintln(w, string(result))
}

// UpdateProfile godoc
// @title Profile update
// @description Method to update user's profile
// @accept json
// @produce json
// @param user_data body models.UserData true "User profile"
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @failure 500 {object} models.SuccessOrErrorMessage
// @router /update [post]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	var message models.SuccessOrErrorMessage
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		message.Status = http.StatusBadRequest
		message.Message = "User is not authorized to refresh session"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	sessionToken := c.Value

	userRecord, err := sessionCache.Get(sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message.Status = http.StatusInternalServerError
		message.Message = "Error while getting session"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	if userRecord == nil {
		w.WriteHeader(http.StatusUnauthorized)
		message.Status = http.StatusUnauthorized
		message.Message = "User is not authorized to refresh session"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}

	var newProfile models.UserData
	err = json.NewDecoder(r.Body).Decode(&newProfile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message.Status = http.StatusBadRequest
		message.Message = "Error while parsing user profile JSON"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	_ = userCache.Delete(userRecord.Username)
	_ = userCache.Set(newProfile.Username, &newProfile)

	newSessionToken := uuid.New().String()
	_ = sessionCache.Delete(sessionToken)
	_ = sessionCache.Set(newSessionToken, &models.UserSessionRecord{
		Username:     newProfile.Username,
		SessionToken: newSessionToken,
		Timeout:      defaultTimeout,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(time.Duration(defaultTimeout) * time.Second),
	})

	message.Status = 200
	message.Message = fmt.Sprintf("profiled changed to %v", newProfile)
	res, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(res))
}

// Refresh godoc
// @title Registration
// @description Method to register new user
// @produce json
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @failure 500 {object} models.SuccessOrErrorMessage
// @router /refresh_token [post]
func Refresh(w http.ResponseWriter, r *http.Request) {
	var message models.SuccessOrErrorMessage
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			message.Status = http.StatusUnauthorized
			message.Message = "User is not authorized to refresh session"
			msg, _ := json.Marshal(message)
			_, _ = fmt.Fprintln(w, string(msg))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		message.Status = 200
		message.Message = "Cookie fucked up, I do not know"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	sessionToken := c.Value

	userRecord, err := sessionCache.Get(sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message.Status = http.StatusInternalServerError
		message.Message = "Error while getting your session"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	if userRecord == nil {
		w.WriteHeader(http.StatusUnauthorized)
		message.Status = http.StatusUnauthorized
		message.Message = "User is not authorized to refresh, no session saved"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}
	// (END) The code uptil this point is the same as the first part of the `Welcome` route

	// Now, create a new session token for the current user
	newSessionToken := uuid.New().String()
	userName := userRecord.Username
	timeout := userRecord.Timeout
	err = sessionCache.Set(newSessionToken, &models.UserSessionRecord{
		Username:     userName,
		Timeout:      timeout,
		SessionToken: newSessionToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message.Status = http.StatusInternalServerError
		message.Message = "Error while saving session"
		msg, _ := json.Marshal(message)
		_, _ = fmt.Fprintln(w, string(msg))
		return
	}

	// Delete the older session token
	err = sessionCache.Delete(sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(time.Duration(defaultTimeout) * time.Second),
	})
	message.Status = 200
	message.Message = "Refreshed successfully"
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
}

//

// тестовый контроллер который сейчас контроллируется сессиями

// Welcome godoc
// @title Cheerful method
// @description Method to check sessions consistency
// @produce text/plain
// @success 200 {string} string "Welcome"
// @failure 401 {string} string "Unauthorized"
// @router /welcome [get]
func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintln(w, "UNAUTHORIZED NO COOKIE SET")
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// We then get the name of the user from our sessionCache, where we set the session token
	user, err := sessionCache.Get(sessionToken)
	if err != nil {
		// If there is an error fetching from sessionCache, return an internal server error status
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "INTERNAL SERVER ERROR")
		return
	}
	if user == nil {
		// If the session token is not present in sessionCache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "UNAUTHORZED NO COOKIE IN CACHE")
		return
	}
	// Finally, return the welcome message to the user
	_, _ = fmt.Fprintln(w, fmt.Sprintf("Welcome %s!", user.Username))
}

//тестовый контроллер чтобы понимать что творится на сервере, получая список активных сессий
// TODO: DEPRECATED (mandatory)
func GetSessions(w http.ResponseWriter, r *http.Request) {
	message, _ := json.Marshal(sessionCache.Data)
	_, _ = fmt.Fprintln(w, string(message))
}
