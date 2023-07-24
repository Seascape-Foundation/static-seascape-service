// Package handler defines the commands and command handlers
// that storage service's reply controller supports
package handler

import (
	"github.com/ahmetson/service-lib/communication/command"
	"github.com/ahmetson/service-lib/controller"
)

const (
	// GetAbi Direct
	GetAbi string = "abi_get"
	// SetAbi Through the router
	SetAbi string = "abi_set"
	// GetConfiguration Through the router
	GetConfiguration string = "configuration_get"
	// SetConfiguration Through the router
	SetConfiguration string = "configuration_set"
	// SetSmartcontract Through the router
	SetSmartcontract string = "smartcontract_set"
	// GetSmartcontract Direct
	GetSmartcontract string = "smartcontract_get"
)

// RegisterCommands registers the commands and their handlers in the controller
func RegisterCommands(c *controller.Controller, extensions ...string) {
	abiGet := command.NewRoute(GetAbi, AbiGet, extensions...)
	abiSet := command.NewRoute(SetAbi, AbiRegister, extensions...)
	smartcontractGet := command.NewRoute(GetSmartcontract, SmartcontractGet, extensions...)
	smartcontractSet := command.NewRoute(SetSmartcontract, SmartcontractRegister, extensions...)
	configurationGet := command.NewRoute(GetConfiguration, ConfigurationGet, extensions...)
	configurationSet := command.NewRoute(SetConfiguration, ConfigurationRegister, extensions...)

	c.AddRoute(abiGet)
	c.AddRoute(abiSet)
	c.AddRoute(smartcontractGet)
	c.AddRoute(smartcontractSet)
	c.AddRoute(configurationGet)
	c.AddRoute(configurationSet)
}
