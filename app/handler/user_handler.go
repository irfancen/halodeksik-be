package handler

import (
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/dto/uriparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc        usecase.UserUseCase
	validator appvalidator.AppValidator
}

func NewUserHandler(uc usecase.UserUseCase, validator appvalidator.AppValidator) *UserHandler {
	return &UserHandler{uc: uc, validator: validator}
}

func (h *UserHandler) AddAdmin(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddAdmin{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	added, err := h.uc.AddAdmin(ctx.Request.Context(), req.ToUser())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToUserResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetById(ctx *gin.Context) {
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

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	user, err := h.uc.GetById(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: user.ToUserResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetAll(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllUserQuery := queryparamdto.GetAllUsersQuery{}
	_ = ctx.ShouldBindQuery(&getAllUserQuery)

	param, err := getAllUserQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAll(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.UserResponse, 0)
	for _, user := range paginatedItems.Items.([]*entity.User) {
		resps = append(resps, user.ToUserResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Edit(ctx *gin.Context) {
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

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	req := requestdto.EditAdmin{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	updated, err := h.uc.Edit(ctx.Request.Context(), uri.Id, req.ToUser())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: updated.ToUserResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Remove(ctx *gin.Context) {
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

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	err = h.uc.Remove(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusNoContent, dto.ResponseDto{})
}
