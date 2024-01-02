package api

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/applogger"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/handler"
	"halodeksik-be/app/handler/middleware"
	"net/http"
	"net/http/pprof"
	"os"
)

type RouterOpts struct {
	DrugClassificationHandler *handler.DrugClassificationHandler
	ManufacturerHandler       *handler.ManufacturerHandler
	ProductCategoryHandler    *handler.ProductCategoryHandler
	ProductHandler            *handler.ProductHandler
}

func InitializeAllRouterOpts(allUC *AllUseCases) *RouterOpts {
	return &RouterOpts{
		DrugClassificationHandler: handler.NewDrugClassificationHandler(allUC.DrugClassificationUseCase),
		ManufacturerHandler:       handler.NewManufacturerHandler(allUC.ManufacturerUseCase),
		ProductCategoryHandler:    handler.NewProductCategoryHandler(allUC.ProductCategoryUseCase),
		ProductHandler:            handler.NewProductHandler(allUC.ProductUseCase, appvalidator.Validator),
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
		drugClassifications := v1.Group("/drug-classifications")
		{
			drugClassifications.GET("/no-params", rOpts.DrugClassificationHandler.GetAllWithoutParams)
		}

		manufacturers := v1.Group("/manufacturers")
		{
			manufacturers.GET("/no-params", rOpts.ManufacturerHandler.GetAllWithoutParams)
		}

		productCategories := v1.Group("/product-categories")
		{
			productCategories.GET("/no-params", rOpts.ProductCategoryHandler.GetAllWithoutParams)
		}

		products := v1.Group("/products")
		{
			products.GET("/:id", rOpts.ProductHandler.GetById)
			products.GET("", rOpts.ProductHandler.GetAll)
			products.POST("", rOpts.ProductHandler.Add)
			products.PUT("/:id", rOpts.ProductHandler.Edit)
			products.DELETE("/:id", rOpts.ProductHandler.Remove)
		}
	}

	return router
}
