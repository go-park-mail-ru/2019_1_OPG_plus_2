package chatservice

import (
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	room *ChatRoom
	conn *websocket.Conn
	// channel from which messages are sent via WS to browser or whatever
	send chan []byte
	Log  *tsLogger.TSLogger
}

func NewClient(room *ChatRoom, conn *websocket.Conn) *Client {
	return &Client{
		room: room,
		conn: conn,
		send: make(chan []byte, 256),
		Log:  room.Log,
	}
}

func (c *Client) reading() {
	defer func() {
		c.room.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		c.room.messageHandler <- Message{
			msg:      message,
			feedback: c,
		}
	}
}

// writing pumps messages from the room to the websocket connection.
//
// A goroutine running writing is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writing() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
		c.Log.LogTrace("CHAT: Connection closed")
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.Log.LogTrace("CHAT: Connection closed")
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.Log.LogErr("CHAT: Connection closed: %v", err)
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				c.Log.LogErr("CHAT: Connection closed: %v", err)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Log.LogErr("CHAT: Connection closed: %v", err)
				return
			}
		}
	}
}
