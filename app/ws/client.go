package ws

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vincent-petithory/dataurl"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appencoder"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"log"
	"os"
	"time"
)

type Client struct {
	Conn    *websocket.Conn
	Message chan *Message
	Id      int64           `json:"id"`
	RoomId  int64           `json:"room_id"`
	Profile *entity.Profile `json:"profile"`
}

type ConsultationMessage struct {
	IsTyping   bool      `json:"is_typing"`
	Message    string    `json:"message"`
	Attachment string    `json:"attachment"`
	CreatedAt  time.Time `json:"created_at"`
}

type Message struct {
	Content ConsultationMessage `json:"content"`
	UserId  int64               `json:"user_id"`
	RoomId  int64               `json:"room_id"`
}

func (c *Client) WriteMessage() {
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			return
		}
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

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		_, jsonMessage, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var consultationMessage ConsultationMessage
		err = appencoder.JsonEncoder.Unmarshal(jsonMessage, &consultationMessage)
		if err != nil {
			break
		}
		consultationMessage.CreatedAt = time.Now()

		msg := &Message{
			Content: consultationMessage,
			UserId:  c.Id,
			RoomId:  c.RoomId,
		}

		if !util.IsEmptyString(msg.Content.Attachment) {
			decodeString, decodeErr := dataurl.DecodeString(msg.Content.Attachment)
			if decodeErr == nil && decodeString.Type == appconstant.DataTypeImage && decodeString.Encoding == appconstant.DataEncodingBase64 {
				myUuid, err2 := uuid.NewRandom()
				if err2 != nil {
					return
				}

				fileName := fmt.Sprintf("%s.%s", myUuid.String(), decodeString.Subtype)
				tempFile, err2 := util.WriteTempFile(decodeString.Data, decodeString.Subtype)

				file, err2 := os.Open(tempFile.Name())
				if err2 != nil {
					return
				}

				ctx, cancel := context.WithTimeout(context.Background(), appconstant.DefaultRequestTimeout*time.Second)
				err2 = appcloud.AppFileUploader.SendToBucketWithFile(ctx, file, "chats/", fileName)
				if err != nil {
					tempFile.Close()
					file.Close()
					os.Remove(tempFile.Name())
					cancel()
				}

				tempFile.Close()
				os.Remove(tempFile.Name())
				cancel()
			}
		}

		hub.Broadcast <- msg
	}
}
