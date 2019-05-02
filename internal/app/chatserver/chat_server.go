package chatserver

import (
	"2019_1_OPG_plus_2/internal/pkg/chatservice"
	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

var Params = struct {
	Port string
}{
	Port: "8003",
}

func Start() (*chatservice.ChatService, error) {

	service := chatservice.NewChatService(chatservice.NewHub(), tsLogger.NewLogger())

	if err := db.Open(); err != nil {
		service.Log.LogErr("%v", err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			err := db.Ping()
			if err != nil {
				service.Log.LogErr("CHAT: Database seems to be down, trying to reconnect... :%v", err)
			}
		}
	}()

	router := service.AddChatServicePaths(mux.NewRouter())
	go service.Hub.Run()

	service.Log.LogTrace("CHAT: running at %v", Params.Port)

	return service, http.ListenAndServe(":"+Params.Port, router)
}
