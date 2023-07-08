package smartcontract

import (
	"fmt"

	"github.com/Seascape-Foundation/mysql-seascape-extension/handler"
	"github.com/Seascape-Foundation/sds-common-lib/blockchain"
	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
	"github.com/Seascape-Foundation/sds-common-lib/smartcontract_key"
	"github.com/Seascape-Foundation/sds-service-lib/remote"
)

// Insert Inserts the smartcontract into the database
//
// Implements common/data_type/database.Crud interface
func (sm *Smartcontract) Insert(dbInterface interface{}) error {
	db := dbInterface.(*remote.ClientSocket)
	request := handler.DatabaseQueryRequest{
		Fields: []string{"network_id",
			"address",
			"abi_id",
			"transaction_id",
			"transaction_index",
			"block_number",
			"block_timestamp",
			"deployer"},
		Tables: []string{"smartcontract"},
		Arguments: []interface{}{
			sm.SmartcontractKey.NetworkId,
			sm.SmartcontractKey.Address,
			sm.AbiId,
			sm.TransactionKey.Id,
			sm.TransactionKey.Index,
			sm.BlockHeader.Number,
			sm.BlockHeader.Timestamp,
			sm.Deployer,
		},
	}
	var reply handler.InsertReply

	err := handler.INSERT.Request(db, request, &reply)
	if err != nil {
		return fmt.Errorf("handler.INSERT.Request: %w", err)
	}
	return nil
}

// SelectAll smartcontracts from database
//
// Implements common/data_type/database.Crud interface
func (sm *Smartcontract) SelectAll(dbInterface interface{}, returnValues interface{}) error {
	db := dbInterface.(*remote.ClientSocket)

	smartcontracts, ok := returnValues.(*[]*Smartcontract)
	if !ok {
		return fmt.Errorf("return_values.(*[]*Smartcontract)")
	}

	request := handler.DatabaseQueryRequest{
		Fields: []string{
			"network_id",
			"address",
			"abi_id",
			"transaction_id",
			"transaction_index",
			"block_number",
			"block_timestamp",
			"deployer",
		},
		Tables: []string{"smartcontract"},
	}
	var reply handler.SelectAllReply

	err := handler.SelectAll.Request(db, request, &reply)
	if err != nil {
		return fmt.Errorf("handler.SELECT_ALL.Request: %w", err)
	}

	*smartcontracts = make([]*Smartcontract, len(reply.Rows))

	// Loop through rows, using Scan to assign column data to struct fields.
	for i, raw := range reply.Rows {
		var sm = Smartcontract{
			SmartcontractKey: smartcontract_key.Key{},
			TransactionKey:   blockchain.TransactionKey{},
			BlockHeader:      blockchain.BlockHeader{},
		}

		err := raw.Interface(&sm.SmartcontractKey)
		if err != nil {
			return fmt.Errorf("raw.ToInterface(SmartcontractKey): %w", err)
		}

		err = raw.Interface(&sm.BlockHeader)
		if err != nil {
			return fmt.Errorf("raw.ToInterface(BlockHeader): %w", err)
		}

		err = raw.Interface(&sm.TransactionKey)
		if err != nil {
			return fmt.Errorf("raw.ToInterface(TransactionKey): %w", err)
		}

		deployer, err := raw.GetString("deployer")
		if err != nil {
			return fmt.Errorf("failed to extract deployer from database result: %w", err)
		}
		sm.Deployer = deployer

		abiId, err := raw.GetString("abi_id")
		if err != nil {
			return fmt.Errorf("failed to extract abi id from database result: %w", err)
		}
		sm.AbiId = abiId

		(*smartcontracts)[i] = &sm
	}

	returnValues = smartcontracts

	return err
}

// Select Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (sm *Smartcontract) Select(_ interface{}) error {
	return fmt.Errorf("not implemented")
}

// SelectAllByCondition Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (sm *Smartcontract) SelectAllByCondition(_ interface{}, _ key_value.KeyValue, _ interface{}) error {
	return fmt.Errorf("not implemented")
}

// Exist Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (sm *Smartcontract) Exist(_ interface{}) bool {
	return false
}

// Update Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (sm *Smartcontract) Update(_ interface{}, _ uint8) error {
	return fmt.Errorf("not implemented")
}
