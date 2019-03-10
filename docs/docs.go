// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-03-11 00:31:29.640740755 +0300 MSK m=+0.045192719

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "Game based on filling field with color cells",
        "title": "Colors service API by OPG+2",
        "contact": {
            "name": "@Xatabch",
            "email": "nikolsky.dan@gmail.com"
        },
        "license": {
            "name": "Apache 2.0"
        },
        "version": "1.0"
    },
    "host": "localhost:8001",
    "basePath": "/api",
    "paths": {
        "/profile": {
            "put": {
                "description": "This method updates info in profile and auth-db record of user, who is making a query",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Updates client's profile",
                "parameters": [
                    {
                        "description": "User new profile data",
                        "name": "profile_data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.UserData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    }
                }
            },
            "post": {
                "description": "This method creates records about new user in auth-bd and profile-db and then sends cookie to user in order to identify",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Registers user",
                "parameters": [
                    {
                        "description": "User profile data",
                        "name": "profile_data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.UserData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    }
                }
            },
            "delete": {
                "description": "This method deletes all information about user, making a query, including profile, game stats and authorization info",
                "produces": [
                    "application/json"
                ],
                "summary": "Deletes profile and user of client",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    }
                }
            }
        },
        "/profile/{id}": {
            "get": {
                "description": "This method provides client with profile data, matching required ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Produces user profile info",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Profile ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.UserProfile"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    }
                }
            }
        },
        "/upload_avatar": {
            "post": {
                "description": "This method saves avatar image in server storage and sets it as clients user avatar",
                "consumes": [
                    "image/png",
                    "image/jpeg"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Saves new avatar image of client's user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.SuccessOrErrorMessage"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.SuccessOrErrorMessage": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Query processed successfully"
                },
                "status": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "models.UserData": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user_test@test.com"
                },
                "password": {
                    "type": "string",
                    "example": "verysecretpasswordwhichnooneknows"
                },
                "username": {
                    "type": "string",
                    "example": "user_test"
                }
            }
        },
        "models.UserProfile": {
            "type": "object",
            "properties": {
                "avatar_url": {
                    "type": "string",
                    "example": "\u003csome avatar url\u003e"
                },
                "email": {
                    "type": "string",
                    "example": "user_test@test.com"
                },
                "games_played": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "lose": {
                    "type": "integer"
                },
                "score": {
                    "type": "integer"
                },
                "username": {
                    "type": "string",
                    "example": "user_test"
                },
                "win": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
