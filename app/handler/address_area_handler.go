package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type AddressAreaHandler struct {
	uc usecase.AddressAreaUseCase
}

func NewAddressAreaHandler(uc usecase.AddressAreaUseCase) *AddressAreaHandler {
	return &AddressAreaHandler{uc: uc}
}

func (h *AddressAreaHandler) GetAllProvince(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	provinces, err := h.uc.GetAllProvinces(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProvinceResponse, 0)
	for _, province := range provinces {
		resps = append(resps, province.ToResponse())
	}
	resp := dto.ResponseDto{Data: resps}
	ctx.JSON(http.StatusOK, resp)
}

func (h *AddressAreaHandler) GetAllCities(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	cities, err := h.uc.GetAllCities(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.CityResponse, 0)
	for _, city := range cities {
		resps = append(resps, city.ToResponse())
	}
	resp := dto.ResponseDto{Data: resps}
	ctx.JSON(http.StatusOK, resp)
}
