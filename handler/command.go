// Package handler defines the commands and command handlers
// that storage service's reply controller supports
package handler

import (
	"github.com/ahmetson/service-lib/communication/command"
	"github.com/ahmetson/service-lib/controller"
)

const (
	// GetAbi Direct
	GetAbi command.Name = "abi_get"
	// SetAbi Through the router
	SetAbi command.Name = "abi_set"
	// GetConfiguration Through the router
	GetConfiguration command.Name = "configuration_get"
	// SetConfiguration Through the router
	SetConfiguration command.Name = "configuration_set"
	// SetSmartcontract Through the router
	SetSmartcontract command.Name = "smartcontract_set"
	// GetSmartcontract Direct
	GetSmartcontract command.Name = "smartcontract_get"
)

// RegisterCommands registers the commands and their handlers in the controller
func RegisterCommands(c *controller.Controller) {
	c.RegisterCommand(GetAbi, AbiGet)
	c.RegisterCommand(SetAbi, AbiRegister)
	c.RegisterCommand(GetSmartcontract, SmartcontractGet)
	c.RegisterCommand(SetSmartcontract, SmartcontractRegister)
	c.RegisterCommand(GetConfiguration, ConfigurationGet)
	c.RegisterCommand(SetConfiguration, ConfigurationRegister)
}
