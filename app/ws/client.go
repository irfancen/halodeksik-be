package ws

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vincent-petithory/dataurl"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/appencoder"
	"halodeksik-be/app/applogger"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"halodeksik-be/app/util"
	"log"
	"os"
	"time"
)

type Client struct {
	Conn      *websocket.Conn
	Message   chan *Message
	Id        int64           `json:"id"`
	SessionId int64           `json:"session_id"`
	Profile   *entity.Profile `json:"profile"`
}

type ConsultationMessage struct {
	IsTyping   bool      `json:"is_typing"`
	Message    string    `json:"message"`
	Attachment string    `json:"attachment"`
	CreatedAt  time.Time `json:"created_at"`
}

type Message struct {
	Content   ConsultationMessage `json:"content"`
	SenderId  int64               `json:"sender_id"`
	SessionId int64               `json:"session_id"`
}

func NewMessage(message *entity.ConsultationMessage) *Message {
	return &Message{
		Content: ConsultationMessage{
			Message:    message.Message.String,
			Attachment: message.Attachment.String,
			CreatedAt:  message.CreatedAt.Time,
		},
		SenderId:  message.SenderId.Int64,
		SessionId: message.SessionId.Int64,
	}
}

func (m *Message) ToEntityConsultationMessage() *entity.ConsultationMessage {
	return &entity.ConsultationMessage{
		SessionId:  appdb.NewSqlNullInt64(m.SessionId),
		SenderId:   appdb.NewSqlNullInt64(m.SenderId),
		Message:    appdb.NewSqlNullString(m.Content.Message),
		Attachment: appdb.NewSqlNullString(m.Content.Attachment),
		CreatedAt:  appdb.NewSqlNullTime(m.Content.CreatedAt),
		UpdatedAt:  appdb.NewSqlNullTime(m.Content.CreatedAt),
	}
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

func (c *Client) ReadMessage(hub *Hub, consultationMessageUC usecase.ConsultationMessageUseCase) {
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
			Content:   consultationMessage,
			SenderId:  c.Id,
			SessionId: c.SessionId,
		}

		hub.Broadcast <- msg

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
				fileUrl, err2 := appcloud.AppFileUploader.UploadFromFile(
					ctx, file, appconfig.Config.GcloudStorageFolderConsultationSessions, fileName,
				)
				if err != nil {
					tempFile.Close()
					file.Close()
					os.Remove(tempFile.Name())
					cancel()
				}

				tempFile.Close()
				os.Remove(tempFile.Name())
				cancel()

				msg.Content.Attachment = fileUrl
			}
		}

		msgToStore := msg.ToEntityConsultationMessage()
		_, err = consultationMessageUC.Add(context.Background(), *msgToStore)
		if err != nil {
			applogger.Log.Error(err)
		}
	}
}
