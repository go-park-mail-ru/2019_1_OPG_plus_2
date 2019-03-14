package server

import (
	"fmt"
	_ "github.com/go-park-mail-ru/2019_1_OPG_plus_2/docs"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/controllers"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/user"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"net/http"
)

type Params struct {
	Port string
}

func StartApp(params Params) error {
	fmt.Println("Server starting at " + params.Port)

	if err := db.Open(); err != nil {
		fmt.Println(err.Error())
	}

	userHandlers := controllers.NewUserHandlers(user.NewStorageAdapter())

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	router.Use(middleware.CorsMiddleware)

	router.HandleFunc("/", controllers.MainHandler)
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	apiRouter.Use(middleware.AuthMiddleware)
	apiRouter.Use(middleware.ApplyJsonContentType)

	apiRouter.HandleFunc("/", controllers.IndexApiHandler)

	apiRouter.HandleFunc("/session", controllers.IsAuth).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/session", controllers.SignIn).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/session", controllers.SignOut).Methods("DELETE", "OPTIONS")
	apiRouter.HandleFunc("/password", controllers.UpdatePassword).Methods("PUT", "OPTIONS")

	apiRouter.HandleFunc("/user", userHandlers.GetUser).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/user/{id:[0-9]+}", userHandlers.GetUser).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/user", userHandlers.CreateUser).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/user", userHandlers.UpdateUser).Methods("PUT", "OPTIONS")
	apiRouter.HandleFunc("/user", userHandlers.RemoveUser).Methods("DELETE", "OPTIONS")
	apiRouter.HandleFunc("/avatar", controllers.UploadAvatar).Methods("POST", "OPTIONS")

	apiRouter.HandleFunc("/users", controllers.GetScoreboard).Methods("GET", "OPTIONS")

	router.PathPrefix("/static").Handler(http.StripPrefix(
		"/static",
		http.FileServer(http.Dir(controllers.StaticPath)),
	))

	return http.ListenAndServe(":"+params.Port, router)
}

func StopApp() {
	fmt.Println("Stopping server...")
	if err := db.Close(); err != nil {
		fmt.Println(err.Error())
	}
}
