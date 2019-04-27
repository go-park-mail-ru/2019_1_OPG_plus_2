package chatservice

import (
	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/models"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	msg      []byte
	feedback *Client
}

type ChatRoom struct {
	id      int
	hub     *Hub
	clients map[*Client]bool

	// channel which messages are broadcasted from
	messageHandler chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister        chan *Client
	currentPlayersNum int
	win               bool
}

func newRoom(hub *Hub, id int) *ChatRoom {
	r := &ChatRoom{
		hub:               hub,
		id:                id,
		messageHandler:    make(chan Message),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		clients:           make(map[*Client]bool),
		currentPlayersNum: 0,
	}
	return r
}

func (r *ChatRoom) Run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			r.currentPlayersNum++
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		case message := <-r.messageHandler:
			m, err := r.handleMessage(message)
			if err != nil {
				var bMsg = NewBroadcastErrorMessage(err.Error())
				errorMsg, _ := json.Marshal(&bMsg)
				message.feedback.send <- errorMsg
			} else {
				r.broadcastMsg(m)
			}
			if r.win {
				r.hub.closeRoom(r.id)
			}
		case id := <-r.hub.closer:
			if id == r.id {
				for client := range r.clients {
					close(client.send)
					delete(r.clients, client)
				}
				return
			}
		}
	}
}

func (r *ChatRoom) broadcastMsg(message []byte) {
	for client := range r.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(r.clients, client)
		}
	}
}

func (r *ChatRoom) handleMessage(message Message) ([]byte, error) {
	var msg GenericMessage
	err := json.Unmarshal(message.msg, &msg)
	tsLogger.LogInfo("ROOM %d: %+v", r.id, msg)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing: %v", err)
	}

	switch msg.MType {
	case "text":
		message.msg, err = r.performChatLogic(message)
		return message.msg, err

	default:
		return nil, fmt.Errorf("message type invalid")
	}
}

func (r *ChatRoom) performChatLogic(message Message) ([]byte, error) {
	var chatMessage ChatMessage
	err := json.Unmarshal(message.msg, &chatMessage)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error")
	}

	outMsg := models.ChatMessage{
		Username: chatMessage.User,
		Content:  chatMessage.Content,
		Type:     chatMessage.MType,
		Date:     models.JSONTime(time.Now()),
	}
	outMsg, err = db.CreateMessage(outMsg)
	if err != nil {
		return nil, err
	}

	message.msg, err = json.Marshal(&outMsg)
	if err != nil {
		return nil, err
	}
	return message.msg, nil
}
