package server

import (
    "../../pkg/controllers"
    "fmt"
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

	return http.ListenAndServe(":" + params.Port, router)
}
