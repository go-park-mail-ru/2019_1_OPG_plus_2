package server

import (
	"fmt"
	_ "github.com/go-park-mail-ru/2019_1_OPG_plus_2/docs"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/controllers"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/middleware"
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

	router.Use(middleware.ApplyJsonContentType)
	router.Use(middleware.CorsMiddleware)

	router.HandleFunc("/", controllers.MainHandler)
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/api/", controllers.IndexApiHandler)
	router.HandleFunc("/api/sign_in", controllers.SignIn).Methods("POST")
	router.HandleFunc("/api/welcome", controllers.Welcome).Methods("GET")
	router.HandleFunc("/api/refresh_token", controllers.Refresh).Methods("POST")
	router.HandleFunc("/api/register", controllers.Register).Methods("POST")
	router.HandleFunc("/api/sign_out", controllers.SignOut).Methods("POST")
	router.HandleFunc("/api/update", controllers.UpdateProfile).Methods("POST")
	router.HandleFunc("/admin/get_sessions", controllers.GetSessions).Methods("GET")

	return http.ListenAndServe(":"+params.Port, router)
}
