package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/responsedto"
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
			err = wrapError(err)
			_ = ctx.Error(err)
		}
	}()
	categories, err := h.uc.GetAllWithoutParams(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductCategoryResponse, 0)
	for _, category := range categories {
		resps = append(resps, category.ToResponse())
	}
	resp := dto.ResponseDto{Data: resps}
	ctx.JSON(http.StatusOK, resp)
}
