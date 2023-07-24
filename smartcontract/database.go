package smartcontract

import (
	"fmt"
	"github.com/ahmetson/common-lib/topic"

	"github.com/ahmetson/common-lib/data_type/key_value"
	databaseExtension "github.com/ahmetson/service-lib/extension/database"
	"github.com/ahmetson/service-lib/remote"
)

// Insert Inserts the smartcontract into the database
//
// Implements common/data_type/database.Crud interface
func (sm *Smartcontract) Insert(dbInterface interface{}) error {
	db := dbInterface.(*remote.ClientSocket)
	request := databaseExtension.QueryRequest{
		Fields: []string{
			"id",
			"topic",
			"transaction_id",
			"owner",
			"verifier",
			"specific",
		},
		Tables: []string{"smartcontract"},
		Arguments: []interface{}{
			sm.Topic.Id().Only("org", "net", "name"),
			sm.Topic,
			sm.TransactionId,
			sm.Owner,
			sm.Verifier,
			sm.Specific,
		},
	}.Request(databaseExtension.Insert)
	var reply databaseExtension.InsertReply

	parameters, err := db.RequestRemoteService(&request)
	if err != nil {
		return fmt.Errorf("databaseExtension.INSERT.Request: %w", err)
	}

	err = parameters.Interface(&reply)
	if err != nil {
		return fmt.Errorf("failed to parse reply: %w", err)
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

	request := databaseExtension.QueryRequest{
		Fields: []string{
			"topic",
			"transaction_id",
			"owner",
			"verifier",
			"specific",
		},
		Tables: []string{"smartcontract"},
	}.Request(databaseExtension.SelectAll)

	parameters, err := db.RequestRemoteService(&request)
	if err != nil {
		return fmt.Errorf("databaseExtension.SELECT_ALL.Request: %w", err)
	}

	var reply databaseExtension.SelectAllReply
	err = parameters.Interface(&reply)
	if err != nil {
		return fmt.Errorf("failed to parse reply")
	}

	*smartcontracts = make([]*Smartcontract, len(reply.Rows))

	// Loop through rows, using Scan to assign column data to struct fields.
	for i, raw := range reply.Rows {
		var sm = Smartcontract{
			Topic:         topic.Topic{},
			TransactionId: "",
			Owner:         "",
			Verifier:      "",
			Specific:      key_value.Empty(),
		}

		topicId, err := raw.GetString("topic_id")
		if err != nil {
			return fmt.Errorf("failed to extract topic_id from database result: %w", err)
		}
		sm.Topic, err = topic.Id(topicId).Unmarshal()
		if err != nil {
			return fmt.Errorf("failed to decode data into topic")
		}

		sm.Specific, err = raw.GetKeyValue("specific")
		if err != nil {
			return fmt.Errorf("raw.GetKeyValue(specific): %w", err)
		}

		owner, err := raw.GetString("owner")
		if err != nil {
			return fmt.Errorf("failed to extract owner from database result: %w", err)
		}
		sm.Owner = owner

		verifier, err := raw.GetString("verifier")
		if err != nil {
			return fmt.Errorf("failed to extract verifier from database result: %w", err)
		}
		sm.Verifier = verifier

		transactionId, err := raw.GetString("transaction_id")
		if err != nil {
			return fmt.Errorf("failed to extract transaction_id from database result: %w", err)
		}
		sm.TransactionId = transactionId

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
