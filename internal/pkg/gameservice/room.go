package gameservice

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	msg      []byte
	feedback *Client
}

type Room struct {
	gameModel *GameModel
	id        int
	hub       *Hub
	clients   map[*Client]bool

	// channel which messages are broadcasted from
	messageHandler chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	maxPlayersNum     int
	currentPlayersNum int
}

func newRoom(hub *Hub, id int) *Room {
	r := &Room{
		hub:               hub,
		id:                id,
		messageHandler:    make(chan Message),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		clients:           make(map[*Client]bool),
		maxPlayersNum:     2,
		currentPlayersNum: 0,
	}
	r.gameModel = NewGameModel(r)
	return r
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			if !(r.currentPlayersNum >= r.maxPlayersNum) {
				r.clients[client] = true
				r.currentPlayersNum++
			}
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		case message := <-r.messageHandler:
			m, err := r.handleMessage(message)
			if err != nil {
				var bMsg = NewBroadcastErrorMessage(err.Error())
				errorm, _ := json.Marshal(&bMsg)
				message.feedback.send <- errorm
			} else {
				r.broadcastMsg(m)
			}
		case id := <-r.hub.closer:
			if id == r.id {
				for client := range r.clients {
					client.send <- []byte("ROOM CLOSES")
					close(client.send)
					delete(r.clients, client)
				}
				return
			}
		}
	}
}

func (r *Room) broadcastMsg(message []byte) {
	for client := range r.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(r.clients, client)
		}
	}
}

func (r *Room) handleMessage(message Message) ([]byte, error) {
	var msg GenericMessage
	err := json.Unmarshal(message.msg, &msg)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing: " + err.Error())
	}

	if msg.User == "SERVICE" {
		return nil, fmt.Errorf("using SERVICE username is illegal))")
	}

	switch msg.MType {
	case "game":
		message.msg, err = r.performGameLogic(message)
		return message.msg, err

	case "chat":
		message.msg, err = r.performChatLogic(message)
		return message.msg, err

	case "register":
		message.msg, err = r.performRegisterLogic(message)
		return message.msg, err

	default:
		return nil, fmt.Errorf("message type invalid")
	}
}

func (r *Room) performGameLogic(message Message) ([]byte, error) {
	if !r.gameModel.IsReady() {
		return nil, fmt.Errorf("room not ready yet")
	}

	var gameAction GameMessage
	err := json.Unmarshal(message.msg, &gameAction)
	if err != nil {
		return nil, err
	}

	err = r.gameModel.DoTurn(gameAction)
	if err != nil {
		return nil, err
	}

	return message.msg, nil
}

func (r *Room) performChatLogic(message Message) ([]byte, error) {
	var chatMessage ChatMessage
	err := json.Unmarshal(message.msg, &chatMessage)

	if err != nil {
		return nil, err
	}

	return message.msg, nil

}

func (r *Room) performRegisterLogic(message Message) ([]byte, error) {
	if r.gameModel.IsReady() {
		return nil, fmt.Errorf("game is already running")
	}
	var registerMessage RegisterMessage
	err := json.Unmarshal(message.msg, &registerMessage)

	if err != nil {
		return nil, err
	}

	if message.feedback.registered {
		return nil, fmt.Errorf("already registered")
	}

	r.gameModel.players = append(r.gameModel.players, registerMessage.User)
	message.feedback.registered = true

	if len(r.gameModel.players) == r.maxPlayersNum {
		r.gameModel.ready = true
	}

	if r.gameModel.IsReady() {
		r.gameModel.Init()
		var dat = NewBroadcastEventMessage("ready", map[string]interface{}{
			"players_num": r.currentPlayersNum,
			"players":     r.gameModel.players,
			"whose_turn":  r.gameModel.players[r.gameModel.whoseTurn],
		})
		m, _ := json.Marshal(&dat)
		return m, nil
	}

	return message.msg, nil
}