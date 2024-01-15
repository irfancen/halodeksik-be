package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type ProductStockMutationRequestHandler struct {
	uc        usecase.ProductStockMutationRequestUseCase
	validator appvalidator.AppValidator
}

func NewProductStockMutationRequestHandler(uc usecase.ProductStockMutationRequestUseCase, validator appvalidator.AppValidator) *ProductStockMutationRequestHandler {
	return &ProductStockMutationRequestHandler{uc: uc, validator: validator}
}

func (h *ProductStockMutationRequestHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddProductStockMutationRequest{}
	if err = ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	if err = h.validator.Validate(req); err != nil {
		return
	}

	added, err := h.uc.Add(ctx, req.ToProductStockMutationRequest())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}
