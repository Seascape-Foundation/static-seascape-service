package configuration

import (
	"fmt"

	"github.com/ahmetson/common-lib/data_type/key_value"
	"github.com/ahmetson/common-lib/topic"
	databaseExtension "github.com/ahmetson/service-lib/extension/database"
	"github.com/ahmetson/service-lib/remote"
)

// Insert Inserts the configuration into the database
//
// It doesn't validate the configuration.
// Call conf.Validate() before calling this
//
// Implements common/data_type/database.Crud interface
func (c *Configuration) Insert(dbInterface interface{}) error {
	ids := make([]topic.Id, len(c.Topics))

	for i, sm := range c.Topics {
		ids[i] = sm.Id()
	}

	db := dbInterface.(*remote.ClientSocket)
	request := databaseExtension.QueryRequest{
		Fields:    []string{"id", "smartcontracts"},
		Tables:    []string{"configuration"},
		Arguments: []interface{}{c.Topic.Id().Only("org", "proj"), ids},
	}
	var reply databaseExtension.InsertReply

	err := databaseExtension.Insert.Request(db, request, &reply)
	if err != nil {
		return fmt.Errorf("databaseExtension.INSERT.Request: %w", err)
	}
	return nil
}

// SelectAll configurations from database
//
// Implements common/data_type/database.Crud interface
func (c *Configuration) SelectAll(dbInterface interface{}, returnValues interface{}) error {
	db := dbInterface.(*remote.ClientSocket)

	configurations, ok := returnValues.(*[]*Configuration)
	if !ok {
		return fmt.Errorf("return_values.(*[]*Configuration)")
	}

	request := databaseExtension.QueryRequest{
		Fields: []string{
			"id",
			"smartcontracts",
		},
		Tables: []string{"configuration"},
	}
	var reply databaseExtension.SelectAllReply

	err := databaseExtension.SelectAll.Request(db, request, &reply)
	if err != nil {
		return fmt.Errorf("databaseExtension.SELECT_ALL.Request: %w", err)
	}

	*configurations = make([]*Configuration, len(reply.Rows))

	// Loop through rows, using Scan to assign column data to struct fields.
	for i, raw := range reply.Rows {
		idString, err := raw.GetString("id")
		if err != nil {
			return fmt.Errorf("parsing topic parameters from database result failed: %w", err)
		}
		confId := topic.Id(idString)
		confTopic, err := confId.Unmarshal()
		if err != nil {
			return fmt.Errorf("failed to convert configuration id %s to topic: %w", confId, err)
		}

		idStrings, err := raw.GetStringList("smartcontracts")
		if err != nil {
			return fmt.Errorf("parsing address parameter from database result failed: %w", err)
		}
		smartcontracts := make([]topic.Topic, len(idStrings))

		for i, idString := range idStrings {
			smartcontractId := topic.Id(idString)
			smartcontractTopic, err := smartcontractId.Unmarshal()
			if err != nil {
				return fmt.Errorf("failed to convert smartcontract id %s to topic in conf %s: %w", smartcontractId, confId, err)
			}
			smartcontracts[i] = smartcontractTopic
		}

		conf, err := NewFromTopic(confTopic, smartcontracts)
		if err != nil {
			return fmt.Errorf("NewFromTopic: %w", err)
		}
		(*configurations)[i] = conf
	}
	returnValues = configurations

	return err
}

// Select Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (c *Configuration) Select(_ interface{}) error {
	return fmt.Errorf("not implemented")
}

// SelectAllByCondition Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (c *Configuration) SelectAllByCondition(_ interface{}, _ key_value.KeyValue, _ interface{}) error {
	return fmt.Errorf("not implemented")
}

// Exist Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (c *Configuration) Exist(_ interface{}) bool {
	return false
}

// Update Not implemented common/data_type/database.Crud interface
//
// Returns an error
func (c *Configuration) Update(_ interface{}, _ uint8) error {
	return fmt.Errorf("not implemented")
}
