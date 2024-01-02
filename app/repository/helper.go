package repository

import (
	"fmt"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/util"
	"log"
	"strings"
)

func buildQuery(initQuery string, param *queryparamdto.GetAllParams, isLimitSet ...bool) (string, []interface{}) {
	var query strings.Builder
	var values []interface{}

	query.WriteString(initQuery)

	if len(param.WhereClauses) > 0 {
		query.WriteString(appdb.AND + " ")
	}

	indexPreparedStatement := 0

	for index, whereClause := range param.WhereClauses {
		if whereClause.Condition == appdb.In {
			query.WriteString(fmt.Sprintf("%s %s (", whereClause.Column, whereClause.Condition))
			val := strings.Split(whereClause.Value.(string), ",")
			for idx, v := range val {
				indexPreparedStatement++
				query.WriteString(fmt.Sprintf("$%d", indexPreparedStatement))
				if idx != len(val)-1 {
					query.WriteString(",")
				}
				values = append(values, v)
			}
			query.WriteString(string(") " + whereClause.Logic))
			continue
		}

		indexPreparedStatement++
		query.WriteString(fmt.Sprintf("%s %s $%d %s ", whereClause.Column, whereClause.Condition, indexPreparedStatement, whereClause.Logic))

		if index != len(param.WhereClauses)-1 && util.IsEmptyString(string(whereClause.Logic)) {
			query.WriteString(appdb.AND + " ")
		}

		values = append(values, whereClause.Value)
	}

	setLimit := true
	if len(isLimitSet) > 0 {
		setLimit = isLimitSet[0]
	}

	if setLimit {
		query.WriteString(" GROUP BY id ")
	}

	if setLimit {
		query.WriteString(" ORDER BY ")
		for _, sortClause := range param.SortClauses {
			query.WriteString(fmt.Sprintf("%s %s,", sortClause.Column, sortClause.Order))
		}
		query.WriteString(` id ASC `)
	}

	if setLimit && param.PageId != nil && param.PageSize != nil {
		size := *param.PageSize
		offset := (*param.PageId - 1) * size

		query.WriteString(fmt.Sprintf("LIMIT $%d OFFSET $%d", indexPreparedStatement+1, indexPreparedStatement+2))
		values = append(values, size, offset)
	}
	log.Println(query.String())
	return query.String(), values
}
