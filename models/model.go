package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// NullToEmptyString :
type NullToEmptyString string

// Value : overide default mysql Value
func (n NullToEmptyString) Value() (driver.Value, error) {
	return string(n), nil
}

// Scan : overide default mysql Scan
func (n *NullToEmptyString) Scan(v interface{}) error {
	if v == nil {
		*n = ""
		return nil
	}

	s, ok := v.([]uint8)
	if !ok {
		return fmt.Errorf("Cannot scan value %v: expected string type, got %T", v, v)
	}

	*n = NullToEmptyString(s)
	return nil
}

// Uint8ArrayToArrayString :
type Uint8ArrayToArrayString []string

// Value : overide default mysql Value
func (n Uint8ArrayToArrayString) Value() (driver.Value, error) {
	s, err := json.Marshal(n)

	if nil != err {
		return "", err
	}

	if nil == s {
		return "", nil
	}

	return string(s), nil
}

// Scan : overide default mysql Scan
func (n *Uint8ArrayToArrayString) Scan(v interface{}) error {
	if v == nil {
		*n = []string{}
		return nil
	}

	val, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("Cannot scan value %v: expected string type, got %T", v, v)
	}

	var value []string

	if 0 != len(val) {
		if err := json.Unmarshal(val, &value); nil != err {
			return fmt.Errorf("Cannot scan value %v: expected string type, got %T", v, v)
		}
	}

	*n = Uint8ArrayToArrayString(value)
	return nil
}
