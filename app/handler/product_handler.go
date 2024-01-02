package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/dto/uriparamdto"
	"halodeksik-be/app/entity"
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

	req := requestdto.AddEditProduct{}
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

func (h *ProductHandler) GetById(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = wrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	product, err := h.uc.GetById(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: product.ToProductResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetAll(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = wrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllProductQuery := queryparamdto.GetAllProductsQuery{}
	_ = ctx.ShouldBindQuery(&getAllProductQuery)

	param, err := getAllProductQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAll(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductResponse, 0)
	for _, product := range paginatedItems.Items.([]*entity.Product) {
		resps = append(resps, product.ToProductResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) Edit(ctx *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			err = wrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	req := requestdto.AddEditProduct{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	toUpdate, err := req.ToProduct()
	if err != nil {
		return
	}
	updated, err := h.uc.Edit(ctx.Request.Context(), uri.Id, toUpdate)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: updated}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) Remove(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = wrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	err = h.uc.Remove(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusNoContent, dto.ResponseDto{})
}
