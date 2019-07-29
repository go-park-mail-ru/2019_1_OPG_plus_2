package auth

import "2019_1_OPG_plus_2/internal/pkg/config"

const CookieName = "jwt"

var secret = []byte(config.Auth.Secret)

var serviceLocation = config.Auth.AuthServiceLocation
var port = config.Auth.AuthPort

var cookieport = config.Auth.CookieServicePort
var cookielocation = config.Auth.CookieServiceLocation
