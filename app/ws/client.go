package ws

import (
	"github.com/gorilla/websocket"
	"halodeksik-be/app/entity"
	"log"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	Id       int64           `json:"id"`
	RoomId   int64           `json:"roomId"`
	Username string          `json:"username"`
	Profile  *entity.Profile `json:"profile"`
}

type Message struct {
	Content  string          `json:"content"`
	UserId   int64           `json:"user_id"`
	RoomId   int64           `json:"roomId"`
	Username string          `json:"username"`
	Profile  *entity.Profile `json:"profile"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		err := c.Conn.WriteJSON(message)
		if err != nil {
			return
		}
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(message),
			UserId:   c.Id,
			RoomId:   c.RoomId,
			Username: c.Username,
			Profile:  c.Profile,
		}

		hub.Broadcast <- msg
	}
}
