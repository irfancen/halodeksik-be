package api

import "halodeksik-be/app/usecase"

type AllUseCases struct {
	ManufacturerUseCase usecase.ManufacturerUseCase
	ProductUseCase      usecase.ProductUseCase
}

func InitializeUseCases(allRepo *AllRepositories) *AllUseCases {
	return &AllUseCases{
		ManufacturerUseCase: usecase.NewManufacturerUseCaseImpl(allRepo.ManufacturerRepository),
		ProductUseCase: usecase.NewProductUseCaseImpl(allRepo.ProductRepository),
	}
}
