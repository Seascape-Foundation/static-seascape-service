package smartcontract

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	db "github.com/Seascape-Foundation/mysql-seascape-extension"
	"github.com/ahmetson/common-lib/blockchain"
	"github.com/ahmetson/common-lib/smartcontract_key"
	"github.com/ahmetson/service-lib/configuration"
	parameter "github.com/ahmetson/service-lib/identity"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/service-lib/remote"
	"github.com/ahmetson/static-service/abi"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestSmartcontractDbSuite struct {
	suite.Suite
	dbName        string
	smartcontract Smartcontract
	container     *mysql.MySQLContainer
	dbCon         *remote.ClientSocket
	ctx           context.Context
}

func (suite *TestSmartcontractDbSuite) SetupTest() {
	// prepare the database creation
	suite.dbName = "test"
	_, filename, _, _ := runtime.Caller(0)
	storageAbi := "20230308171023_storage_abi.sql"
	storageSmartcontract := "20230308173919_storage_smartcontract.sql"
	abiSqlPath := filepath.Join(filepath.Dir(filename), "..", "..", "_db", "migrations", storageAbi)
	smartcontractSqlPath := filepath.Join(filepath.Dir(filename), "..", "..", "_db", "migrations", storageSmartcontract)
	suite.T().Log("storage smartcontract sql table path", smartcontractSqlPath)

	// run the container
	ctx := context.TODO()
	container, err := mysql.RunContainer(ctx,
		mysql.WithDatabase(suite.dbName),
		mysql.WithUsername("root"),
		mysql.WithPassword("tiger"),
		mysql.WithScripts(abiSqlPath, smartcontractSqlPath),
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
	// we should have the username and password

	// Overwrite the default parameters to use test container
	host, err := container.Host(ctx)
	suite.Require().NoError(err)
	ports, err := container.Ports(ctx)
	suite.Require().NoError(err)
	exposedPort := ports["3306/tcp"][0].HostPort

	db.DatabaseConfigurations.Parameters["SDS_DATABASE_HOST"] = host
	db.DatabaseConfigurations.Parameters["SDS_DATABASE_PORT"] = exposedPort
	db.DatabaseConfigurations.Parameters["SDS_DATABASE_NAME"] = suite.dbName

	// wait for initiation of the controller
	time.Sleep(time.Second * 1)

	databaseService, err := parameter.Inprocess("db")
	suite.Require().NoError(err)
	client, err := remote.InprocRequestSocket(databaseService.Url(), logger, appConfig)
	suite.Require().NoError(err)

	suite.dbCon = client

	_, _ = smartcontract_key.New("1", "0xaddress")
	_, _ = blockchain.NewHeader(uint64(1), uint64(23))
	suite.smartcontract = Smartcontract{}

	suite.T().Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			suite.T().Fatalf("failed to terminate container: %s", err)
		}
		if err := suite.dbCon.Close(); err != nil {
			suite.T().Fatalf("failed to terminate database connection: %s", err)
		}
	})
}

func (suite *TestSmartcontractDbSuite) TestSmartcontract() {
	var smartcontracts []*Smartcontract
	err := suite.smartcontract.SelectAll(suite.dbCon, &smartcontracts)
	suite.Require().NoError(err)
	suite.Require().Len(smartcontracts, 0)

	// Insert into the database
	// it should fail, since the smartcontract depends on the
	// abi
	err = suite.smartcontract.Insert(suite.dbCon)
	suite.Require().Error(err)

	sampleAbi := abi.Abi{
		Body: "[{}]",
		Id:   "",
	}
	err = sampleAbi.Insert(suite.dbCon)
	suite.Require().NoError(err)

	// inserting a smartcontract should be successful
	err = suite.smartcontract.Insert(suite.dbCon)
	suite.Require().NoError(err)

	// duplicate key in the database
	// it should fail
	err = suite.smartcontract.Insert(suite.dbCon)
	suite.Require().Error(err)

	// all from database
	err = suite.smartcontract.SelectAll(suite.dbCon, &smartcontracts)
	suite.Require().NoError(err)
	suite.Require().Len(smartcontracts, 1)
	suite.Require().EqualValues(suite.smartcontract, *smartcontracts[0])
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSmartcontractDb(t *testing.T) {
	suite.Run(t, new(TestSmartcontractDbSuite))
}
