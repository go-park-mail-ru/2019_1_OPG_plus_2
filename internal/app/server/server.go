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
	apiRouter := router.PathPrefix("/api").Subrouter()

	router.HandleFunc("/", controllers.MainHandler)
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	router.Use(middleware.CorsMiddleware)
	apiRouter.Use(middleware.ApplyJsonContentType)

	apiRouter.HandleFunc("/", controllers.IndexApiHandler)
	apiRouter.HandleFunc("/sign_in", controllers.SignIn).Methods("POST")
	apiRouter.HandleFunc("/sign_out", controllers.SignOut).Methods("POST")
	apiRouter.HandleFunc("/register", controllers.Register).Methods("POST")
	apiRouter.HandleFunc("/update", controllers.UpdateProfile).Methods("POST")
	apiRouter.HandleFunc("/refresh_token", controllers.Refresh).Methods("POST")

	router.HandleFunc("/api/welcome", controllers.Welcome).Methods("GET")
	apiRouter.HandleFunc("/admin/get_sessions", controllers.GetSessions).Methods("GET")
	return http.ListenAndServe(":"+params.Port, router)
}
