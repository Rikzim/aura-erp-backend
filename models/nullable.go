package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// FlexInt is an int that unmarshals from either a bare JSON number (42)
// or a quoted JSON string ("42"). The frontend sends select values as
// strings, so plain int fields like client_id need this flexibility.
type FlexInt int

func (f FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(f))
}

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" || s == `""` || s == "" {
		*f = 0
		return nil
	}
	// Bare number
	if s[0] != '"' {
		var n int
		if err := json.Unmarshal(data, &n); err != nil {
			return err
		}
		*f = FlexInt(n)
		return nil
	}
	// Quoted string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		*f = 0
		return nil
	}
	n, err := strconv.Atoi(str)
	if err != nil {
		return fmt.Errorf("FlexInt: cannot parse %q as integer", str)
	}
	*f = FlexInt(n)
	return nil
}

// NullInt64 wraps sql.NullInt64 with proper JSON marshaling (null vs number)
type NullInt64 struct {
	sql.NullInt64
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Int64)
}

func (n *NullInt64) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" || s == `""` || s == "" {
		n.Valid = false
		return nil
	}
	// Bare number: 42
	if len(s) > 0 && s[0] != '"' {
		n.Valid = true
		return json.Unmarshal(data, &n.Int64)
	}
	// Quoted string sent by the frontend select: "42"
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		n.Valid = false
		return nil
	}
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("NullInt64: cannot parse %q as integer", str)
	}
	n.Valid = true
	n.Int64 = v
	return nil
}

// NullString wraps sql.NullString with proper JSON marshaling (null vs string)
type NullString struct {
	sql.NullString
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.String)
}

func (n *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(data, &n.String)
}

// NullTime wraps sql.NullTime with proper JSON marshaling (null vs RFC3339 string)
type NullTime struct {
	sql.NullTime
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Time.Format(time.RFC3339))
}

func (n *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	// Accept both full RFC3339 ("2026-03-03T00:00:00Z") and plain date ("2026-03-03")
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, s); err == nil {
			n.Valid = true
			n.Time = t
			return nil
		}
	}
	return fmt.Errorf("cannot parse %q as a date/time value", s)
}
