package api

import (
	"database/sql"
	"halodeksik-be/app/repository"
)

type AllRepositories struct {
	ManufacturerRepository repository.ManufacturerRepository
	ProductRepository      repository.ProductRepository
}

func InitializeRepositories(db *sql.DB) *AllRepositories {
	return &AllRepositories{
		ManufacturerRepository: repository.NewManufacturerRepositoryImpl(db),
		ProductRepository:      repository.NewProductRepositoryImpl(db),
	}
}
