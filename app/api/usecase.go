package api

import (
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/usecase"
)

type AllUseCases struct {
	AuthUsecase                 usecase.AuthUsecase
	CartItemUseCase             usecase.CartItemUseCase
	DrugClassificationUseCase   usecase.DrugClassificationUseCase
	ManufacturerUseCase         usecase.ManufacturerUseCase
	PharmacyUseCase             usecase.PharmacyUseCase
	PharmacyProductUseCase      usecase.PharmacyProductUseCase
	ProductCategoryUseCase      usecase.ProductCategoryUseCase
	ProductStockMutation        usecase.ProductStockMutationUseCase
	ProductUseCase              usecase.ProductUseCase
	UserUseCase                 usecase.UserUseCase
	ProfileUseCase              usecase.ProfileUseCase
	DoctorSpecializationUseCase usecase.DoctorSpecializationUseCase
	ForgotTokenUseCase          usecase.ForgotTokenUseCase
	RegisterTokenUseCase        usecase.RegisterTokenUseCase
}

func InitializeUseCases(allRepo *AllRepositories, allUtil *AllUtil) *AllUseCases {

	forgotTokenUseCase := usecase.NewForgotTokenUsecase(allRepo.UserRepository, allRepo.ForgotTokenRepository, allUtil.AuthUtil, allUtil.MailUtil)
	registerTokenUseCase := usecase.NewRegisterTokenUseCase(allRepo.UserRepository, allRepo.RegisterTokenRepository, allUtil.AuthUtil, allUtil.MailUtil)
	authRepos := usecase.AuthRepos{
		UserRepo:      allRepo.UserRepository,
		TForgotRepo:   allRepo.ForgotTokenRepository,
		TRegisterRepo: allRepo.RegisterTokenRepository,
		ProfileRepo:   allRepo.ProfileRepository,
	}
	authCases := usecase.AuthUseCases{TForgotUseCase: forgotTokenUseCase, TRegisterUseCase: registerTokenUseCase}

	return &AllUseCases{
		AuthUsecase:                 usecase.NewAuthUsecase(authRepos, allUtil.AuthUtil, appcloud.AppFileUploader, authCases),
		CartItemUseCase:             usecase.NewCartItemUseCaseImpl(allRepo.CartItemRepository, allRepo.ProductRepository, allRepo.PharmacyProductRepository),
		DrugClassificationUseCase:   usecase.NewDrugClassificationUseCaseImpl(allRepo.DrugClassificationRepository),
		ManufacturerUseCase:         usecase.NewManufacturerUseCaseImpl(allRepo.ManufacturerRepository),
		PharmacyUseCase:             usecase.NewPharmacyUseCseImpl(allRepo.PharmacyRepository),
		PharmacyProductUseCase:      usecase.NewPharmacyProductUseCaseImpl(allRepo.PharmacyProductRepository, allRepo.PharmacyRepository, allRepo.ProductRepository),
		ProductCategoryUseCase:      usecase.NewProductCategoryUseCaseImpl(allRepo.ProductCategoryRepository),
		ProductUseCase:              usecase.NewProductUseCaseImpl(allRepo.ProductRepository, allRepo.PharmacyRepository, appcloud.AppFileUploader),
		ProductStockMutation:        usecase.NewProductStockMutationUseCaseImpl(allRepo.ProductStockMutationRepository, allRepo.PharmacyProductRepository),
		UserUseCase:                 usecase.NewUserUseCaseImpl(allRepo.UserRepository, allUtil.AuthUtil),
		ProfileUseCase:              usecase.NewProfileUseCaseImpl(allRepo.ProfileRepository),
		DoctorSpecializationUseCase: usecase.NewDoctorSpecializationUseCaseImpl(allRepo.DoctorSpecializationRepository),
		ForgotTokenUseCase:          forgotTokenUseCase,
		RegisterTokenUseCase:        registerTokenUseCase,
	}
}
