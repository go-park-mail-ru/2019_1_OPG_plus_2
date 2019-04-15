package controllers

import (
	"encoding/json"
	"fmt"
	"2019_1_OPG_plus_2/internal/pkg/adapters"
	"2019_1_OPG_plus_2/internal/pkg/auth"
	"2019_1_OPG_plus_2/internal/pkg/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Response []map[string]interface{}
}

func GenerateApiUrl(accessToken string, fields ...string) string {
	baseUrl := "https://api.vk.com/method/users.get?fields=%s&access_token=%s&v=5.52"
	fieldsString := strings.Join(fields, ",")
	return fmt.Sprintf(baseUrl, fieldsString, accessToken)
}

var vkConfig = oauth2.Config{
	ClientID:     AppId,
	ClientSecret: AppKey,
	RedirectURL:  "http://127.0.0.1:8002/api/callback",
	Endpoint:     vk.Endpoint,
	Scopes:       []string{"email", "friends"},
}

type VkAuthHandlers struct{}

func NewVkAuthHandlers() *VkAuthHandlers {
	return &VkAuthHandlers{}
}

// Code retrieval
// Return redirect to 2nd login stage
func (*VkAuthHandlers) Login1stStageRetrieveCode(w http.ResponseWriter, r *http.Request) {
	url := vkConfig.AuthCodeURL("", oauth2.SetAuthURLParam("display", "popup"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (*VkAuthHandlers) Login2ndStageRetrieveTokenGetData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := r.FormValue("code")

	//retrieving user's access-token
	token, err := vkConfig.Exchange(ctx, code)
	if err != nil {
		log.Println("cannot exchange", err)
		_, _ = w.Write([]byte("=("))
		return
	}

	//creating client with user privileges
	client := vkConfig.Client(ctx, token)

	//getting data
	ApiUrl := GenerateApiUrl(token.AccessToken, "domain", "id", "photo_100", "connections", "site", "email")
	resp, err := client.Get(ApiUrl)
	if err != nil {
		log.Println("cannot request data", err)
		_, _ = w.Write([]byte("=("))
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("cannot read buffer", err)
		_, _ = w.Write([]byte("=("))
		return
	}

	email := token.Extra("email").(string)

	data := &Response{}
	_ = json.Unmarshal(body, data)

	username := data.Response[0]["domain"].(string)
	password := "nil"

	sData := models.SignUpData{
		Username: username,
		Email:    email,
		Password: password,
	}

	jwtData, err, _ := adapters.GetStorages().User.CreateUser(sData)

	if err == models.AlreadyExists {
		jwtData, _, _ = adapters.GetStorages().Auth.SignIn(
			models.SignInData{
				Login:    username,
				Password: password,
			})
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	http.Redirect(w, r, "/api/user", http.StatusTemporaryRedirect)
}
