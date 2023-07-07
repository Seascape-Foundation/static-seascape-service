package handler

import (
	"testing"

	"github.com/Seascape-Foundation/sds-common-lib/blockchain"
	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
	"github.com/Seascape-Foundation/sds-common-lib/smartcontract_key"
	"github.com/Seascape-Foundation/sds-service-lib/communication/message"
	"github.com/Seascape-Foundation/sds-service-lib/log"
	"github.com/Seascape-Foundation/static-seascape-service/smartcontract"
	"github.com/stretchr/testify/suite"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestSmartcontractSuite struct {
	suite.Suite
	logger log.Logger
	abi0Id string
	sm0Key smartcontract_key.Key
	sm1Key smartcontract_key.Key
	sm     smartcontract.Smartcontract
	smList *key_value.List
}

func (suite *TestSmartcontractSuite) SetupTest() {
	logger, err := log.New("test", false)
	suite.Require().NoError(err)
	suite.logger = logger

	suite.abi0Id = "hello"
	suite.sm0Key = smartcontract_key.Key{
		NetworkId: "1",
		Address:   "0xaddress",
	}
	suite.sm1Key = smartcontract_key.Key{
		NetworkId: "1",
		Address:   "0xsm_key",
	}

	sm0 := smartcontract.Smartcontract{
		SmartcontractKey: suite.sm0Key,
		AbiId:            suite.abi0Id,
		TransactionKey: blockchain.TransactionKey{
			Id:    "0x1",
			Index: 0,
		},
		BlockHeader: blockchain.BlockHeader{
			Number:    blockchain.Number(1),
			Timestamp: blockchain.Timestamp(2),
		},
		Deployer: "0xdeployer",
	}
	suite.sm = sm0

	sm1 := smartcontract.Smartcontract{
		SmartcontractKey: suite.sm1Key,
		AbiId:            suite.abi0Id,
		TransactionKey: blockchain.TransactionKey{
			Id:    "0x1",
			Index: 0,
		},
		BlockHeader: blockchain.BlockHeader{
			Number:    blockchain.Number(1),
			Timestamp: blockchain.Timestamp(2),
		},
		Deployer: "0xdeployer",
	}

	list := key_value.NewList()
	err = list.Add(sm0.SmartcontractKey, &sm0)
	suite.Require().NoError(err)

	err = list.Add(sm1.SmartcontractKey, &sm1)
	suite.Require().NoError(err)
	suite.smList = list
}

func (suite *TestSmartcontractSuite) TestGet() {
	// valid request
	validKv, err := key_value.NewFromInterface(suite.sm0Key)
	suite.Require().NoError(err)

	request := message.Request{
		Command:    "",
		Parameters: validKv,
	}
	reply := SmartcontractGet(request, suite.logger, nil, nil, suite.smList)
	suite.Require().True(reply.IsOK())

	var repliedSm GetSmartcontractReply
	err = reply.Parameters.Interface(&repliedSm)
	suite.Require().NoError(err)

	suite.Require().EqualValues(suite.sm, repliedSm)

	// request with empty parameter should fail
	request = message.Request{
		Command:    "",
		Parameters: key_value.Empty(),
	}
	reply = SmartcontractGet(request, suite.logger, nil, nil, suite.smList)
	suite.Require().False(reply.IsOK())

	// request of smartcontract that doesn't exist in the list
	// should fail
	request = message.Request{
		Command: "",
		Parameters: key_value.Empty().
			Set("network_id", "56").
			Set("address", "0xsm_key"),
	}
	reply = SmartcontractGet(request, suite.logger, nil, nil, suite.smList)
	suite.Require().False(reply.IsOK())

	// requesting with invalid type for abi id should fail
	request = message.Request{
		Command: "",
		Parameters: key_value.Empty().
			Set("network_id", 1).
			Set("address", "0xsm_key"),
	}
	reply = SmartcontractGet(request, suite.logger, nil, nil, suite.smList)
	suite.Require().False(reply.IsOK())
}

func (suite *TestSmartcontractSuite) TestSet() {
	// valid request
	validRequest := smartcontract.Smartcontract{
		SmartcontractKey: smartcontract_key.Key{
			NetworkId: "imx",
			Address:   "0xnft",
		},
		AbiId: suite.abi0Id,
		TransactionKey: blockchain.TransactionKey{
			Id:    "0x1",
			Index: 0,
		},
		BlockHeader: blockchain.BlockHeader{
			Number:    blockchain.Number(1),
			Timestamp: blockchain.Timestamp(2),
		},
		Deployer: "0xdeployer",
	}
	validKv, err := key_value.NewFromInterface(validRequest)
	suite.Require().NoError(err)

	request := message.Request{
		Command:    "",
		Parameters: validKv,
	}
	reply := SmartcontractRegister(request, suite.logger, nil, nil, suite.smList)
	suite.T().Log(reply.Message)
	suite.Require().True(reply.IsOK())

	var repliedSm SetSmartcontractReply
	err = reply.Parameters.Interface(&repliedSm)
	suite.Require().NoError(err)
	suite.Require().EqualValues(validRequest, repliedSm)

	// the abi list should have the item
	smInList, err := suite.smList.Get(repliedSm.SmartcontractKey)
	suite.Require().NoError(err)
	suite.Require().EqualValues(&repliedSm, smInList)

	// registering with empty parameter should fail
	request = message.Request{
		Command:    "",
		Parameters: key_value.Empty(),
	}
	reply = SmartcontractRegister(request, suite.logger, nil, nil, suite.smList)
	suite.Require().False(reply.IsOK())

	// request of abi that already exist in the list
	// should fail
	request = message.Request{
		Command:    "",
		Parameters: validKv,
	}
	reply = SmartcontractRegister(request, suite.logger, nil, nil, suite.smList)
	suite.Require().False(reply.IsOK())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSmartcontract(t *testing.T) {
	suite.Run(t, new(TestSmartcontractSuite))
}
