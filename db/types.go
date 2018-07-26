package db

import (
	"database/sql/driver"
	"encoding/json"
)

type JSONArray []interface{}

func (t *JSONArray) Scan(v interface{}) error {
	// Should be more strictly to check this type.
	if b, ok := v.([]byte); ok {
		ret := JSONArray{}
		if err := json.Unmarshal(b, &ret); err == nil {
			*t = ret
			return nil
		} else {
			panic(err)
		}
	}
	return nil
}

func (t JSONArray) Value() (driver.Value, error) {
	if b, err := json.Marshal(t); err == nil {
		return driver.Value(string(b)), nil
	} else {
		return nil, err
	}
}

type Filterable interface {
	FilterableFields() []string
}

// 指定model哪些字段可以修改
type Updatable interface {
	UpdatableFields() []string
}

type Insertable interface {
	InsertableFields() []string
}

type Sortable interface {
	SortableFields() []string
}
