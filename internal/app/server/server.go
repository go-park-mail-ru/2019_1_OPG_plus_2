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

	router.Use(middleware.CorsMiddleware)

	router.HandleFunc("/", controllers.MainHandler)
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	apiRouter.Use(middleware.AuthMiddleware)
	apiRouter.Use(middleware.ApplyJsonContentType)

	apiRouter.HandleFunc("/", controllers.IndexApiHandler)

	apiRouter.HandleFunc("/user", controllers.SignUp).Methods("POST")
	apiRouter.HandleFunc("/session", controllers.IsAuth).Methods("GET")
	apiRouter.HandleFunc("/session", controllers.SignIn).Methods("POST")
	apiRouter.HandleFunc("/session", controllers.SignOut).Methods("DELETE")

	return http.ListenAndServe(":"+params.Port, router)
}
