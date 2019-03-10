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

	apiRouter.HandleFunc("/session", controllers.IsAuth).Methods("GET")
	apiRouter.HandleFunc("/session", controllers.SignIn).Methods("POST")
	apiRouter.HandleFunc("/session", controllers.SignOut).Methods("DELETE")

	// apiRouter.HandleFunc("/profile", controllers.GetProfile).Methods("GET")
	apiRouter.HandleFunc("/profile/{id}", controllers.GetProfile).Methods("GET")
	apiRouter.HandleFunc("/profile", controllers.CreateProfile).Methods("POST")
	apiRouter.HandleFunc("/profile", controllers.UpdateProfile).Methods("PUT")
	apiRouter.HandleFunc("/profile", controllers.DeleteProfile).Methods("DELETE")

	apiRouter.HandleFunc("/upload_avatar", controllers.UploadAvatar).Methods("POST")

	apiRouter.HandleFunc("/profiles", controllers.GetProfiles).Methods("GET")
	apiRouter.HandleFunc("/profiles/score?page={page}", controllers.ScoreBoardByPage).Methods("GET")

	staticHandler := http.StripPrefix(
		"/static",
		http.FileServer(http.Dir("/home/daniknik/colors_static/")),
	)
	router.PathPrefix("/img").Handler(staticHandler)

	return http.ListenAndServe(":"+params.Port, router)
}
