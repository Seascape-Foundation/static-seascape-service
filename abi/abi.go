// Package abi defines the abi of the smartcontract
//
// The db.go contains the database related functions of the ABI
package abi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/ahmetson/common-lib/data_type"
)

type Abi struct {
	Body string `json:"body"`
	Id   string `json:"abi_id"`
}

// ToString Returns the abi content as a string.
// The abi bytes are first formatted.
// If the abi parameters are invalid, then
// the ToString() returns empty string.
func (a *Abi) ToString() string {
	if err := a.formatBytes(); err != nil {
		return ""
	}
	return a.Body
}

// GenerateId Creates the abi hash from the abi body
// The Abi ID is the unique identifier of the abi
//
//	ID is the first 8 characters of the
//
// sha256 checksum
// representation of the abi.
//
// If the bytes field is invalid, then the id will be empty
func (a *Abi) GenerateId() error {
	a.Id = ""

	// re-serialize to remove the empty spaces
	if err := a.formatBytes(); err != nil {
		return fmt.Errorf("format_bytes: %w", err)
	}
	encoded := sha256.Sum256([]byte(a.Body))
	a.Id = hex.EncodeToString(encoded[0:8])

	return nil
}

func (a *Abi) formatBytes() error {
	// re-serialize to remove the empty spaces
	var json interface{}
	err := a.Interface(&json)
	if err != nil {
		return fmt.Errorf("failed to deserialize: %w", err)
	}
	bytes, err := data_type.Serialize(json)
	if err != nil {
		return fmt.Errorf("failed to re-serialize: %w", err)
	}
	a.Body = string(bytes)

	return nil
}

// Interface Get the interface from the bytes
// It converts the bytes into the JSON value
func (a *Abi) Interface(body interface{}) error {
	err := data_type.Deserialize([]byte(a.Body), body)
	if err != nil {
		return fmt.Errorf("data_type.Deserialize: %w", err)
	}

	return nil
}
