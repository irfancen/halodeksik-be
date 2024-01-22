package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/util"
	"strconv"
)

type GetAllTransactionsQuery struct {
	Limit string `form:"limit"`
	Page  string `form:"page"`
}

func (q *GetAllTransactionsQuery) ToGetAllParams() *GetAllParams {
	param := NewGetAllParams()

	pageSize := appconstant.DefaultGetAllPageSize
	if !util.IsEmptyString(q.Limit) {
		noPageSize, err := strconv.Atoi(q.Limit)
		if err == nil && noPageSize > 0 {
			pageSize = noPageSize
		}
	}
	param.PageSize = &pageSize

	pageId := 1
	if !util.IsEmptyString(q.Page) {
		noPageId, err := strconv.Atoi(q.Page)
		if err == nil && noPageId > 0 {
			pageId = noPageId
		}
	}
	param.PageId = &pageId
	return param
}
