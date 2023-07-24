package main

import (
	"github.com/ahmetson/service-lib/configuration"
	"github.com/ahmetson/service-lib/controller"
	"github.com/ahmetson/service-lib/independent"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/static-service/handler"
)

func main() {
	logger, err := log.New("static", true)
	if err != nil {
		logger.Fatal("log.New(`main`)", "error", err)
	}

	appConfig, err := configuration.New(logger)
	if err != nil {
		logger.Fatal("configuration.NewAppConfig", "error", err)
	}

	/////////////////////////////////////////////////////////////////////////
	//
	// Create the controller
	//
	/////////////////////////////////////////////////////////////////////////

	replier, err := controller.NewReplier(logger)
	if err != nil {
		logger.Fatal("failed to create a controller", "error", err)
	}
	// prepare the dependencies that controller needs

	replier.RequireExtension("github.com/ahmetson/w3storage-extension")
	// add to replier the commands
	handler.RegisterCommands(replier, "github.com/ahmetson/w3storage-extension")

	service, err := independent.New(appConfig, logger)
	if err != nil {
		logger.Fatal("failed to create an independent service", "error", err)
	}
	service.AddController("main", replier)
	service.RequireProxy("github.com/ahmetson/web-proxy", configuration.DefaultContext)

	err = service.Prepare(configuration.IndependentType)
	if err != nil {
		logger.Fatal("service.Prepare", "error", err)
	}

	service.Run()
}
