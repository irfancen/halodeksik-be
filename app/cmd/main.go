package main

import (
	"fmt"
	"halodeksik-be/app/api"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/applogger"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/ws"
	"os"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		applogger.Log.Error(err)
	}
	tmpdir := fmt.Sprintf("%s/%s", pwd, "tmp")
	err = os.Setenv("TMPDIR", tmpdir)
	if err != nil {
		applogger.Log.Error(err)
	}

	logger, logFile := applogger.NewLogger()
	if logFile != nil {
		defer logFile.Close()
	}
	applogger.SetLogger(logger)

	validator := appvalidator.NewAppValidatorImpl()
	if err := appvalidator.AddCustomValidators(validator); err != nil {
		applogger.Log.Error(err)
	}
	appvalidator.SetValidator(validator)

	fileUploader := appcloud.NewFileUploaderImpl()
	appcloud.SetAppFileUploader(fileUploader)

	db, dbErr := appdb.Connect()
	if dbErr != nil {
		applogger.Log.Error(dbErr)
	}

	allRepositories := api.InitializeRepositories(db)
	allUtil := api.InitializeUtil()
	allUseCases := api.InitializeUseCases(allRepositories, allUtil)
	hub := ws.NewHub()
	routerOpts := api.InitializeAllRouterOpts(allUseCases, hub)

	go hub.Run()
	ginMode := api.GetGinMode()
	router := api.NewRouter(routerOpts, ginMode)
	server := api.NewServer(router)
	api.StartAndSetupGracefulShutdown(server)
}
