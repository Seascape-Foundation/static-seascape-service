package handler

import (
	"github.com/Seascape-Foundation/sds-common-lib/data_type/database"
	"github.com/Seascape-Foundation/sds-common-lib/topic"
	"github.com/Seascape-Foundation/sds-service-lib/log"
	"github.com/Seascape-Foundation/sds-service-lib/remote"
	"github.com/Seascape-Foundation/static-seascape-service/configuration"

	"github.com/Seascape-Foundation/sds-service-lib/communication/command"
	"github.com/Seascape-Foundation/sds-service-lib/communication/message"
)

type GetConfigurationRequest = topic.Topic
type GetConfigurationReply = configuration.Configuration

type SetConfigurationRequest = configuration.Configuration
type SetConfigurationReply = configuration.Configuration

// ConfigurationRegister Register a new smartcontract in the configuration.
// It requires smartcontract address. First call smartcontract_register command.
var ConfigurationRegister = func(request message.Request, _ log.Logger, clients remote.Clients) message.Reply {
	var conf SetConfigurationRequest
	err := request.Parameters.Interface(&conf)
	if err != nil {
		return message.Fail("failed to parse data")
	}
	if err := conf.Validate(); err != nil {
		return message.Fail("validation: " + err.Error())
	}

	dbCon := remote.GetClient(clients, "database")
	var crud database.Crud = &conf
	if err = crud.Insert(dbCon); err != nil {
		return message.Fail("Configuration saving in the database failed: " + err.Error())
	}

	var reply = conf
	replyMessage, err := command.Reply(&reply)
	if err != nil {
		return message.Fail("failed to reply")
	}

	return replyMessage
}

// ConfigurationGet Returns configuration and smartcontract information related to the configuration
func ConfigurationGet(request message.Request, _ log.Logger, clients remote.Clients) message.Reply {
	var confTopic GetConfigurationRequest
	err := request.Parameters.Interface(&confTopic)
	if err != nil {
		return message.Fail("failed to parse data")
	}
	if err := confTopic.Validate(); err != nil {
		return message.Fail("invalid topic: " + err.Error())
	}
	if confTopic.Level() != topic.SmartcontractLevel {
		return message.Fail("topic level is not at SMARTCONTRACT LEVEL")
	}

	dbCon := remote.GetClient(clients, "database")

	var selectedConf = configuration.Configuration{}
	err = selectedConf.Select(dbCon)
	if err != nil {
		return message.Fail("failed to get configuration from the database: " + err.Error())
	}

	replyMessage, err := command.Reply(&selectedConf)
	if err != nil {
		return message.Fail("failed to reply")
	}

	return replyMessage
}
