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
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	return http.ListenAndServe(":"+params.Port, router)
}
