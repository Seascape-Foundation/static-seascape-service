package handler

import (
	"github.com/ahmetson/common-lib/data_type/database"
	"github.com/ahmetson/common-lib/topic"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/service-lib/remote"
	"github.com/ahmetson/static-service/configuration"

	"github.com/ahmetson/service-lib/communication/command"
	"github.com/ahmetson/service-lib/communication/message"
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
	if !confTopic.Has("org", "proj") {
		return message.Fail("missing org or proj property in the topic")
	}

	dbCon := remote.GetClient(clients, "database")

	var selectedConf = configuration.Configuration{
		Topic: confTopic,
	}
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
