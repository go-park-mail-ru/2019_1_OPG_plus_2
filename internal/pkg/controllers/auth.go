package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

//Authentification controllers
var defaultTimeout = 300

var _ = `
/login
/logout
/sign_up
/update_profile
/update_session
`
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type UserRecordDummy struct {
	Username     string `json:"username, string"`
	SessionToken string `json:"session_token, string"`
	Timeout      int    `json:"timeout, int"` // seconds
}

type Cache interface {
	Get(sessionToken string) (*UserRecordDummy, error)
	Set(sessionToken string, dummy *UserRecordDummy) error
	Delete(sessionToken string) error
}

type CookieCacheDummy struct {
	Data map[string]*UserRecordDummy
}

func NewCookieCacheDummy() *CookieCacheDummy {
	return &CookieCacheDummy{Data: make(map[string]*UserRecordDummy)}
}

func (cache *CookieCacheDummy) Set(sessionToken string, dummy *UserRecordDummy) error {
	cache.Data[sessionToken] = dummy
	return nil
}

func (cache *CookieCacheDummy) Get(sessionToken string) (*UserRecordDummy, error) {
	result := cache.Data[sessionToken]
	if result == nil {
		return nil, fmt.Errorf("No cookie in cache")
	}
	return result, nil
}

func (cache *CookieCacheDummy) Delete(sessionToken string) error {
	delete(cache.Data, sessionToken)
	return nil
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

var cache = NewCookieCacheDummy()

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(creds)
	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]
	fmt.Println(expectedPassword)

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new random session token
	sessionToken := uuid.New().String()
	fmt.Println(sessionToken)
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of ${defaultTimeout} seconds
	err = cache.Set(sessionToken, &UserRecordDummy{
		Username:     creds.Username,
		SessionToken: sessionToken,
		Timeout:      defaultTimeout,
	})
	fmt.Println(cache.Get(sessionToken))
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("vse ok")
	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(time.Duration(defaultTimeout) * time.Second),
	})
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintln(w, "UNAUTHORZED NO COOKIE SET")
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// We then get the name of the user from our cache, where we set the session token
	user, err := cache.Get(sessionToken)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "INTERNAL SERVER ERROR")
		return
	}
	if user == nil {
		// If the session token is not present in cache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "UNAUTHORZED NO COOKIE IN CACHE")
		return
	}
	// Finally, return the welcome message to the user
	_, _ = fmt.Fprintln(w, fmt.Sprintf("Welcome %s!", user.Username))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	userRecord, err := cache.Get(sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userRecord == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code uptil this point is the same as the first part of the `Welcome` route

	// Now, create a new session token for the current user
	newSessionToken := uuid.New().String()
	userName := userRecord.Username
	timeout := userRecord.Timeout
	err = cache.Set(newSessionToken, &UserRecordDummy{
		Username:     userName,
		Timeout:      timeout,
		SessionToken: newSessionToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete the older session token
	err = cache.Delete(sessionToken)
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
}
