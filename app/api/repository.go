package api

import "database/sql"

type AllRepositories struct {
}

func InitializeRepositories(db *sql.DB) *AllRepositories {
	return &AllRepositories{}
}
