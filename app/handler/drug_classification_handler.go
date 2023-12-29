package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type DrugClassificationHandler struct {
	uc usecase.DrugClassificationUseCase
}

func NewDrugClassificationHandler(uc usecase.DrugClassificationUseCase) *DrugClassificationHandler {
	return &DrugClassificationHandler{uc: uc}
}

func (h *DrugClassificationHandler) GetAllWithoutParams(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = wrapError(err)
			_ = ctx.Error(err)
		}
	}()
	drugClassifications, err := h.uc.GetAllWithoutParams(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.DrugClassificationResponse, 0)
	for _, dc := range drugClassifications {
		resps = append(resps, dc.ToResponse())
	}
	resp := dto.ResponseDto{Data: resps}
	ctx.JSON(http.StatusOK, resp)
}
