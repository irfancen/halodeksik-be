package handler

import (
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PharmacyProductHandler struct {
	uc        usecase.PharmacyProductUseCase
	validator appvalidator.AppValidator
}

func NewPharmacyProductHAndler(uc usecase.PharmacyProductUseCase, validator appvalidator.AppValidator) *PharmacyProductHandler {
	return &PharmacyProductHandler{uc: uc, validator: validator}
}

func (h *PharmacyProductHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			var notFoundError *apperror.NotFound
			if errors.As(err, &notFoundError) {
				err = WrapError(err, http.StatusBadRequest)
			} else {
				err = WrapError(err)
			}
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddPharmacyProduct{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	added, err := h.uc.Add(ctx.Request.Context(), req.ToPharmacyProduct())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToPharmacyProductResponse()}
	ctx.JSON(http.StatusOK, resp)
}
