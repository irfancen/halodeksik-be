package queryparamdto

import "halodeksik-be/app/appdb"

type GetAllParams struct {
	WhereClauses []appdb.WhereClause
	SortClauses  []appdb.SortClause
	Search       string
	PageId       *int
	PageSize     *int
}

func NewGetAllParams() *GetAllParams {
	return &GetAllParams{
		WhereClauses: make([]appdb.WhereClause, 0),
	}
}
