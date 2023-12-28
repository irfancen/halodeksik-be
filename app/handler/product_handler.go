package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type ProductHandler struct {
	uc usecase.ProductUseCase
}

func NewProductHandler(uc usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

func (h *ProductHandler) GetAll(ctx *gin.Context) {
	resp := dto.ResponseDto{}

	getAllProductQuery := queryparamdto.GetAllProductsQuery{}
	_ = ctx.ShouldBindQuery(&getAllProductQuery)

	param, err := getAllProductQuery.ToGetAllParams()
	if err != nil {
		err = wrapError(err)
		_ = ctx.Error(err)
		return
	}

	products, err := h.uc.GetAll(ctx, param)
	if err != nil {
		err = wrapError(err)
		_ = ctx.Error(err)
		return
	}

	resp.Data = products
	ctx.JSON(http.StatusOK, resp)
}
