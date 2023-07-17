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
	Topic  topic.Topic   `json:"topic"`
	Topics []topic.Topic `json:"topics"`
}

func (c *Configuration) Validate() error {
	if err := c.Topic.Validate(); err != nil {
		return fmt.Errorf("Topic.Validate: %w", err)
	}
	if !c.Topic.Has("org", "proj") {
		return fmt.Errorf("topic id should missing org or proj")
	}
	if len(c.Topics) == 0 {
		return fmt.Errorf("missing Address parameter")
	}

	for _, topicId := range c.Topics {
		if !topicId.Has("org", "net", "name") {
			return fmt.Errorf("smartcontract is missing org, net and name")
		}
	}

	return nil
}
