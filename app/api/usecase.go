package api

import (
	"halodeksik-be/app/usecase"
	"halodeksik-be/app/util"
)

type AllUseCases struct {
	DrugClassificationUseCase usecase.DrugClassificationUseCase
	ManufacturerUseCase       usecase.ManufacturerUseCase
	ProductCategoryUseCase    usecase.ProductCategoryUseCase
	ProductUseCase            usecase.ProductUseCase
	UserUseCase               usecase.UserUseCase
}

func InitializeUseCases(allRepo *AllRepositories) *AllUseCases {
	return &AllUseCases{
		DrugClassificationUseCase: usecase.NewDrugClassificationUseCaseImpl(allRepo.DrugClassificationRepository),
		ManufacturerUseCase:       usecase.NewManufacturerUseCaseImpl(allRepo.ManufacturerRepository),
		ProductCategoryUseCase:    usecase.NewProductCategoryUseCaseImpl(allRepo.ProductCategoryRepository),
		ProductUseCase:            usecase.NewProductUseCaseImpl(allRepo.ProductRepository),
		UserUseCase:               usecase.NewUserUseCaseImpl(allRepo.UserRepository, util.NewAuthUtil()),
	}
}
