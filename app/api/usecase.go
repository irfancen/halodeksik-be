package api

import (
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/usecase"
)

type AllUseCases struct {
	AuthUsecase               usecase.AuthUsecase
	DrugClassificationUseCase usecase.DrugClassificationUseCase
	ManufacturerUseCase       usecase.ManufacturerUseCase
	PharmacyUseCase           usecase.PharmacyUseCase
	PharmacyProductUseCase    usecase.PharmacyProductUseCase
	ProductCategoryUseCase    usecase.ProductCategoryUseCase
	ProductUseCase            usecase.ProductUseCase
	UserUseCase               usecase.UserUseCase
}

func InitializeUseCases(allRepo *AllRepositories, allUtil *AllUtil) *AllUseCases {
	return &AllUseCases{
		AuthUsecase:               usecase.NewAuthUsecase(allRepo.UserRepository, allRepo.VerifyTokenRepository, allUtil.AuthUtil, allUtil.MailUtil),
		DrugClassificationUseCase: usecase.NewDrugClassificationUseCaseImpl(allRepo.DrugClassificationRepository),
		ManufacturerUseCase:       usecase.NewManufacturerUseCaseImpl(allRepo.ManufacturerRepository),
		PharmacyUseCase:           usecase.NewPharmacyUseCseImpl(allRepo.PharmacyRepository),
		PharmacyProductUseCase:    usecase.NewPharmacyProductUseCaseImpl(allRepo.PharmacyProductRepository, allRepo.PharmacyRepository, allRepo.ProductRepository),
		ProductCategoryUseCase:    usecase.NewProductCategoryUseCaseImpl(allRepo.ProductCategoryRepository),
		ProductUseCase:            usecase.NewProductUseCaseImpl(allRepo.ProductRepository, appcloud.AppFileUploader),
		UserUseCase:               usecase.NewUserUseCaseImpl(allRepo.UserRepository, allUtil.AuthUtil),
	}
}
