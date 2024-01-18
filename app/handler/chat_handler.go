package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/uriparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"halodeksik-be/app/util"
	"halodeksik-be/app/ws"
	"net/http"
)

type ChatHandler struct {
	hub       *ws.Hub
	profileUC usecase.ProfileUseCase
	validator appvalidator.AppValidator
}

func NewChatHandler(hub *ws.Hub, profileUC usecase.ProfileUseCase, validator appvalidator.AppValidator) *ChatHandler {
	return &ChatHandler{hub: hub, profileUC: profileUC, validator: validator}
}

func (h *ChatHandler) CreateRoom(ctx *gin.Context) {
	var err error

	req := requestdto.AddChatRoom{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	roomId := int64(len(h.hub.Rooms) + 1)
	h.hub.Rooms[roomId] = &ws.Room{
		Id:        roomId,
		DoctorId:  req.DoctorId,
		PatientId: req.PatientId,
		Clients:   make(map[int64]*ws.Client),
	}

	ctx.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		//return origin == "http://localhost:3000"
		return true
	},
}

func (h *ChatHandler) JoinRoom(ctx *gin.Context) {
	var err error

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	roomId := uri.Id
	room, isRoomExisted := h.hub.Rooms[roomId]
	if !isRoomExisted {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": []string{"room not found"},
		})
		return
	}

	clientIdCtx := ctx.Request.Context().Value(appconstant.ContextKeyUserId)
	clientId := clientIdCtx.(int64)

	roleIdCtx := ctx.Request.Context().Value(appconstant.ContextKeyRoleId)
	roleId := roleIdCtx.(int64)

	var user *entity.User

	if roleId == appconstant.UserRoleIdDoctor {
		if room.DoctorId != clientId {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": []string{"you are not allowed here"},
			})
			return
		}
		user, err = h.profileUC.GetDoctorProfileByUserId(ctx, room.DoctorId)
		if err != nil {
			return
		}
	}

	if roleId == appconstant.UserRoleIdUser {
		if room.PatientId != clientId {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": []string{"you are not allowed here"},
			})
			return
		}
		user, err = h.profileUC.GetUserProfileByUserId(ctx, room.PatientId)
		if err != nil {
			return
		}
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &ws.Client{
		Conn:    conn,
		Message: make(chan *ws.Message, 10),
		Id:      clientId,
		RoomId:  roomId,
		Profile: user.GetProfile(),
	}

	message := &ws.Message{
		Content: "A new user has joined the room",
		UserId:  client.Id,
		RoomId:  roomId,
	}

	h.hub.Register <- client
	h.hub.Broadcast <- message

	go client.WriteMessage()
	go client.ReadMessage(h.hub)
}

type RoomRes struct {
	Id        int64 `json:"id"`
	DoctorId  int64 `json:"doctor_id"`
	PatientId int64 `json:"patient_id"`
}

func (h *ChatHandler) GetRooms(ctx *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, room := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			Id:        room.Id,
			DoctorId:  room.DoctorId,
			PatientId: room.PatientId,
		})
	}

	ctx.JSON(http.StatusOK, rooms)
}

type ClientRes struct {
	Id      int64          `json:"id"`
	Profile *entity.Profile `json:"profile"`
}

func (h *ChatHandler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomIdQuery := c.Param("roomId")
	roomId, _ := util.ParseInt64(roomIdQuery)

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			Id:      c.Id,
			Profile: c.Profile,
		})
	}

	c.JSON(http.StatusOK, clients)
}
