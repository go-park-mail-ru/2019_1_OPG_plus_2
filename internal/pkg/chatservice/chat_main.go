package chatservice

import (
	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/middleware"
	"2019_1_OPG_plus_2/internal/pkg/models"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *ChatService) AddChatServicePaths(router *mux.Router) *mux.Router {
	router.HandleFunc("/chat/{id}", func(w http.ResponseWriter, r *http.Request) {
		const pageSize = 10
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			tsLogger.LogWarn("could not parse %d", id)
			return
		}
		if s.Hub.rooms[int(id)] == nil {
			_, _ = fmt.Fprint(w, "no such room with id ", id)
			return
		}

		limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil || limit < 1 {
			limit = pageSize
		}

		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		if err != nil || page < 1 {
			page = 1
		}

		messages, count, err := db.GetMessages(limit, (page-1)*limit)
		if err != nil {
			models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
			return
		}

		models.Send(w, http.StatusOK, models.GetMessageListAnswer(models.ChatMessageList{
			Messages: messages,
			Count:    count,
		}))
	}).Methods("GET")

	router.HandleFunc("/chat/{id}/room", s.ConnectToRoom)

	router.HandleFunc("/chat/{id}", s.DeleteRoom).Methods("DELETE")

	router.HandleFunc("/chat/{id}", s.CreateRoom).Methods("CREATE")

	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.AuthMiddleware)
	router.Use(middleware.PanicMiddleware)
	return router
}

func serveClientConnection(room *ChatRoom, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		tsLogger.LogErr("CONNECTION UPGRADE ERROR: %s", err)
		return err
	}
	client := NewClient(room, conn)
	client.room.register <- client

	go client.writePump()
	go client.readPump()
	return nil
}
