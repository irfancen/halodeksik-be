package appdb

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func NewNullInt64(value ...int64) *NullInt64 {
	if len(value) > 0 {
		return &NullInt64{
			sql.NullInt64{Int64: value[0], Valid: true},
		}
	}
	return &NullInt64{
		sql.NullInt64{Valid: false},
	}
}

func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}
