package main

import (
	"halodeksik-be/app/api"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/applogger"
)

func main() {
	logger, logFile := applogger.NewLogger()
	if logFile != nil {
		defer logFile.Close()
	}
	applogger.SetLogger(logger)

	db, dbErr := appdb.Connect()
	if dbErr != nil {
		applogger.Log.Error(dbErr)
	}

	allRepositories := api.InitializeRepositories(db)
	allUseCases := api.InitializeUseCases(allRepositories)
	routerOpts := api.InitializeAllRouterOpts(allUseCases)

	ginMode := api.GetGinMode()
	router := api.NewRouter(routerOpts, ginMode)
	server := api.NewServer(router)
	api.StartAndSetupGracefulShutdown(server)
}
