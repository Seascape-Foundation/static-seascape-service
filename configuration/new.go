package configuration

import (
	"fmt"

	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
	"github.com/Seascape-Foundation/sds-common-lib/topic"
)

// NewFromTopic Converts the Topic to the Configuration
// Note that you should set the address as well
func NewFromTopic(t topic.Topic, address string) (*Configuration, error) {
	c := &Configuration{
		Topic:   t,
		Address: address,
	}
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("Validate: %w", err)
	}

	return c, nil
}

// New Creates a new storage.Configuration class based on the given data
func New(parameters key_value.KeyValue) (*Configuration, error) {
	var conf Configuration
	err := parameters.Interface(&conf)
	if err != nil {
		return nil, fmt.Errorf("failed to convert key-value of Configuration to interface %v", err)
	}

	if err := conf.Validate(); err != nil {
		return nil, fmt.Errorf("Validate: %w", err)
	}

	return &conf, nil
}
