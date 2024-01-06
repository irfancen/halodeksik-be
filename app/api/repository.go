package api

import (
	"database/sql"
	"halodeksik-be/app/repository"
)

type AllRepositories struct {
	CartItemRepository           repository.CartItemRepository
	DrugClassificationRepository repository.DrugClassificationRepository
	ManufacturerRepository       repository.ManufacturerRepository
	PharmacyRepository           repository.PharmacyRepository
	PharmacyProductRepository    repository.PharmacyProductRepository
	ProductCategoryRepository    repository.ProductCategoryRepository
	ProductRepository            repository.ProductRepository
	UserRepository               repository.UserRepository
	VerifyTokenRepository        repository.VerifyTokenRepository
}

func InitializeRepositories(db *sql.DB) *AllRepositories {
	return &AllRepositories{
		CartItemRepository:           repository.NewCartItemRepositoryImpl(db),
		DrugClassificationRepository: repository.NewDrugClassificationRepositoryImpl(db),
		ManufacturerRepository:       repository.NewManufacturerRepositoryImpl(db),
		PharmacyRepository:           repository.NewPharmacyRepository(db),
		PharmacyProductRepository:    repository.NewPharmacyProductRepository(db),
		ProductCategoryRepository:    repository.NewProductCategoryRepositoryImpl(db),
		ProductRepository:            repository.NewProductRepositoryImpl(db),
		UserRepository:               repository.NewUserRepository(db),
		VerifyTokenRepository:        repository.NewVerifyTokenRepository(db),
	}
}
