package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type DoctorSpecializationHandler struct {
	uc        usecase.DoctorSpecializationUseCase
	validator appvalidator.AppValidator
}

func NewDoctorSpecializationHandler(uc usecase.DoctorSpecializationUseCase, validator appvalidator.AppValidator) *DoctorSpecializationHandler {
	return &DoctorSpecializationHandler{uc: uc, validator: validator}
}

func (h *DoctorSpecializationHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddEditDoctorSpecialization{}
	err = ctx.Bind(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	fileHeader, err := ctx.FormFile(appconstant.FormImage)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return
	}
	if fileHeader != nil {
		reqImage := requestdto.AddEditDoctorSpecializationImage{}
		err = ctx.ShouldBind(&reqImage)
		if err != nil {
			return
		}

		err = h.validator.Validate(reqImage)
		if err != nil {
			return
		}

		reqCtx1 := ctx.Request.Context()
		reqCtx2 := context.WithValue(reqCtx1, appconstant.FormImage, fileHeader)
		ctx.Request = ctx.Request.WithContext(reqCtx2)
	}

	added, err := h.uc.Add(ctx.Request.Context(), req.ToDoctorSpecialization())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
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
