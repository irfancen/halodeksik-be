package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type ManufacturerHandler struct {
	uc usecase.ManufacturerUseCase
}

func NewManufacturerHandler(uc usecase.ManufacturerUseCase) *ManufacturerHandler {
	return &ManufacturerHandler{uc: uc}
}

func (h *ManufacturerHandler) GetAllWithoutParams(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = wrapError(err)
			_ = ctx.Error(err)
		}
	}()
	paginatedItems, err := h.uc.GetAllManufacturersWithoutParams(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.ManufacturerResponse, 0)
	for _, manufacturer := range paginatedItems.Items.([]*entity.Manufacturer) {
		resps = append(resps, manufacturer.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}
