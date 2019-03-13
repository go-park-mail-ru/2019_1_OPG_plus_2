package server

import (
	"fmt"
	_ "github.com/go-park-mail-ru/2019_1_OPG_plus_2/docs"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/controllers"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
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

	if err := db.Open(); err != nil {
	    fmt.Println(err.Error())
    }

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	router.Use(middleware.CorsMiddleware)

	router.HandleFunc("/", controllers.MainHandler)
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	apiRouter.Use(middleware.AuthMiddleware)
	apiRouter.Use(middleware.ApplyJsonContentType)

	apiRouter.HandleFunc("/", controllers.IndexApiHandler)

	apiRouter.HandleFunc("/session", controllers.IsAuth).Methods("GET")
	apiRouter.HandleFunc("/session", controllers.SignIn).Methods("POST")
	apiRouter.HandleFunc("/session", controllers.SignOut).Methods("DELETE")
	apiRouter.HandleFunc("/password", controllers.UpdatePassword).Methods("PUT")

	apiRouter.HandleFunc("/user", controllers.GetUser).Methods("GET")
	apiRouter.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	apiRouter.HandleFunc("/user", controllers.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/user", controllers.UpdateUser).Methods("PUT")
	apiRouter.HandleFunc("/user", controllers.RemoveUser).Methods("DELETE")
	apiRouter.HandleFunc("/avatar", controllers.UploadAvatar).Methods("POST")

	apiRouter.HandleFunc("/users", controllers.GetScoreBoard).Methods("GET")

	staticHandler := http.StripPrefix(
		"/static",
		http.FileServer(http.Dir("/home/daniknik/colors_static/")),
	)
	router.PathPrefix("/static").Handler(staticHandler)

	return http.ListenAndServe(":"+params.Port, router)
}

func StopApp() {
	fmt.Println("Stopping server...")
	if err := db.Close(); err != nil {
        fmt.Println(err.Error())
    }
}
