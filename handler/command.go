// Package handler defines the commands and command handlers
// that storage service's reply controller supports
package handler

import (
	"github.com/Seascape-Foundation/sds-service-lib/communication/command"
	"github.com/Seascape-Foundation/sds-service-lib/controller"
)

const (
	// Direct
	GET_ABI command.Name = "abi_get"
	// Through the router
	SET_ABI command.Name = "abi_set"
	// Through the router
	GET_CONFIGURATION command.Name = "configuration_get"
	// Through the router
	SET_CONFIGURATION command.Name = "configuration_set"
	// Direct
	FILTER_SMARTCONTRACTS command.Name = "smartcontract_filter"
	// Through the router
	FILTER_SMARTCONTRACT_KEYS command.Name = "smartcontract_key_filter"
	// Through the router
	SET_SMARTCONTRACT command.Name = "smartcontract_set"
	// Direct
	GET_SMARTCONTRACT command.Name = "smartcontract_get"
)

// RegisterCommands registers the commands and their handlers in the controller
func RegisterCommands(c *controller.Controller) {
	c.RegisterCommand(GET_ABI, AbiGet)
	c.RegisterCommand(SET_ABI, AbiRegister)
	c.RegisterCommand(GET_SMARTCONTRACT, SmartcontractGet)
	c.RegisterCommand(SET_SMARTCONTRACT, SmartcontractRegister)
	c.RegisterCommand(FILTER_SMARTCONTRACTS, SmartcontractFilter)
	c.RegisterCommand(FILTER_SMARTCONTRACT_KEYS, SmartcontractKeyFilter)
	c.RegisterCommand(GET_CONFIGURATION, ConfigurationGet)
	c.RegisterCommand(SET_CONFIGURATION, ConfigurationRegister)
}
