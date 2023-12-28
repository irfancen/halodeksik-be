package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type ProductHandler struct {
	uc        usecase.ProductUseCase
	validator appvalidator.AppValidator
}

func NewProductHandler(uc usecase.ProductUseCase, validator appvalidator.AppValidator) *ProductHandler {
	return &ProductHandler{uc: uc, validator: validator}
}

func (h *ProductHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = wrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddProduct{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	toAdd, err := req.ToProduct()
	if err != nil {
		return
	}
	added, err := h.uc.Add(ctx.Request.Context(), toAdd)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToProductResponse()}
	ctx.JSON(http.StatusOK, resp)
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

	products, err := h.uc.GetAll(ctx.Request.Context(), param)
	if err != nil {
		err = wrapError(err)
		_ = ctx.Error(err)
		return
	}

	resp.Data = products
	ctx.JSON(http.StatusOK, resp)
}
