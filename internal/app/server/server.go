package server

import (
	"2019_1_OPG_plus_2/internal/pkg/gameservice"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"

	_ "2019_1_OPG_plus_2/docs"
	a "2019_1_OPG_plus_2/internal/pkg/adapters"
	"2019_1_OPG_plus_2/internal/pkg/auth"
	"2019_1_OPG_plus_2/internal/pkg/controllers"
	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/middleware"
	"2019_1_OPG_plus_2/internal/pkg/user"
)

type Params struct {
	Port string
}

func StartApp(params Params) error {
	fmt.Println("Server starting at " + params.Port)

	if err := db.Open(); err != nil {
		fmt.Println(err.Error())
	}

	a.SetStorages(user.NewStorage(), auth.NewStorage())
	a.SetHandlers(controllers.NewUserHandlers(), controllers.NewAuthHandlers(), controllers.NewVkAuthHandlers())

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.PanicMiddleware)

	router.HandleFunc("/", controllers.MainHandler)
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	apiRouter.Use(middleware.AuthMiddleware)
	apiRouter.Use(middleware.ApplyJsonContentType)

	apiRouter.HandleFunc("/", controllers.IndexApiHandler)

	apiRouter.HandleFunc("/session", a.GetHandlers().Auth.IsAuth).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/session", a.GetHandlers().Auth.SignIn).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/session", a.GetHandlers().Auth.SignOut).Methods("DELETE", "OPTIONS")
	apiRouter.HandleFunc("/password", a.GetHandlers().Auth.UpdatePassword).Methods("PUT", "OPTIONS")

	apiRouter.HandleFunc("/user", a.GetHandlers().User.GetUser).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/user/{id:[0-9]+}", a.GetHandlers().User.GetUser).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/user", a.GetHandlers().User.CreateUser).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/user", a.GetHandlers().User.UpdateUser).Methods("PUT", "OPTIONS")
	apiRouter.HandleFunc("/user", a.GetHandlers().User.RemoveUser).Methods("DELETE", "OPTIONS")

	apiRouter.HandleFunc("/avatar", controllers.UploadAvatar).Methods("POST", "OPTIONS")

	apiRouter.HandleFunc("/users", controllers.GetScoreboard).Methods("GET", "OPTIONS")

	gameRouter := router.PathPrefix("/game").Subrouter()

	router.PathPrefix("/static").Handler(http.StripPrefix(
		"/static",
		http.FileServer(http.Dir(controllers.StaticPath)),
	))

	apiRouter.HandleFunc("/vk_login", a.GetHandlers().OAuth.Login1stStageRetrieveCode)
	apiRouter.HandleFunc("/callback", a.GetHandlers().OAuth.Login2ndStageRetrieveTokenGetData)
	gameservice.AddGameServicePaths(gameRouter)

	return http.ListenAndServe(":"+params.Port, router)
}

func StopApp() {
	fmt.Println("Stopping server...")
	if err := db.Close(); err != nil {
		fmt.Println(err.Error())
	}
}
