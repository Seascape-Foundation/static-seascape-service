package handler

import (
	"testing"

	"github.com/ahmetson/common-lib/data_type/key_value"
	"github.com/ahmetson/common-lib/topic"
	"github.com/ahmetson/service-lib/communication/message"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/static-service/configuration"
	"github.com/stretchr/testify/suite"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestConfigurationSuite struct {
	suite.Suite
	logger    log.Logger
	conf      configuration.Configuration
	conf_list *key_value.List
}

func (suite *TestConfigurationSuite) SetupTest() {
	logger, err := log.New("test", true)
	suite.Require().NoError(err)
	suite.logger = logger

	conf_0 := configuration.Configuration{
		Topic: topic.Topic{
			Organization: "test_org",
			Project:      "test_proj",
			NetworkId:    "test_1",
			Group:        "test_group",
			Name:         "test_name",
		},
	}
	suite.conf = conf_0

	conf_1 := configuration.Configuration{
		Topic: topic.Topic{
			Organization: "test_org_1",
			Project:      "test_proj_1",
			NetworkId:    "test_1",
			Group:        "test_group_1",
			Name:         "test_name_1",
		},
	}

	list := key_value.NewList()
	err = list.Add(conf_0.Topic, &conf_0)
	suite.Require().NoError(err)

	err = list.Add(conf_1.Topic, &conf_1)
	suite.Require().NoError(err)
	suite.conf_list = list
}

func (suite *TestConfigurationSuite) TestGet() {
	// valid request
	valid_kv, err := key_value.NewFromInterface(suite.conf.Topic)
	suite.Require().NoError(err)

	request := message.Request{
		Command:    "",
		Parameters: valid_kv,
	}
	reply := ConfigurationGet(request, suite.logger, nil)
	suite.Require().True(reply.IsOK())

	var replied_sm GetConfigurationReply
	err = reply.Parameters.Interface(&replied_sm)
	suite.Require().NoError(err)

	suite.Require().EqualValues(suite.conf, replied_sm)

	// request with empty parameter should fail
	request = message.Request{
		Command:    "",
		Parameters: key_value.Empty(),
	}
	reply = ConfigurationGet(request, suite.logger, nil)
	suite.Require().False(reply.IsOK())

	// request of configuration that
	// doesn't exist in the list
	// should fail
	no_topic := topic.Topic{
		Organization: "test_org_2",
		Project:      "test_proj_2",
		NetworkId:    "test_1",
		Group:        "test_group_2",
		Name:         "test_name_2",
	}
	topic_kv, err := key_value.NewFromInterface(no_topic)
	suite.Require().NoError(err)

	request = message.Request{
		Command:    "",
		Parameters: topic_kv,
	}
	reply = ConfigurationGet(request, suite.logger, nil)
	suite.Require().False(reply.IsOK())

	// requesting with invalid type for abi id should fail
	no_topic = topic.Topic{
		Organization: "test_org_2",
		Project:      "test_proj_2",
		NetworkId:    "test_1",
		Group:        "test_group_2",
	}
	topic_kv, err = key_value.NewFromInterface(no_topic)
	suite.Require().NoError(err)
	request = message.Request{
		Command:    "",
		Parameters: topic_kv,
	}
	reply = ConfigurationGet(request, suite.logger, nil)
	suite.Require().False(reply.IsOK())
}

func (suite *TestConfigurationSuite) TestSet() {
	// valid request
	no_topic := topic.Topic{
		Organization: "test_org_2",
		Project:      "test_proj_2",
		NetworkId:    "test_1",
		Group:        "test_group_2",
		Name:         "test_name_2",
	}
	valid_request := configuration.Configuration{
		Topic: no_topic,
	}
	valid_kv, err := key_value.NewFromInterface(valid_request)
	suite.Require().NoError(err)

	request := message.Request{
		Command:    "",
		Parameters: valid_kv,
	}
	reply := ConfigurationRegister(request, suite.logger, nil)
	suite.T().Log(reply.Message)
	suite.Require().True(reply.IsOK())

	var repliedSm GetConfigurationReply
	err = reply.Parameters.Interface(&repliedSm)
	suite.Require().NoError(err)
	suite.Require().EqualValues(valid_request, repliedSm)

	// the abi list should have the item
	sm_in_list, err := suite.conf_list.Get(repliedSm)
	suite.Require().NoError(err)
	suite.Require().EqualValues(&repliedSm, sm_in_list)

	// registering with empty parameter should fail
	request = message.Request{
		Command:    "",
		Parameters: key_value.Empty(),
	}
	reply = ConfigurationRegister(request, suite.logger, nil)
	suite.Require().False(reply.IsOK())

	// registering of abi that already exist in the list
	// should fail
	request = message.Request{
		Command:    "",
		Parameters: valid_kv,
	}
	reply = ConfigurationRegister(request, suite.logger, nil)
	suite.Require().False(reply.IsOK())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestConfiguration(t *testing.T) {
	suite.Run(t, new(TestConfigurationSuite))
}
