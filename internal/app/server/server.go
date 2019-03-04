package server

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

type Params struct {
	Port string
}

func StartApp(params Params) error {
	fmt.Println("Server starting at " + params.Port)

	router := mux.NewRouter()

	router.HandleFunc("/", controllers.MainHandler)

	return http.ListenAndServe(":"+params.Port, router)
}
