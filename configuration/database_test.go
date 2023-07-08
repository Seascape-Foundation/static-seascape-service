package configuration

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	database "github.com/Seascape-Foundation/mysql-seascape-extension"
	"github.com/Seascape-Foundation/sds-common-lib/blockchain"
	"github.com/Seascape-Foundation/sds-common-lib/smartcontract_key"
	"github.com/Seascape-Foundation/sds-common-lib/topic"
	"github.com/Seascape-Foundation/sds-service-lib/configuration"
	parameter "github.com/Seascape-Foundation/sds-service-lib/identity"
	"github.com/Seascape-Foundation/sds-service-lib/log"
	"github.com/Seascape-Foundation/sds-service-lib/remote"
	"github.com/Seascape-Foundation/static-seascape-service/abi"
	"github.com/Seascape-Foundation/static-seascape-service/smartcontract"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestConfigurationDbSuite struct {
	suite.Suite
	dbName        string
	configuration Configuration
	container     *mysql.MySQLContainer
	dbCon         *remote.ClientSocket
	ctx           context.Context
}

func (suite *TestConfigurationDbSuite) SetupTest() {
	// prepare the database creation
	suite.dbName = "test"
	_, filename, _, _ := runtime.Caller(0)
	// configuration depends on smartcontract.
	// smartcontract depends on abi.
	fileDir := filepath.Dir(filename)
	storageAbi := "20230308171023_storage_abi.sql"
	storageSmartcontract := "20230308173919_storage_smartcontract.sql"
	storageConfiguration := "20230308173943_storage_configuration.sql"
	changeGroupType := "20230314150414_storage_configuration_group_type.sql"

	abiSqlPath := filepath.Join(fileDir, "..", "..", "_db", "migrations", storageAbi)
	smartcontractSqlPath := filepath.Join(fileDir, "..", "..", "_db", "migrations", storageSmartcontract)
	configurationSqlPath := filepath.Join(fileDir, "..", "..", "_db", "migrations", storageConfiguration)
	changeGroupPath := filepath.Join(fileDir, "..", "..", "_db", "migrations", changeGroupType)

	suite.T().Log("the configuration table path", configurationSqlPath)

	// run the container
	ctx := context.TODO()
	container, err := mysql.RunContainer(ctx,
		mysql.WithDatabase(suite.dbName),
		mysql.WithUsername("root"),
		mysql.WithPassword("tiger"),
		mysql.WithScripts(abiSqlPath, smartcontractSqlPath, configurationSqlPath, changeGroupPath),
	)
	suite.Require().NoError(err)
	suite.container = container
	suite.ctx = ctx

	logger, err := log.New("mysql-suite", false)
	suite.Require().NoError(err)
	appConfig, err := configuration.NewAppConfig(logger)
	suite.Require().NoError(err)

	// Creating a database client
	// after settings the default parameters
	// we should have the user name and password
	appConfig.SetDefaults(database.DatabaseConfigurations)

	// Overwrite the default parameters to use test container
	host, err := container.Host(ctx)
	suite.Require().NoError(err)
	ports, err := container.Ports(ctx)
	suite.Require().NoError(err)
	exposedPort := ports["3306/tcp"][0].HostPort

	database.DatabaseConfigurations.Parameters["SDS_DATABASE_HOST"] = host
	database.DatabaseConfigurations.Parameters["SDS_DATABASE_PORT"] = exposedPort
	database.DatabaseConfigurations.Parameters["SDS_DATABASE_NAME"] = suite.dbName

	//go database.Run(appConfig, logger)
	// wait for initiation of the controller
	time.Sleep(time.Second * 1)

	database_service, err := parameter.Inprocess("database")
	suite.Require().NoError(err)
	client, err := remote.InprocRequestSocket(database_service.Url(), logger, appConfig)
	suite.Require().NoError(err)

	suite.dbCon = client

	// add the storage abi
	abiId := "base64="
	sampleAbi := abi.Abi{
		Body: []byte("[{}]"),
		Id:   abiId,
	}
	err = sampleAbi.Insert(suite.dbCon)
	suite.Require().NoError(err)

	// add the storage smartcontract
	key, _ := smartcontract_key.New("1", "0xaddress")
	txKey := blockchain.TransactionKey{
		Id:    "0xtx_id",
		Index: 0,
	}
	header, _ := blockchain.NewHeader(uint64(1), uint64(23))
	deployer := "0xahmetson"

	sm := smartcontract.Smartcontract{
		SmartcontractKey: key,
		AbiId:            abiId,
		TransactionKey:   txKey,
		BlockHeader:      header,
		Deployer:         deployer,
	}
	err = sm.Insert(suite.dbCon)
	suite.Require().NoError(err)

	sample := topic.Topic{
		Organization:  "seascape",
		Project:       "sds-core",
		NetworkId:     "1",
		Group:         "test-suite",
		Smartcontract: "TestErc20",
	}
	suite.configuration = Configuration{
		Topic:   sample,
		Address: key.Address,
	}

	suite.T().Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			suite.T().Fatalf("failed to terminate container: %s", err)
		}
		if err := suite.dbCon.Close(); err != nil {
			suite.T().Fatalf("failed to terminate database connection: %s", err)
		}
	})
}

func (suite *TestConfigurationDbSuite) TestConfiguration() {
	var configs []*Configuration

	err := suite.configuration.SelectAll(suite.dbCon, &configs)
	suite.Require().NoError(err)
	suite.Require().Len(configs, 0)

	err = suite.configuration.Insert(suite.dbCon)
	suite.Require().NoError(err)

	err = suite.configuration.SelectAll(suite.dbCon, &configs)
	suite.Require().NoError(err)
	suite.Require().Len(configs, 1)
	suite.Require().EqualValues(suite.configuration, *configs[0])

	// inserting a configuration
	// that links to the non existing smartcontract
	// should fail
	sample := topic.Topic{
		Organization:  "seascape",
		Project:       "sds-core",
		NetworkId:     "1",
		Group:         "test-suite",
		Smartcontract: "TestToken",
	}
	conf := Configuration{
		Topic:   sample,
		Address: "not_inserted",
	}
	err = conf.Insert(suite.dbCon)
	suite.Require().Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestConfigurationDb(t *testing.T) {
	suite.Run(t, new(TestConfigurationDbSuite))
}
