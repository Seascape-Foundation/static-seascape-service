// Package configuration defines the link between smartcontract and topic.
package configuration

import (
	"fmt"

	"github.com/Seascape-Foundation/sds-common-lib/topic"
)

// Configuration The Storage Configuration is the relationship
// between the topic and the smartcontract.
//
// The database part depends on the Storage Smartcontract
type Configuration struct {
	Id             topic.Topic   `json:"id"`
	Smartcontracts []topic.Topic `json:"smartcontract"`
}

func (c *Configuration) Validate() error {
	if err := c.Id.Validate(); err != nil {
		return fmt.Errorf("Topic.Validate: %w", err)
	}
	if !c.Id.Has("org", "proj") {
		return fmt.Errorf("topic id should missing org or proj")
	}
	if len(c.Smartcontracts) == 0 {
		return fmt.Errorf("missing Address parameter")
	}

	for _, topicId := range c.Smartcontracts {
		if !topicId.Has("org", "net", "name") {
			return fmt.Errorf("smartcontract is missing org, net and name")
		}
	}

	return nil
}
