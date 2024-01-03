package handler

import (
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PharmacyHandler struct {
	uc        usecase.PharmacyUseCase
	validator appvalidator.AppValidator
}

func NewPharmacyHandler(uc usecase.PharmacyUseCase, validator appvalidator.AppValidator) *PharmacyHandler {
	return &PharmacyHandler{uc: uc, validator: validator}
}

func (h *PharmacyHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = wrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddEditPharmacy{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	added, err := h.uc.Add(ctx.Request.Context(), req.ToPharmacy())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToPharmacyResponse()}
	ctx.JSON(http.StatusOK, resp)
}
