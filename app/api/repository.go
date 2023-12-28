package api

import (
	"database/sql"
	"halodeksik-be/app/repository"
)

type AllRepositories struct {
	ProductRepository repository.ProductRepository
}

func InitializeRepositories(db *sql.DB) *AllRepositories {
	return &AllRepositories{
		ProductRepository: repository.NewProductRepositoryImpl(db),
	}
}
