package util

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB map[string]interface{}

// Value mengubah JSONB ke bentuk driver.Value (untuk insert ke DB)
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan membaca nilai dari DB dan mengubahnya ke bentuk map
func (j *JSONB) Scan(src interface{}) error {
	if src == nil {
		*j = nil
		return nil
	}

	switch s := src.(type) {
	case []byte:
		return json.Unmarshal(s, j)
	case string:
		return json.Unmarshal([]byte(s), j)
	default:
		return errors.New("incompatible type for JSONB")
	}
}
