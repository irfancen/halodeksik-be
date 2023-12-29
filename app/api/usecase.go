package api

import "halodeksik-be/app/usecase"

type AllUseCases struct {
	ProductUseCase usecase.ProductUseCase
}

func InitializeUseCases(allRepo *AllRepositories) *AllUseCases {
	return &AllUseCases{
		ProductUseCase: usecase.NewProductUseCaseImpl(allRepo.ProductRepository),
	}
}
