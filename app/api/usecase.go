package api

import "halodeksik-be/app/usecase"

type AllUseCases struct {
	DrugClassificationUseCase usecase.DrugClassificationUseCase
	ManufacturerUseCase       usecase.ManufacturerUseCase
	ProductCategoryUseCase    usecase.ProductCategoryUseCase
	ProductUseCase            usecase.ProductUseCase
	AuthUsecase               usecase.AuthUsecase
}

func InitializeUseCases(allRepo *AllRepositories, allUtil *AllUtil) *AllUseCases {
	return &AllUseCases{
		DrugClassificationUseCase: usecase.NewDrugClassificationUseCaseImpl(allRepo.DrugClassificationRepository),
		ManufacturerUseCase:       usecase.NewManufacturerUseCaseImpl(allRepo.ManufacturerRepository),
		ProductCategoryUseCase:    usecase.NewProductCategoryUseCaseImpl(allRepo.ProductCategoryRepository),
		ProductUseCase:            usecase.NewProductUseCaseImpl(allRepo.ProductRepository),
		AuthUsecase:               usecase.NewAuthUsecase(allRepo.UserRepository, allRepo.VerifyTokenRepository, allUtil.AuthUtil, allUtil.MailUtil),
	}
}
