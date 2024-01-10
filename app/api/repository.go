package api

import (
	"database/sql"
	"halodeksik-be/app/repository"
)

type AllRepositories struct {
	AddressAreaRepository          repository.AddressAreaRepository
	CartItemRepository             repository.CartItemRepository
	DrugClassificationRepository   repository.DrugClassificationRepository
	ManufacturerRepository         repository.ManufacturerRepository
	PharmacyRepository             repository.PharmacyRepository
	PharmacyProductRepository      repository.PharmacyProductRepository
	ProductCategoryRepository      repository.ProductCategoryRepository
	ProductRepository              repository.ProductRepository
	ProductStockMutationRepository repository.ProductStockMutationRepository
	UserRepository                 repository.UserRepository
	ForgotTokenRepository          repository.ForgotTokenRepository
	RegisterTokenRepository        repository.RegisterTokenRepository
	ProfileRepository              repository.ProfileRepository
	DoctorSpecializationRepository repository.DoctorSpecializationRepository
}

func InitializeRepositories(db *sql.DB) *AllRepositories {
	return &AllRepositories{
		AddressAreaRepository:          repository.NewAddressAreaRepositoryImpl(db),
		CartItemRepository:             repository.NewCartItemRepositoryImpl(db),
		DrugClassificationRepository:   repository.NewDrugClassificationRepositoryImpl(db),
		ManufacturerRepository:         repository.NewManufacturerRepositoryImpl(db),
		PharmacyRepository:             repository.NewPharmacyRepository(db),
		PharmacyProductRepository:      repository.NewPharmacyProductRepository(db),
		ProductCategoryRepository:      repository.NewProductCategoryRepositoryImpl(db),
		ProductRepository:              repository.NewProductRepositoryImpl(db),
		ProductStockMutationRepository: repository.NewProductStockMutationRepositoryImpl(db),
		UserRepository:                 repository.NewUserRepository(db),
		ForgotTokenRepository:          repository.NewForgotTokenRepository(db),
		ProfileRepository:              repository.NewProfileRepository(db),
		DoctorSpecializationRepository: repository.NewDoctorSpecializationRepositoryImpl(db),
		RegisterTokenRepository:        repository.NewRegisterTokenRepository(db),
	}
}
