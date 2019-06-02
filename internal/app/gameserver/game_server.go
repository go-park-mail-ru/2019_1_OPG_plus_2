package gameserver

import (
	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/gameservice"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"github.com/gorilla/mux"
	"net/http"
)

var service = gameservice.NewService(tsLogger.NewLogger())

func Start() error {
	router := mux.NewRouter()
	router = service.AddGameServicePaths(router)

	service.Log.Run()

	if err := db.Open(); err != nil {
		service.Log.LogErr("%v", err)
		return err
	}
	go service.Hub.Run()
	service.Log.LogTrace("GAME: Starting server at %v", gameservice.Port)
	return http.ListenAndServe(gameservice.Port, router)
}

//func Stop(){}
