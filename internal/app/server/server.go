package server

import (
	"fmt"
	_ "github.com/go-park-mail-ru/2019_1_OPG_plus_2/docs"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/controllers"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"net/http"
)

type Params struct {
	Port string
}

func StartApp(params Params) error {
	fmt.Println("Server starting at " + params.Port)

	router := mux.NewRouter()

	router.HandleFunc("/", controllers.MainHandler)
	router.HandleFunc("/api.", controllers.IndexApiHandler)
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/sign_in", controllers.SignIn).Methods("POST")
	router.HandleFunc("/welcome", controllers.Welcome).Methods("GET")
	router.HandleFunc("/refresh_token", controllers.Refresh).Methods("POST")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/sign_out", controllers.SignOut).Methods("POST")
	router.HandleFunc("/update", controllers.UpdateProfile).Methods("POST")

	return http.ListenAndServe(":"+params.Port, router)
}
