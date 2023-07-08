package configuration

import (
	"fmt"

	"github.com/Seascape-Foundation/mysql-seascape-extension/handler"
	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
	"github.com/Seascape-Foundation/sds-common-lib/topic"
	"github.com/Seascape-Foundation/sds-service-lib/remote"
)

// Insert Inserts the configuration into the database
//
// It doesn't validate the configuration.
// Call conf.Validate() before calling this
//
// Implements common/data_type/database.Crud interface
func (c *Configuration) Insert(dbInterface interface{}) error {
	db := dbInterface.(*remote.ClientSocket)
	request := handler.DatabaseQueryRequest{
		Fields:    []string{"organization", "project", "network_id", "group_name", "smartcontract_name", "address"},
		Tables:    []string{"configuration"},
		Arguments: []interface{}{c.Topic.Organization, c.Topic.Project, c.Topic.NetworkId, c.Topic.Group, c.Topic.Smartcontract, c.Address},
	}
	var reply handler.InsertReply

	err := handler.INSERT.Request(db, request, &reply)
	if err != nil {
		return fmt.Errorf("handler.INSERT.Request: %w", err)
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

	request := handler.DatabaseQueryRequest{
		Fields: []string{
			"organization as o",
			"project as p",
			"network_id as n",
			"group_name as g",
			"smartcontract_name as s",
			"address",
		},
		Tables: []string{"configuration"},
	}
	var reply handler.SelectAllReply

	err := handler.SelectAll.Request(db, request, &reply)
	if err != nil {
		return fmt.Errorf("handler.SELECT_ALL.Request: %w", err)
	}

	*configurations = make([]*Configuration, len(reply.Rows))

	// Loop through rows, using Scan to assign column data to struct fields.
	for i, raw := range reply.Rows {
		confTopic, err := topic.ParseJSON(raw)
		if err != nil {
			return fmt.Errorf("parsing topic parameters from database result failed: %w", err)
		}
		address, err := raw.GetString("address")
		if err != nil {
			return fmt.Errorf("parsing address parameter from database result failed: %w", err)
		}
		conf, err := NewFromTopic(*confTopic, address)
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
