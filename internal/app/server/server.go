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

	apiRouter.HandleFunc("/create_profile", controllers.CreateProfile).Methods("POST")
	apiRouter.HandleFunc("/get_profile/{id}", controllers.GetProfile).Methods("GET")
	apiRouter.HandleFunc("/update_profile", controllers.UpdateProfile).Methods("PUT")
	apiRouter.HandleFunc("/delete_profile", controllers.DeleteProfile).Methods("DELETE")

	apiRouter.HandleFunc("/scoreboard/{page}", controllers.ScoreBoardByPage).Methods("GET")
	apiRouter.HandleFunc("/admin/get_profiles", controllers.GetProfiles).Methods("GET")

	staticHandler := http.StripPrefix(
		"/img",
		http.FileServer(http.Dir("/home/daniknik/colors_static/")),
	)

	router.HandleFunc("/api/upload_avatar", controllers.UploadAvatar).Methods("POST")
	router.PathPrefix("/img").Handler(staticHandler)
	return http.ListenAndServe(":"+params.Port, router)
}
