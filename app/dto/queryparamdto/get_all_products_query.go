package queryparamdto

import (
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
)

type GetAllProductsQuery struct {
	Search   string `form:"search"`
	SortBy   string `form:"sort_by"`
	Sort     string `form:"sort"`
	FilterBy string `form:"filter_by"`
	Filter   string `form:"filter"`
	Limit    string `form:"limit"`
	Page     string `form:"page"`
}

func (q *GetAllProductsQuery) ToGetAllParams() (*GetAllParams, error) {
	param := NewGetAllParams()

	if q.Search != "" {
		words := strings.Split(q.Search, " ")
		wordToSearch := ""
		for _, word := range words {
			wordToSearch += "%" + word + "%"
		}
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere("products.name", appdb.ILike, wordToSearch, appdb.OR),
			appdb.NewWhere("products.generic_name", appdb.ILike, wordToSearch, appdb.OR),
			appdb.NewWhere("products.description", appdb.ILike, wordToSearch, appdb.OR),
			appdb.NewWhere("products.content", appdb.ILike, wordToSearch),
		)
	}

	switch q.SortBy {
	case "name":
		q.SortBy = "products.name"
	case "price":
		q.SortBy = "products.price"
	case "date":
		q.SortBy = "products.created_at"
	case "":
		q.SortBy = ""
	}
	sortClause := appdb.NewSort(q.SortBy)
	switch q.Sort {
	case strings.ToLower(string(appdb.OrderAsc)):
		sortClause.Order = appdb.OrderAsc
	default:
		sortClause.Order = appdb.OrderDesc
	}
	if !util.IsEmptyString(q.SortBy) {
		param.SortClauses = append(param.SortClauses, sortClause)
	}

	switch q.FilterBy {
	case "drug_class":
		q.FilterBy = "products.drug_classification_id"
	default:
		q.FilterBy = ""
	}
	if !util.IsEmptyString(q.FilterBy) {
		q.Filter = fmt.Sprintf("%v", q.Filter)
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(q.FilterBy, appdb.In, q.Filter))
	}

	pageSize := appconstant.DefaultGetAllPageSize
	if !util.IsEmptyString(q.Limit) {
		noPageSize, err := strconv.Atoi(q.Limit)
		if err == nil {
			pageSize = noPageSize
		}
	}
	param.PageSize = &pageSize

	pageId := 1
	if !util.IsEmptyString(q.Page) {
		noPageId, err := strconv.Atoi(q.Page)
		if err == nil {
			pageId = noPageId
		}
	}
	param.PageId = &pageId

	return param, nil
}
