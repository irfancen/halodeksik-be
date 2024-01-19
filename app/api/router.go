package api

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/applogger"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/handler"
	"halodeksik-be/app/handler/middleware"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/gin-gonic/gin"
)

type RouterOpts struct {
	AddressAreaHandler                 *handler.AddressAreaHandler
	AuthHandler                        *handler.AuthHandler
	CartItemHandler                    *handler.CartItemHandler
	DoctorSpecsHandler                 *handler.DoctorSpecializationHandler
	DrugClassificationHandler          *handler.DrugClassificationHandler
	ForgotTokenHandler                 *handler.ForgotTokenHandler
	ManufacturerHandler                *handler.ManufacturerHandler
	PharmacyHandler                    *handler.PharmacyHandler
	PharmacyProductsHandler            *handler.PharmacyProductHandler
	ProductCategoryHandler             *handler.ProductCategoryHandler
	ProductHandler                     *handler.ProductHandler
	ProductStockMutationHandler        *handler.ProductStockMutationHandler
	ProductStockMutationRequestHandler *handler.ProductStockMutationRequestHandler
	ProfileHandler                     *handler.ProfileHandler
	RegisterTokenHandler               *handler.RegisterTokenHandler
	ShippingMethodHandler              *handler.ShippingMethodHandler
	StockReportHandler                 *handler.StockReportHandler
	UserAddressHandler                 *handler.UserAddressHandler
	UserHandler                        *handler.UserHandler
}

func InitializeAllRouterOpts(allUC *AllUseCases) *RouterOpts {
	return &RouterOpts{
		AddressAreaHandler:                 handler.NewAddressAreaHandler(allUC.AddressAreaUseCase),
		AuthHandler:                        handler.NewAuthHandler(allUC.AuthUseCase, appvalidator.Validator),
		CartItemHandler:                    handler.NewCartItemHandler(allUC.CartItemUseCase, appvalidator.Validator),
		DoctorSpecsHandler:                 handler.NewDoctorSpecializationHandler(allUC.DoctorSpecializationUseCase, appvalidator.Validator),
		DrugClassificationHandler:          handler.NewDrugClassificationHandler(allUC.DrugClassificationUseCase),
		ForgotTokenHandler:                 handler.NewForgotTokenHandler(allUC.ForgotTokenUseCase, appvalidator.Validator),
		ManufacturerHandler:                handler.NewManufacturerHandler(allUC.ManufacturerUseCase, appvalidator.Validator),
		PharmacyHandler:                    handler.NewPharmacyHandler(allUC.PharmacyUseCase, appvalidator.Validator),
		PharmacyProductsHandler:            handler.NewPharmacyProductHAndler(allUC.PharmacyProductUseCase, appvalidator.Validator),
		ProductCategoryHandler:             handler.NewProductCategoryHandler(allUC.ProductCategoryUseCase, appvalidator.Validator),
		ProductHandler:                     handler.NewProductHandler(allUC.ProductUseCase, appvalidator.Validator),
		ProductStockMutationHandler:        handler.NewProductStockMutationHandler(allUC.ProductStockMutation, appvalidator.Validator),
		ProductStockMutationRequestHandler: handler.NewProductStockMutationRequestHandler(allUC.ProductStockMutationRequest, appvalidator.Validator),
		ProfileHandler:                     handler.NewProfileHandler(allUC.ProfileUseCase, appvalidator.Validator),
		RegisterTokenHandler:               handler.NewRegisterTokenHandler(allUC.RegisterTokenUseCase, appvalidator.Validator),
		ShippingMethodHandler:              handler.NewShippingMethodHandler(allUC.ShippingMethodUseCase, appvalidator.Validator),
		StockReportHandler:                 handler.NewStockReportHandler(allUC.ProductStockMutation, appvalidator.Validator),
		UserAddressHandler:                 handler.NewAddressHandler(allUC.UserAddressUseCase, appvalidator.Validator),
		UserHandler:                        handler.NewUserHandler(allUC.UserUseCase, appvalidator.Validator),
	}
}

func GetGinMode() string {
	ginMode := os.Getenv("APP_MODE")
	if ginMode == "" {
		return gin.DebugMode
	}
	return ginMode
}

func NewRouter(rOpts *RouterOpts, ginMode string) *gin.Engine {
	if ginMode == "" {
		applogger.Log.Fatal("gin mode should have some value like \"k6\", \"debug\", or\"release\"")
	}
	gin.SetMode(ginMode)
	router := gin.New()
	router.ContextWithFallback = true

	router.GET("/debug/pprof/", gin.WrapH(http.HandlerFunc(pprof.Index)))
	router.GET("/debug/pprof/profile", gin.WrapH(http.HandlerFunc(pprof.Profile)))
	router.GET("/debug/pprof/heap", gin.WrapH(http.HandlerFunc(pprof.Handler("heap").ServeHTTP)))
	router.GET("/debug/pprof/block", gin.WrapH(http.HandlerFunc(pprof.Handler("block").ServeHTTP)))
	router.GET("/debug/pprof/goroutine", gin.WrapH(http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP)))

	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.TimeoutHandler)
	router.Use(middleware.LogHandler)
	router.Use(middleware.ErrorHandler)

	router.NoRoute(func(ctx *gin.Context) {
		resp := dto.ResponseDto{
			Errors: []string{
				"page not found",
			},
		}
		ctx.JSON(http.StatusNotFound, resp)
	})

	v1 := router.Group("/v1")
	{
		addressArea := v1.Group("/address-area")
		{
			addressArea.GET("/provinces/no-params", rOpts.AddressAreaHandler.GetAllProvince)
			addressArea.GET("/cities/no-params", rOpts.AddressAreaHandler.GetAllCities)
			addressArea.POST("/validate", rOpts.AddressAreaHandler.ValidateLatLong)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/register-token", rOpts.RegisterTokenHandler.SendRegisterToken)
			auth.GET("/verify-register", rOpts.RegisterTokenHandler.VerifyRegisterToken)
			auth.POST("/register", rOpts.AuthHandler.Register)
			auth.POST("/login", rOpts.AuthHandler.Login)
			auth.POST("/forgot-token", rOpts.ForgotTokenHandler.SendForgotToken)
			auth.GET("/verify-forgot", rOpts.ForgotTokenHandler.VerifyForgotToken)
			auth.POST("/reset-password", rOpts.AuthHandler.ResetPassword)
		}

		cartItems := v1.Group("/cart-items")
		{
			cartItems.GET(
				"",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdUser),
				rOpts.CartItemHandler.GetAllByUserId,
			)

			cartItems.GET(
				"/checkout",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdUser),
				rOpts.CartItemHandler.Checkout,
			)

			cartItems.POST(
				"",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdUser),
				rOpts.CartItemHandler.Add,
			)

			cartItems.DELETE(
				"",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdUser),
				rOpts.CartItemHandler.Remove,
			)
		}

		drugClassifications := v1.Group("/drug-classifications")
		{
			drugClassifications.GET("/no-params", rOpts.DrugClassificationHandler.GetAllWithoutParams)
		}

		specs := v1.Group("/doctor-specs")
		{
			specs.GET("/:id", rOpts.DoctorSpecsHandler.GetById)
			specs.GET("/no-params", rOpts.DoctorSpecsHandler.GetAllWithoutParams)
			specs.GET("", rOpts.DoctorSpecsHandler.GetAll)
			specs.POST(
				"",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdAdmin),
				rOpts.DoctorSpecsHandler.Add,
			)
			specs.PUT(
				"/:id",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdAdmin),
				rOpts.DoctorSpecsHandler.Edit,
			)
			specs.DELETE(
				"/:id",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdAdmin),
				rOpts.DoctorSpecsHandler.Remove,
			)
		}

		manufacturers := v1.Group("/manufacturers")
		{
			manufacturers.GET("/:id", rOpts.ManufacturerHandler.GetById)
			manufacturers.GET("/no-params", rOpts.ManufacturerHandler.GetAllWithoutParams)
			manufacturers.GET("", rOpts.ManufacturerHandler.GetAll)
			manufacturers.POST(
				"",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdAdmin),
				rOpts.ManufacturerHandler.Add,
			)
			manufacturers.PUT(
				"/:id",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdAdmin),
				rOpts.ManufacturerHandler.Edit,
			)
			manufacturers.DELETE(
				"/:id",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdAdmin),
				rOpts.ManufacturerHandler.Remove,
			)
		}

		pharmacy := v1.Group(
			"/pharmacies",
			middleware.LoginMiddleware(),
			middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
		)
		{
			pharmacy.GET("", rOpts.PharmacyHandler.GetAll)
			pharmacy.GET("/:id", rOpts.PharmacyHandler.GetById)
			pharmacy.POST("", rOpts.PharmacyHandler.Add)
			pharmacy.PUT("/:id", rOpts.PharmacyHandler.Edit)
			pharmacy.DELETE("/:id", rOpts.PharmacyHandler.Remove)
		}

		pharmacyProducts := v1.Group(
			"/pharmacy-products",
			middleware.LoginMiddleware(),
			middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
		)
		{
			pharmacyProducts.GET("", rOpts.PharmacyProductsHandler.GetAllByPharmacy)
			pharmacyProducts.GET("/:id", rOpts.PharmacyProductsHandler.GetById)
			pharmacyProducts.POST("", rOpts.PharmacyProductsHandler.Add)
			pharmacyProducts.PUT("/:id", rOpts.PharmacyProductsHandler.Edit)
			pharmacyProducts.GET("/:id/request", rOpts.PharmacyProductsHandler.GetAllByProductId)
		}

		productCategories := v1.Group("/product-categories")
		{
			productCategories.GET("/:id", rOpts.ProductCategoryHandler.GetById)
			productCategories.GET("/no-params", rOpts.ProductCategoryHandler.GetAllWithoutParams)
			productCategories.GET("", rOpts.ProductCategoryHandler.GetAll)
			productCategories.POST(
				"",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
				rOpts.ProductCategoryHandler.Add,
			)
			productCategories.PUT(
				"/:id",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
				rOpts.ProductCategoryHandler.Edit,
			)
			productCategories.DELETE(
				"/:id",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
				rOpts.ProductCategoryHandler.Remove,
			)
		}

		products := v1.Group("/products")
		{
			products.GET(
				"/:id/admin",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
				rOpts.ProductHandler.GetById,
			)
			products.GET("/:id", rOpts.ProductHandler.GetByIdForUser)
			products.GET("", rOpts.ProductHandler.GetAll)
			products.GET(
				"/admin",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
				rOpts.ProductHandler.GetAllForAdmin,
			)
			products.POST(
				"",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
				rOpts.ProductHandler.Add,
			)
			products.PUT(
				"/:id",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
				rOpts.ProductHandler.Edit,
			)
			products.DELETE(
				"/:id",
				middleware.LoginMiddleware(),
				middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
				rOpts.ProductHandler.Remove,
			)
		}

		report := v1.Group(
			"/report-stock-mutations",
			middleware.LoginMiddleware(),
			middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
		)
		{
			report.GET("", rOpts.StockReportHandler.FindAll)
		}

		shippingMethod := v1.Group(
			"/shipping-methods",
			middleware.LoginMiddleware(),
			middleware.AllowRoles(appconstant.UserRoleIdUser),
		)
		{
			shippingMethod.POST("", rOpts.ShippingMethodHandler.GetAll)
		}

		stockMutation := v1.Group(
			"/stock-mutations",
			middleware.LoginMiddleware(),
			middleware.AllowRoles(appconstant.UserRoleIdPharmacyAdmin),
		)
		{
			stockMutation.POST("", rOpts.ProductStockMutationHandler.Add)

			stockMutationRequest := stockMutation.Group("/requests")
			{
				stockMutationRequest.GET("/in", rOpts.ProductStockMutationRequestHandler.GetAllIncoming)
				stockMutationRequest.GET("/out", rOpts.ProductStockMutationRequestHandler.GetAllOutgoing)
				stockMutationRequest.POST("", rOpts.ProductStockMutationRequestHandler.Add)
				stockMutationRequest.PATCH("/:id", rOpts.ProductStockMutationRequestHandler.EditStatus)
			}
		}

		users := v1.Group(
			"/users",
			middleware.LoginMiddleware(),
		)
		{
			users.GET("/:id", rOpts.UserHandler.GetById)
			users.GET(
				"",
				middleware.AllowRoles(appconstant.UserRoleIdAdmin),
				rOpts.UserHandler.GetAll,
			)

			admin := users.Group(
				"/admin",
				middleware.AllowRoles(appconstant.UserRoleIdAdmin),
			)
			{
				admin.POST("", rOpts.UserHandler.AddAdmin)
				admin.PATCH("/:id", rOpts.UserHandler.EditAdmin)
				admin.DELETE("/:id", rOpts.UserHandler.RemoveAdmin)
			}
		}

		doctors := v1.Group("/users/doctor")
		{
			doctors.GET("", rOpts.UserHandler.GetAllDoctors)
			doctors.GET("/:id", rOpts.UserHandler.GetDoctorById)
		}

		profile := v1.Group("/profile",
			middleware.LoginMiddleware())
		{
			profileDoctor := profile.Group("/doctor")
			{
				profileDoctor.GET("",
					middleware.AllowRoles(appconstant.UserRoleIdDoctor), rOpts.ProfileHandler.GetProfile)
				profileDoctor.PUT("", middleware.AllowRoles(appconstant.UserRoleIdDoctor), rOpts.ProfileHandler.EditDoctorProfile)
				profileDoctor.POST("/set-online", middleware.AllowRoles(appconstant.UserRoleIdDoctor), rOpts.ProfileHandler.EditDoctorIsOnline)
			}
			profileUser := profile.Group("/user")
			{
				profileUser.GET("",
					middleware.AllowRoles(appconstant.UserRoleIdUser), rOpts.ProfileHandler.GetProfile)
				profileUser.PUT("", middleware.AllowRoles(appconstant.UserRoleIdUser), rOpts.ProfileHandler.EditUserProfile)
			}
			addressProfile := profile.Group("/addresses", middleware.AllowRoles(appconstant.UserRoleIdUser))
			{
				addressProfile.GET("", rOpts.UserAddressHandler.GetAll)
				addressProfile.POST("", rOpts.UserAddressHandler.Add)
				addressProfile.PUT("/:id", rOpts.UserAddressHandler.Update)
				addressProfile.GET("/:id", rOpts.UserAddressHandler.GetById)
				addressProfile.DELETE("/:id", rOpts.UserAddressHandler.Remove)
				addressProfile.GET("/main", rOpts.UserAddressHandler.GetMain)
				addressProfile.POST(":id/main", rOpts.UserAddressHandler.SetMain)
			}
		}

	}

	return router
}
