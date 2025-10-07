package helpers

import (
	"database/sql"
	"encoding/json"
)

// NullableString is a wrapper for sql.NullString to handle JSON marshaling.
type NullableString struct {
	sql.NullString
}

// MarshalJSON implements the json.Marshaler interface to marshal NullableString as an empty string if not valid.
func (ns NullableString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal("")
}