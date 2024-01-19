package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/dto/uriparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"halodeksik-be/app/util"
	"halodeksik-be/app/ws"
	"net/http"
)

type ChatHandler struct {
	hub                   *ws.Hub
	consultationSessionUC usecase.ConsultationSessionUseCase
	consultationMessageUC usecase.ConsultationMessageUseCase
	profileUC             usecase.ProfileUseCase
	validator             appvalidator.AppValidator
}

func NewChatHandler(
	hub *ws.Hub,
	consultationSessionUC usecase.ConsultationSessionUseCase,
	consultationMessageUC usecase.ConsultationMessageUseCase,
	profileUC usecase.ProfileUseCase,
	validator appvalidator.AppValidator,
) *ChatHandler {
	return &ChatHandler{
		hub:                   hub,
		consultationSessionUC: consultationSessionUC,
		consultationMessageUC: consultationMessageUC,
		profileUC:             profileUC,
		validator:             validator,
	}
}

func (h *ChatHandler) CreateRoom(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddConsultationSession{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	addedOrFound, err := h.consultationSessionUC.Add(ctx, req.ToConsultationSessionUseCase())
	if err != nil && errors.Is(err, apperror.ErrChatStillOngoing) {
		roomId := addedOrFound.Id
		_, isRoomExisted := h.hub.ConsultationSessions[roomId]
		if !isRoomExisted {
			h.hub.ConsultationSessions[roomId] = &ws.ConsultationSession{
				Id:        roomId,
				DoctorId:  req.DoctorId,
				PatientId: req.UserId,
				Clients:   make(map[int64]*ws.Client),
			}
		}
		return
	}

	if err != nil && !errors.Is(err, apperror.ErrChatStillOngoing) {
		return
	}

	roomId := addedOrFound.Id
	h.hub.ConsultationSessions[roomId] = &ws.ConsultationSession{
		Id:        roomId,
		DoctorId:  req.DoctorId,
		PatientId: req.UserId,
		Clients:   make(map[int64]*ws.Client),
	}

	ctx.JSON(http.StatusOK, addedOrFound.ToResponse())
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
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	sessionId := uri.Id
	sessionDb, err := h.consultationSessionUC.GetById(ctx, sessionId)
	if err != nil {
		return
	}

	if sessionDb.ConsultationSessionStatusId != appconstant.ConsultationSessionStatusOngoing {
		err = apperror.ErrChatAlreadyEnded
		return
	}

	_, isRoomExisted := h.hub.ConsultationSessions[sessionId]
	if !isRoomExisted {
		h.hub.ConsultationSessions[sessionId] = &ws.ConsultationSession{
			Id:        sessionId,
			DoctorId:  sessionDb.DoctorId,
			PatientId: sessionDb.UserId,
			Clients:   make(map[int64]*ws.Client),
		}
	}

	clientIdCtx := ctx.Request.Context().Value(appconstant.ContextKeyUserId)
	clientId := clientIdCtx.(int64)

	roleIdCtx := ctx.Request.Context().Value(appconstant.ContextKeyRoleId)
	roleId := roleIdCtx.(int64)

	var user *entity.User

	if roleId == appconstant.UserRoleIdDoctor {
		if sessionDb.DoctorId != clientId {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": []string{"you are not allowed here"},
			})
			return
		}
		user, err = h.profileUC.GetDoctorProfileByUserId(ctx, sessionDb.DoctorId)
		if err != nil {
			return
		}
	}

	if roleId == appconstant.UserRoleIdUser {
		if sessionDb.UserId != clientId {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": []string{"you are not allowed here"},
			})
			return
		}
		user, err = h.profileUC.GetUserProfileByUserId(ctx, sessionDb.UserId)
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
		Conn:      conn,
		Message:   make(chan *ws.Message, 10),
		Id:        clientId,
		SessionId: sessionId,
		Profile:   user.GetProfile(),
	}

	message := &ws.Message{
		Content:   ws.ConsultationMessage{Message: "A new user has joined the session"},
		SenderId:  client.Id,
		SessionId: sessionId,
	}

	h.hub.Register <- client
	h.hub.Broadcast <- message

	go client.WriteMessage()
	go client.ReadMessage(h.hub, h.consultationMessageUC)
}

type RoomRes struct {
	Id        int64 `json:"id"`
	DoctorId  int64 `json:"doctor_id"`
	PatientId int64 `json:"patient_id"`
}

func (h *ChatHandler) GetAllByUserIdOrDoctorId(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	query := queryparamdto.GetAllConsultationSessions{}
	err = ctx.ShouldBindQuery(&query)
	if err != nil {
		return
	}

	err = h.validator.Validate(query)
	if err != nil {
		return
	}

	param := query.ToGetAllParams()
	paginatedItems, err := h.consultationSessionUC.GetAllByUserIdOrDoctorId(ctx, param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ConsultationSessionResponse, 0)
	for _, session := range paginatedItems.Items.([]*entity.ConsultationSession) {
		resps = append(resps, session.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

type ClientRes struct {
	Id      int64           `json:"id"`
	Profile *entity.Profile `json:"profile"`
}

func (h *ChatHandler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomIdQuery := c.Param("roomId")
	roomId, _ := util.ParseInt64(roomIdQuery)

	if _, ok := h.hub.ConsultationSessions[roomId]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.ConsultationSessions[roomId].Clients {
		clients = append(clients, ClientRes{
			Id:      c.Id,
			Profile: c.Profile,
		})
	}

	c.JSON(http.StatusOK, clients)
}
