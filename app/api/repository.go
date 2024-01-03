package api

import (
	"database/sql"
	"halodeksik-be/app/repository"
)

type AllRepositories struct {
	DrugClassificationRepository repository.DrugClassificationRepository
	ManufacturerRepository       repository.ManufacturerRepository
	PharmacyRepository           repository.PharmacyRepository
	ProductCategoryRepository    repository.ProductCategoryRepository
	ProductRepository            repository.ProductRepository
	UserRepository               repository.UserRepository
}

func InitializeRepositories(db *sql.DB) *AllRepositories {
	return &AllRepositories{
		DrugClassificationRepository: repository.NewDrugClassificationRepositoryImpl(db),
		ManufacturerRepository:       repository.NewManufacturerRepositoryImpl(db),
		PharmacyRepository:           repository.NewPharmacyRepository(db),
		ProductCategoryRepository:    repository.NewProductCategoryRepositoryImpl(db),
		ProductRepository:            repository.NewProductRepositoryImpl(db),
		UserRepository:               repository.NewUserRepository(db),
	}
}
