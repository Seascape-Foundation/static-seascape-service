package handler

import (
	"testing"

	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
	"github.com/Seascape-Foundation/sds-common-lib/smartcontract_key"
	"github.com/Seascape-Foundation/sds-common-lib/topic"
	"github.com/Seascape-Foundation/sds-service-lib/log"
	"github.com/Seascape-Foundation/static-seascape-service/configuration"
	"github.com/Seascape-Foundation/static-seascape-service/smartcontract"
	"github.com/stretchr/testify/suite"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestFilterSuite struct {
	suite.Suite
	logger   log.Logger
	conf     configuration.Configuration
	sm       smartcontract.Smartcontract
	confList *key_value.List
	smList   *key_value.List
}

/*
Two organization

	first one has 1 conf
	second one has many
*/
func (suite *TestFilterSuite) SetupTest() {
	logger, err := log.New("test", false)
	suite.Require().NoError(err)
	suite.logger = logger

	sm0 := smartcontract.Smartcontract{
		SmartcontractKey: smartcontract_key.Key{
			NetworkId: "test_1",
			Address:   "0xaddr_0",
		},
		AbiId: "abi",
	}
	suite.sm = sm0

	sm1 := smartcontract.Smartcontract{
		SmartcontractKey: smartcontract_key.Key{
			NetworkId: "test_1",
			Address:   "0xaddr_1",
		},
		AbiId: "abi",
	}

	sm2 := smartcontract.Smartcontract{
		SmartcontractKey: smartcontract_key.Key{
			NetworkId: "test_2",
			Address:   "0xaddr_2",
		},
		AbiId: "abi",
	}

	conf0 := configuration.Configuration{
		Topic: topic.Topic{
			Organization:  "test_org",
			Project:       "test_proj",
			NetworkId:     "test_1",
			Group:         "test_group",
			Smartcontract: "test_name",
		},
		Address: "0xaddr_0",
	}
	suite.conf = conf0

	conf1 := configuration.Configuration{
		Topic: topic.Topic{
			Organization:  "test_org_1",
			Project:       "test_proj_1",
			NetworkId:     "test_1",
			Group:         "test_group_1",
			Smartcontract: "test_name_1",
		},
		Address: "0xaddr_1",
	}

	conf2 := configuration.Configuration{
		Topic: topic.Topic{
			Organization:  "test_org",
			Project:       "test_proj_2",
			NetworkId:     "test_2",
			Group:         "test_group_2",
			Smartcontract: "test_name_2",
		},
		Address: "0xaddr_2",
	}

	list := key_value.NewList()
	err = list.Add(conf0.Topic, &conf0)
	suite.Require().NoError(err)

	err = list.Add(conf1.Topic, &conf1)
	suite.Require().NoError(err)
	suite.confList = list

	err = list.Add(conf2.Topic, &conf2)
	suite.Require().NoError(err)
	suite.confList = list

	smList := key_value.NewList()
	err = smList.Add(sm0.SmartcontractKey, &sm0)
	suite.Require().NoError(err)

	err = smList.Add(sm1.SmartcontractKey, &sm1)
	suite.Require().NoError(err)

	err = smList.Add(sm2.SmartcontractKey, &sm2)
	suite.Require().NoError(err)

	suite.smList = smList
}

func (suite *TestFilterSuite) TestOrganizationFilter() {
	// empty paths should return all configurations
	var paths []string
	newList := filterOrganization(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 3)
	suite.Require().False(newList.IsEmpty())

	// fetching the non-existing paths should return empty list
	paths = []string{"no_org"}
	newList = filterOrganization(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 0)
	suite.Require().True(newList.IsEmpty())

	// fetching the org that has one element
	paths = []string{"test_org_1"}
	newList = filterOrganization(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 1)
	suite.Require().False(newList.IsEmpty())

	// fetching the org that has two element
	paths = []string{"test_org"}
	newList = filterOrganization(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 2)
	suite.Require().False(newList.IsEmpty())

	paths = []string{"test_org", "test_org_1"}
	newList = filterOrganization(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 3)
	suite.Require().False(newList.IsEmpty())

	paths = []string{"test_org", "test_org_1", "non_exist"}
	newList = filterOrganization(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 3)
	suite.Require().False(newList.IsEmpty())
}

func (suite *TestFilterSuite) TestNetworkIdFilter() {
	// empty paths should return all configurations
	var paths []string
	newList := filterNetworkId(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 3)
	suite.Require().False(newList.IsEmpty())

	// fetching the non-existing paths should return empty list
	paths = []string{"ideal_blockchain"}
	newList = filterNetworkId(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 0)
	suite.Require().True(newList.IsEmpty())

	// fetching the org that has one element
	paths = []string{"test_2"}
	newList = filterNetworkId(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 1)
	suite.Require().False(newList.IsEmpty())

	// fetching the org that has two element
	paths = []string{"test_1"}
	newList = filterNetworkId(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 2)
	suite.Require().False(newList.IsEmpty())

	paths = []string{"test_1", "test_2"}
	newList = filterNetworkId(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 3)
	suite.Require().False(newList.IsEmpty())

	paths = []string{"test_org"}
	newList = filterOrganization(suite.confList, paths)
	suite.Require().Equal(newList.Len(), 2)
	suite.Require().False(newList.IsEmpty())

	// fetching from new list should be successful
	paths = []string{"test_1"}
	newList = filterNetworkId(newList, paths)
	suite.Require().Equal(newList.Len(), 1)
	suite.Require().False(newList.IsEmpty())
}

func (suite *TestFilterSuite) TestConfigurationFiltering() {
	topicFilter := topic.Filter{
		Organizations: []string{"test_org"},
		NetworkIds:    []string{"test_1"},
	}
	newList := filterConfiguration(suite.confList, &topicFilter)
	suite.Require().Len(newList, 1)
}

func (suite *TestFilterSuite) TestSmartcontractFiltering() {
	topicFilter := topic.Filter{
		Organizations: []string{"test_org"},
		NetworkIds:    []string{"test_1"},
	}
	newList := filterConfiguration(suite.confList, &topicFilter)

	suite.T().Log("configs", newList[0])

	filteredSm, filteredTopics, err := filterSmartcontract(newList, suite.smList)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(filteredSm)
	suite.Require().EqualValues(suite.conf.Topic.String(topic.SmartcontractLevel), filteredTopics[0])
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFilter(t *testing.T) {
	suite.Run(t, new(TestFilterSuite))
}
