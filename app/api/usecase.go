package api

import "halodeksik-be/app/usecase"

type AllUseCases struct {
	DrugClassificationUseCase usecase.DrugClassificationUseCase
	ManufacturerUseCase       usecase.ManufacturerUseCase
	ProductUseCase            usecase.ProductUseCase
}

func InitializeUseCases(allRepo *AllRepositories) *AllUseCases {
	return &AllUseCases{
		DrugClassificationUseCase: usecase.NewDrugClassificationUseCaseImpl(allRepo.DrugClassificationRepository),
		ManufacturerUseCase:       usecase.NewManufacturerUseCaseImpl(allRepo.ManufacturerRepository),
		ProductUseCase:            usecase.NewProductUseCaseImpl(allRepo.ProductRepository),
	}
}
