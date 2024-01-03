package main

import (
	"fmt"
	"halodeksik-be/app/api"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/applogger"
	"halodeksik-be/app/appvalidator"
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
	err = validator.AddNewCustomValidation("filesize", appvalidator.FileSizeValidation)
	if err != nil {
		applogger.Log.Error(err)
	}
	err = validator.AddNewCustomValidation("filetype", appvalidator.FileTypeValidation)
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
	routerOpts := api.InitializeAllRouterOpts(allUseCases)

	ginMode := api.GetGinMode()
	router := api.NewRouter(routerOpts, ginMode)
	server := api.NewServer(router)
	api.StartAndSetupGracefulShutdown(server)
}
