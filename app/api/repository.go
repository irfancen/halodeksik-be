package api

import (
	"database/sql"
	"halodeksik-be/app/repository"
)

type AllRepositories struct {
	DrugClassificationRepository repository.DrugClassificationRepository
	ManufacturerRepository       repository.ManufacturerRepository
	ProductCategoryRepository    repository.ProductCategoryRepository
	ProductRepository            repository.ProductRepository
}

func InitializeRepositories(db *sql.DB) *AllRepositories {
	return &AllRepositories{
		DrugClassificationRepository: repository.NewDrugClassificationRepositoryImpl(db),
		ManufacturerRepository:       repository.NewManufacturerRepositoryImpl(db),
		ProductCategoryRepository:    repository.NewProductCategoryRepositoryImpl(db),
		ProductRepository:            repository.NewProductRepositoryImpl(db),
	}
}
