package abi

import (
	"fmt"
	"github.com/ahmetson/common-lib/data_type"
	"github.com/ahmetson/common-lib/data_type/key_value"
	databaseExtension "github.com/ahmetson/service-lib/extension/database"
	"github.com/ahmetson/service-lib/remote"
)

// Insert into database
//
// Implements common/data_type/database.Crud interface
func (a *Abi) Insert(dbInterface interface{}) error {
	db := dbInterface.(*remote.ClientSocket)
	request := databaseExtension.QueryRequest{
		Fields:    []string{"abi_id", "body"},
		Tables:    []string{"abi"},
		Arguments: []interface{}{a.Id, data_type.AddJsonPrefix([]byte(a.Body))},
	}
	var reply databaseExtension.InsertReply

	err := databaseExtension.Insert.Request(db, request, &reply)
	if err != nil {
		return fmt.Errorf("databaseExtension.Insert.Request: %w", err)
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

	request := databaseExtension.QueryRequest{
		Fields: []string{"abi_id as id", "body as bytes"},
		Tables: []string{"storage_abi"},
	}
	var reply databaseExtension.SelectAllReply

	err := databaseExtension.SelectAll.Request(dbClient, request, &reply)
	if err != nil {
		return fmt.Errorf("databaseExtension.SELECT_ALL.Push: %w", err)
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
func (a *Abi) Select(dbInterface interface{}) error {
	dbClient := dbInterface.(*remote.ClientSocket)

	request := databaseExtension.QueryRequest{
		Where:     "abi_id=?",
		Tables:    []string{"abi"},
		Arguments: []interface{}{a.Id},
	}
	var reply databaseExtension.SelectRowReply

	err := databaseExtension.SelectRow.Request(dbClient, request, &reply)
	if err != nil {
		return fmt.Errorf("databaseExtension.SELECT_ROW.Push: %w", err)
	}

	abi, err := New(reply.Outputs)
	if err != nil {
		return fmt.Errorf("failed to parse the database reply into Abi: %w", err)
	}
	a.Id = abi.Id
	a.Body = abi.Body

	return err
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
