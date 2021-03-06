package coreserver

import (
	"2019_1_OPG_plus_2/internal/pkg/monitoring"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	tsLogger.Logger.Run()

	if err := db.Open(); err != nil {
		tsLogger.LogErr("%v", err)
	}

	a.SetStorages(
		user.NewStorage(),
		auth.NewStorage(),
	)

	a.SetHandlers(
		controllers.NewUserHandlers(),
		controllers.NewAuthHandlers(),
		controllers.NewVkAuthHandlers(),
	)

	prometheus.MustRegister(
		monitoring.ActiveRooms,
		monitoring.ActiveConns,
		monitoring.AccessSummary,
	)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.PanicMiddleware)

	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/", controllers.MainHandler)
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	apiRouter.Use(middleware.AuthMiddleware)
	apiRouter.Use(middleware.ApplyJsonContentType)
	apiRouter.Use(middleware.AccessMonitoringMiddleware)
	//apiRouter.Use(middleware.AccessLoggingMiddleware)

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

	router.PathPrefix("/upload").Handler(http.StripPrefix(
		"/upload",
		http.FileServer(http.Dir(controllers.StaticPath)),
	))

	//apiRouter.HandleFunc("/vk_login", a.GetHandlers().OAuth.Login1stStageRetrieveCode)
	//apiRouter.HandleFunc("/callback", a.GetHandlers().OAuth.Login2ndStageRetrieveTokenGetData)

	tsLogger.LogTrace("Server starting at " + params.Port)
	return http.ListenAndServe(":"+params.Port, router)
}

func StopApp() {
	tsLogger.LogTrace("Stopping core...")
	if err := db.Close(); err != nil {
		tsLogger.LogErr("%s", err)
	}
}
