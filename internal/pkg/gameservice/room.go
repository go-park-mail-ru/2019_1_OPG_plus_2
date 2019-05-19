package gameservice

import (
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"encoding/json"
	"fmt"
)

type Message struct {
	msg      []byte
	feedback *Client
}

type Room struct {
	gameModel *GameModel
	id        string
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
	win               bool
}

func newRoom(hub *Hub, id string) *Room {
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
				r.currentPlayersNum--
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
					var cMsg = NewBroadcastEventMessage("room_close", fmt.Sprintf("room %q closes", r.id))
					closeMsg, _ := json.Marshal(&cMsg)
					client.send <- closeMsg
					close(client.send)
					delete(r.clients, client)
					r.currentPlayersNum--
				}
				return
			}
		}
		m, err := r.CheckReady()
		if err == nil {
			r.broadcastMsg(m)
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
			for i, c := range r.gameModel.players {
				if c == client.username {
					r.gameModel.players = append(r.gameModel.players[:i], r.gameModel.players[i+1:]...)
				}
			}
			r.currentPlayersNum--

		}
	}
}

func (r *Room) handleMessage(message Message) ([]byte, error) {
	var msg GenericMessage
	err := json.Unmarshal(message.msg, &msg)
	tsLogger.LogInfo("ROOM %q: %+v", r.id, msg)
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
	msg := message.msg
	var gameAction GameMessage
	err := json.Unmarshal(message.msg, &gameAction)
	if err != nil {
		return nil, err
	}
	if gameAction.User != r.gameModel.players[r.gameModel.whoseTurn] {
		return nil, fmt.Errorf("it's not your turn")
	}

	err = r.gameModel.DoTurn(gameAction)
	if err != nil {
		return nil, err
	}
	if r.gameModel.Check() {
		u := NewBroadcastEventMessage("win", map[string]interface{}{
			"winner": r.gameModel.players[r.gameModel.whoseTurn],
		})
		msg, _ = json.Marshal(&u)
		r.win = true
		return msg, nil
	}

	return msg, nil
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
	message.feedback.username = registerMessage.User
	message.feedback.registered = true

	if len(r.gameModel.players) == r.maxPlayersNum {
		r.gameModel.ready = true
	}

	return message.msg, nil
}

func (r *Room) CheckReady() ([]byte, error) {
	if r.gameModel.IsRunning() {
		return nil, fmt.Errorf("game is running")
	}
	if r.gameModel.IsReady() {
		r.gameModel.Init()
		r.gameModel.running = true
		var dat = NewBroadcastEventMessage("ready", map[string]interface{}{
			"field":       r.gameModel.GetField(),
			"players_num": r.currentPlayersNum,
			"players":     r.gameModel.players,
			"whose_turn":  r.gameModel.players[r.gameModel.whoseTurn],
		})
		m, _ := json.Marshal(&dat)
		return m, nil
	}
	return nil, fmt.Errorf("not ready yet")
}
