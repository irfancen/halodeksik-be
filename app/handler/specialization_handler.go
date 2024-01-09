package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type DoctorSpecializationHandler struct {
	uc usecase.DoctorSpecializationUseCase
}

func NewDoctorSpecializationHandler(uc usecase.DoctorSpecializationUseCase) *DoctorSpecializationHandler {
	return &DoctorSpecializationHandler{uc: uc}
}

func (h *DoctorSpecializationHandler) GetAllWithoutParams(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()
	paginatedItems, err := h.uc.GetAllSpecsWithoutParams(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.SpecializationResponse, 0)
	for _, specialization := range paginatedItems.Items.([]*entity.DoctorSpecialization) {
		resps = append(resps, specialization.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}
