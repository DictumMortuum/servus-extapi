package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type JsonNullString sql.NullFloat64

func NewJsonNullString(i float64) JsonNullString {
	return JsonNullString{
		Float64: i,
		Valid:   true,
	}
}

func (obj *JsonNullString) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		obj.Valid = true
		tmp := strings.Trim(string(v), "\"")
		f, _ := strconv.ParseFloat(tmp, 64)
		obj.Float64 = f
		return nil
	case nil:
		obj.Valid = false
		obj.Float64 = 0.0
		return nil
	case string:
		obj.Valid = true
		tmp := strings.Trim(v, "\"")
		f, _ := strconv.ParseFloat(tmp, 64)
		obj.Float64 = f
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (obj JsonNullString) Value() (driver.Value, error) {
	if !obj.Valid {
		return nil, nil
	}

	return obj.Float64, nil
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullString) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x != nil {
		v.Valid = true
		v.Float64 = *x
	} else {
		v.Valid = false
	}
	return nil
}
