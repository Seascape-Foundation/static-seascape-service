package main

import (
	"github.com/ahmetson/service-lib/configuration"
	"github.com/ahmetson/service-lib/controller"
	"github.com/ahmetson/service-lib/independent"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/static-service/handler"
)

func main() {
	/////////////////////////////////////////////////////////////////
	//
	// Initiating the data for our service
	//
	//////////////////////////////////////////////////////////////////
	logger, err := log.New("main", true)
	if err != nil {
		logger.Fatal("log.New(`main`)", "error", err)
	}

	logger.Info("load app configuration begin...")
	appConfig, err := configuration.NewAppConfig(logger)
	if err != nil {
		logger.Fatal("configuration.NewAppConfig", "error", err)
	}
	logger.Info("load app configuration end!")

	//
	// Validations
	//
	if len(appConfig.Services) == 0 {
		logger.Fatal("no seascape.yml or the yaml file doesn't contain services")
	}

	logger.Info("for quick development, we assume the seascape.yml has one controller, service and extension")
	serviceConfig := appConfig.Services[0]
	if serviceConfig.Type != configuration.IndependentType {
		logger.Fatal("first service in seascape.yml is not an independent type", "type", serviceConfig.Type)
	}
	if len(serviceConfig.Controllers) == 0 {
		logger.Fatal("first service in seascape.yml doesn't contain controllers")
	}
	controllerConfig := serviceConfig.Controllers[0]
	if controllerConfig.Type != configuration.ReplierType {
		logger.Fatal("the first controller is not a replier",
			"service name", serviceConfig.Name,
			"controller name", controllerConfig.Name,
			"controller type", controllerConfig.Type,
		)
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
	// add to replier the commands
	handler.RegisterCommands(replier)
	// prepare the dependencies that controller needs
	replier.RequireExtension("database")

	service, err := independent.New(appConfig.Services[0])
	if err != nil {
		logger.Fatal("failed to create an independent service", "error", err)
	}
	err = service.AddController("main", replier)
	if err != nil {
		logger.Fatal("failed to add controller into the independent service", "error", err)
	}

	service.Run()
}
