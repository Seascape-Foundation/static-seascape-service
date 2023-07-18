// Package abi The new.go keeps the functions that creates a new Abi from given parameters
package abi

import (
	"fmt"

	"github.com/ahmetson/common-lib/data_type"
	"github.com/ahmetson/common-lib/data_type/key_value"
)

// New Wraps the JSON abi interface to the internal data type.
// It's blockchain agnostic.
func New(kv key_value.KeyValue) (*Abi, error) {
	var abi Abi
	id, err := kv.GetString("abi_id")
	if err != nil {
		return nil, fmt.Errorf("key_value.GetString(id): %w", err)
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("missing `id` parameter")
	} else {
		abi.Id = id
	}
	body, err := kv.GetString("body")
	if err != nil {
		return nil, fmt.Errorf("key_value.GetString(bytes): %w", err)
	}

	abi.Body = body

	return &abi, nil
}

// NewFromInterface The bytes data are given as a JSON
// It will generate ID.
func NewFromInterface(body interface{}) (*Abi, error) {
	bytes, err := data_type.Serialize(body)
	if err != nil {
		return nil, err
	}
	return NewFromBytes(bytes)
}

// NewFromBytes creates the Abi data based on the JSON string. This function calculates the abi hash
// but won't set it in the database.
func NewFromBytes(bytes []byte) (*Abi, error) {
	abi := Abi{Body: string(bytes)}
	err := abi.GenerateId()
	if err != nil {
		return nil, fmt.Errorf("GenerateId: %w", err)
	}

	return &abi, nil
}
