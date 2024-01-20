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
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"halodeksik-be/app/util"
	"log"
	"os"
	"time"
)

type Client struct {
	Conn      *websocket.Conn
	Message   chan *responsedto.WsConsultationMessage
	SenderId  int64           `json:"id"`
	SessionId int64           `json:"session_id"`
	Profile   *entity.Profile `json:"profile"`
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

		var consultationMessage requestdto.WsConsultationMessage
		err = appencoder.JsonEncoder.Unmarshal(jsonMessage, &consultationMessage)
		if err != nil {
			break
		}
		consultationMessage.SenderId = c.SenderId
		consultationMessage.SessionId = c.SessionId

		msg := &responsedto.WsConsultationMessage{
			IsTyping:   consultationMessage.IsTyping,
			Message:    consultationMessage.Message,
			Attachment: consultationMessage.Attachment,
			CreatedAt:  time.Now(),
			SenderId:   c.SenderId,
			SessionId:  c.SessionId,
		}

		hub.Broadcast <- msg
		msgToStoreInDb := consultationMessage.ToConsultationMessage()

		if !util.IsEmptyString(msg.Attachment) {
			decodeString, decodeErr := dataurl.DecodeString(msg.Attachment)
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

				msgToStoreInDb.Attachment = appdb.NewSqlNullString(fileUrl)
			}
		}

		_, err = consultationMessageUC.Add(context.Background(), *msgToStoreInDb)
		if err != nil {
			applogger.Log.Error(err)
		}
	}
}
