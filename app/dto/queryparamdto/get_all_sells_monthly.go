package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strings"
)

type GetAllPharmacySellsMonthlyForAdminQuery struct {
	Search            string `form:"search"`
	ProductCategoryId string `form:"product_category_id"`
	Year              int64  `form:"year"`
}

func (q GetAllPharmacySellsMonthlyForAdminQuery) ToGetAllParams() (*GetAllParams, error) {
	param := NewGetAllParams()
	product := new(entity.Product)

	if q.Search != "" {
		words := strings.Split(q.Search, " ")
		wordToSearch := ""
		for _, word := range words {
			wordToSearch += "%" + word + "%"
		}
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(product.GetSqlColumnFromField("Name"), appdb.ILike, wordToSearch),
		)
	}

	param.GroupClauses = append(
		param.GroupClauses,
		appdb.NewGroupClause("month"),
	)

	sortClause := appdb.NewSort("month")
	sortClause.Order = appdb.OrderAsc

	param.SortClauses = append(param.SortClauses, sortClause)

	if !util.IsEmptyString(q.ProductCategoryId) {
		column := product.GetSqlColumnFromField("ProductCategoryId")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.ProductCategoryId))
	}
	monthPageSize := appconstant.MonthInAYearPageSize
	param.PageSize = &monthPageSize

	pageId := 1
	param.PageId = &pageId

	return param, nil
}
