package handler

import (
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc        usecase.AuthUsecase
	validator appvalidator.AppValidator
}

func NewAuthHandler(uc usecase.AuthUsecase, v appvalidator.AppValidator) *AuthHandler {
	return &AuthHandler{uc: uc, validator: v}
}

func (h *AuthHandler) SendRegisterToken(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.RequestRegisterToken{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	_, err = h.uc.SendRegisterToken(ctx.Request.Context(), req.Email)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: "Verification link has been sent."}
	ctx.JSON(http.StatusOK, resp)

}

func (h *AuthHandler) VerifyRegisterToken(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.RequestTokenUrl{}
	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	token, err := h.uc.VerifyRegisterToken(ctx.Request.Context(), req.Token)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: token.Email}
	ctx.JSON(http.StatusOK, resp)

}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.RequestRegisterUser{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	reqUri := requestdto.RequestTokenUrl{}
	err = ctx.ShouldBindQuery(&reqUri)
	if err != nil {
		return
	}

	err = h.validator.Validate(reqUri)
	if err != nil {
		return
	}

	user, err := h.uc.Register(ctx.Request.Context(), req.ToUser(), reqUri.Token)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: user.ToUserResponse()}
	ctx.JSON(http.StatusOK, resp)

}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	var req requestdto.LoginRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	log.Println(h.uc)

	user, token, err := h.uc.Login(ctx, req)
	if errors.Is(err, apperror.ErrRecordNotFound) {
		err = apperror.ErrWrongCredentials
		return
	}

	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: responsedto.LoginResponse{
		UserId:     user.Id,
		Email:      user.Email,
		UserRoleId: user.UserRoleId,
		Image:      "",
		Token:      token,
	}}
	ctx.JSON(http.StatusOK, resp)
}
