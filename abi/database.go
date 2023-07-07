package abi

import (
	"fmt"

	"github.com/Seascape-Foundation/mysql-seascape-extension/handler"
	"github.com/Seascape-Foundation/sds-common-lib/data_type"
	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
	"github.com/Seascape-Foundation/sds-service-lib/remote"
)

// Insert into database
//
// Implements common/data_type/database.Crud interface
func (a *Abi) Insert(dbInterface interface{}) error {
	db := dbInterface.(*remote.ClientSocket)
	request := handler.DatabaseQueryRequest{
		Fields:    []string{"abi_id", "body"},
		Tables:    []string{"abi"},
		Arguments: []interface{}{a.Id, data_type.AddJsonPrefix(a.Bytes)},
	}
	var reply handler.InsertReply

	err := handler.INSERT.Request(db, request, &reply)
	if err != nil {
		return fmt.Errorf("handler.INSERT.Request: %w", err)
	}
	return nil
}

// SelectAll abi from database
//
// Implements common/data_type/database.Crud interface
func (a *Abi) SelectAll(dbInterface interface{}, returnValues interface{}) error {
	dbClient := dbInterface.(*remote.ClientSocket)
	abis, ok := returnValues.(*[]*Abi)
	if !ok {
		return fmt.Errorf("return_values.(*[]*Abi)")
	}

	request := handler.DatabaseQueryRequest{
		Fields: []string{"abi_id as id", "body as bytes"},
		Tables: []string{"storage_abi"},
	}
	var reply handler.SelectAllReply

	err := handler.SelectAll.Request(dbClient, request, &reply)
	if err != nil {
		return fmt.Errorf("handler.SELECT_ALL.Push: %w", err)
	}
	*abis = make([]*Abi, len(reply.Rows))

	// Loop through rows, using Scan to assign column data to struct fields.
	for i, raw := range reply.Rows {
		abi, err := New(raw)
		if err != nil {
			return fmt.Errorf("new Abi from database result: %w", err)
		}
		(*abis)[i] = abi
	}
	returnValues = abis

	return err
}

// Select Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (a *Abi) Select(_ interface{}, _ interface{}) error {
	return fmt.Errorf("not implemented")
}

// SelectAllByCondition Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (a *Abi) SelectAllByCondition(_ interface{}, _ key_value.KeyValue, _ interface{}) error {
	return fmt.Errorf("not implemented")
}

// Exist Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (a *Abi) Exist(_ interface{}) bool {
	return false
}

// Update Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (a *Abi) Update(_ interface{}, _ uint8) error {
	return fmt.Errorf("not implemented")
}
