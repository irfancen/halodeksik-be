package main

import (
	"halodeksik-be/app/api"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/applogger"
	"halodeksik-be/app/appvalidator"
)

func main() {
	logger, logFile := applogger.NewLogger()
	if logFile != nil {
		defer logFile.Close()
	}
	applogger.SetLogger(logger)

	appValidator := appvalidator.NewAppValidatorImpl()
	appvalidator.SetValidator(appValidator)

	db, dbErr := appdb.Connect()
	if dbErr != nil {
		applogger.Log.Error(dbErr)
	}

	allRepositories := api.InitializeRepositories(db)
	allUtil := api.InitializeUtil()
	allUseCases := api.InitializeUseCases(allRepositories, allUtil)
	routerOpts := api.InitializeAllRouterOpts(allUseCases)

	ginMode := api.GetGinMode()
	router := api.NewRouter(routerOpts, ginMode)
	server := api.NewServer(router)
	api.StartAndSetupGracefulShutdown(server)
}
