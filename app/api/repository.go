package api

import (
	"database/sql"
	"halodeksik-be/app/repository"
)

type AllRepositories struct {
	AddressAreaRepository                 repository.AddressAreaRepository
	CartItemRepository                    repository.CartItemRepository
	DoctorSpecializationRepository        repository.DoctorSpecializationRepository
	DrugClassificationRepository          repository.DrugClassificationRepository
	ForgotTokenRepository                 repository.ForgotTokenRepository
	ManufacturerRepository                repository.ManufacturerRepository
	OrderRepository                       repository.OrderRepository
	PharmacyRepository                    repository.PharmacyRepository
	PharmacyProductRepository             repository.PharmacyProductRepository
	ProductCategoryRepository             repository.ProductCategoryRepository
	ProductRepository                     repository.ProductRepository
	ProductStockMutationRepository        repository.ProductStockMutationRepository
	ProductStockMutationRequestRepository repository.ProductStockMutationRequestRepository
	ProfileRepository                     repository.ProfileRepository
	RegisterTokenRepository               repository.RegisterTokenRepository
	TransactionRepository                 repository.TransactionRepository
	ShippingMethodRepository              repository.ShippingMethodRepository
	UserAddressRepository                 repository.UserAddressRepository
	UserRepository                        repository.UserRepository
}

func InitializeRepositories(db *sql.DB) *AllRepositories {
	return &AllRepositories{
		AddressAreaRepository:                 repository.NewAddressAreaRepositoryImpl(db),
		CartItemRepository:                    repository.NewCartItemRepositoryImpl(db),
		DoctorSpecializationRepository:        repository.NewDoctorSpecializationRepositoryImpl(db),
		DrugClassificationRepository:          repository.NewDrugClassificationRepositoryImpl(db),
		ForgotTokenRepository:                 repository.NewForgotTokenRepository(db),
		ManufacturerRepository:                repository.NewManufacturerRepositoryImpl(db),
		OrderRepository:                       repository.NewOrderRepositoryImpl(db),
		PharmacyRepository:                    repository.NewPharmacyRepository(db),
		PharmacyProductRepository:             repository.NewPharmacyProductRepository(db),
		ProductCategoryRepository:             repository.NewProductCategoryRepositoryImpl(db),
		ProductRepository:                     repository.NewProductRepositoryImpl(db),
		ProductStockMutationRepository:        repository.NewProductStockMutationRepositoryImpl(db),
		ProductStockMutationRequestRepository: repository.NewProductStockMutationRequestRepositoryImpl(db),
		ProfileRepository:                     repository.NewProfileRepository(db),
		RegisterTokenRepository:               repository.NewRegisterTokenRepository(db),
		TransactionRepository:                 repository.NewTransactionRepositoryImpl(db),
		ShippingMethodRepository:              repository.NewShippingMethodRepositoryImpl(db),
		UserRepository:                        repository.NewUserRepository(db),
		UserAddressRepository:                 repository.NewUserAddressRepositoryImpl(db),
	}
}
