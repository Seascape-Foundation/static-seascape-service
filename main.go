package main

import (
	"github.com/Seascape-Foundation/sds-service-lib/configuration"
	"github.com/Seascape-Foundation/sds-service-lib/log"
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

	logger.Info("Load app configuration begin...")
	appConfig, err := configuration.NewAppConfig(logger)
	if err != nil {
		logger.Fatal("configuration.NewAppConfig", "error", err)
	}
	logger.Info("Load app configuration end!")

	/////////////////////////////////////////////////////////////////////////
	//
	// Run the Core services:
	//
	/////////////////////////////////////////////////////////////////////////

	// Start the core services
	go Run(appConfig)
}
