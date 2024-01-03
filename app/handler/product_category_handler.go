package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type ProductCategoryHandler struct {
	uc usecase.ProductCategoryUseCase
}

func NewProductCategoryHandler(uc usecase.ProductCategoryUseCase) *ProductCategoryHandler {
	return &ProductCategoryHandler{uc: uc}
}

func (h *ProductCategoryHandler) GetAllWithoutParams(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()
	paginatedItems, err := h.uc.GetAllProductCategoriesWithoutParams(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductCategoryResponse, 0)
	for _, drugClassification := range paginatedItems.Items.([]*entity.ProductCategory) {
		resps = append(resps, drugClassification.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}
