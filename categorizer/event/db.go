package event

import (
	"encoding/json"
	"fmt"

	"github.com/blocklords/sds/common/blockchain"
	"github.com/blocklords/sds/common/smartcontract_key"
	"github.com/blocklords/sds/db"
)

func Save(db *db.Database, t *Log) error {
	byt, err := t.Parameters.ToBytes()
	if err != nil {
		return fmt.Errorf("event.Parameters.ToBytes %v: %w", t.Parameters, err)
	}

	_, err = db.Connection.Exec(`INSERT IGNORE INTO categorizer_event 
	(address, transaction_id, transaction_index, network_id, block_number, block_timestamp, log_index, event_name, event_parameters)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, t.SmartcontractKey.Address, t.TransactionKey.Id, t.TransactionKey.Index, t.SmartcontractKey.NetworkId, t.BlockHeader.Number, t.BlockHeader.Timestamp, t.Index, t.Name, byt)

	if err != nil {
		return fmt.Errorf("database exec: %w", err)
	}

	return nil
}

// returns list of logs for smartcontracts
func GetLogsFromDb(con *db.Database, smartcontracts []smartcontract_key.Key, block_timestamp blockchain.Timestamp, limit uint64) ([]*Log, error) {
	var logs []*Log = make([]*Log, 0)
	sm_amount := len(smartcontracts)

	if sm_amount == 0 {
		return logs, nil
	}

	args := make([]interface{}, (sm_amount*2)+2)
	offset := 0
	args[offset] = block_timestamp
	offset++

	smartcontracts_clause := ""
	for i, sm := range smartcontracts {
		network_id := sm.NetworkId
		address := sm.Address

		smartcontracts_clause += "(network_id = ? AND address = ?) "
		if i < sm_amount-1 {
			smartcontracts_clause += " OR "
		}

		args[offset] = network_id
		offset++
		args[offset] = address
		offset++
	}
	args[offset] = limit

	query := `
	SELECT
		block_number, 
		block_timestamp,
		transaction_id,
		transaction_index,
		log_index,
		address,
		network_id,
		event_name,
		event_parameters
	FROM 
		categorizer_event
	WHERE 
		block_timestamp >= ? AND ` + smartcontracts_clause + " LIMIT ? "

	rows, err := con.Connection.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query: %w", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var s Log
		var output_bytes []byte
		if err := rows.Scan(&s.BlockHeader.Number, &s.BlockHeader.Timestamp, &s.TransactionKey.Id, &s.TransactionKey.Index, &s.Index, &s.SmartcontractKey.Address, &s.SmartcontractKey.NetworkId, &s.Name, &output_bytes); err != nil {
			return nil, fmt.Errorf("database row scan: %w", err)
		}

		jsonErr := json.Unmarshal(output_bytes, &s.Parameters)
		if jsonErr != nil {
			return nil, fmt.Errorf("json.deserialize %s: %w", string(output_bytes), err)
		}

		logs = append(logs, &s)
	}

	return logs, nil
}
