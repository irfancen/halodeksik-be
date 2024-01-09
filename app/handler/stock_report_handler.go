package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type StockReportHandler struct {
	uc        usecase.ProductStockMutationUseCase
	validator appvalidator.AppValidator
}

func NewStockReportHandler(uc usecase.ProductStockMutationUseCase, validator appvalidator.AppValidator) *StockReportHandler {
	return &StockReportHandler{uc: uc, validator: validator}
}

func (h *StockReportHandler) FindAll(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllStockMutationQuery := queryparamdto.GetAllStockMutationsQuery{}
	_ = ctx.ShouldBindQuery(&getAllStockMutationQuery)

	err = h.validator.Validate(getAllStockMutationQuery)
	if err != nil {
		return
	}

	param, err := getAllStockMutationQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAll(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductStockMutationResponse, 0)
	for _, stockMutation := range paginatedItems.Items.([]*entity.ProductStockMutation) {
		resps = append(resps, stockMutation.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)

}
