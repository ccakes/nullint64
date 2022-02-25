package nullint64

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/volatiletech/null/v9/convert"
)

// NullBytes is a global byte slice of JSON null
var NullBytes = []byte("null")

// Int64 is an nullable int64.
type Int64 struct {
	Int64 int64
	Valid bool
	Set   bool
}

// NewInt64 creates a new Int64
func NewInt64(i int64, valid bool) Int64 {
	return Int64{
		Int64: i,
		Valid: valid,
		Set:   true,
	}
}

// Int64From creates a new Int64 that will always be valid.
func Int64From(i int64) Int64 {
	return NewInt64(i, true)
}

// Int64FromPtr creates a new Int64 that be null if i is nil.
func Int64FromPtr(i *int64) Int64 {
	if i == nil {
		return NewInt64(0, false)
	}
	return NewInt64(*i, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (i Int64) IsValid() bool {
	return i.Set && i.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (i Int64) IsSet() bool {
	return i.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int64) UnmarshalJSON(data []byte) error {
	i.Set = true
	if bytes.Equal(data, NullBytes) {
		i.Valid = false
		i.Int64 = 0
		return nil
	}

	var (
		v   interface{}
		err error
	)
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch x := v.(type) {
	case float64:
		// Unmarshal again direct to int64 to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int64)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		i.Int64, err = strconv.ParseInt(str, 10, 64)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %T into Go value of type nullint64.Int64", v)
	}

	i.Valid = (err == nil) && (i.Int64 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int64) UnmarshalText(text []byte) error {
	i.Set = true
	if len(text) == 0 {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
func (i Int64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int64) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// SetValid changes this Int64's value and also sets it to be non-null.
func (i *Int64) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
	i.Set = true
}

// Ptr returns a pointer to this Int64's value, or a nil pointer if this Int64 is null.
func (i Int64) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsZero returns true for invalid Int64's, for future omitempty support (Go 1.4?)
func (i Int64) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (i *Int64) Scan(value interface{}) error {
	if value == nil {
		i.Int64, i.Valid, i.Set = 0, false, false
		return nil
	}
	i.Valid, i.Set = true, true
	return convert.ConvertAssign(&i.Int64, value)
}

// Value implements the driver Valuer interface.
func (i Int64) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.Int64, nil
}
