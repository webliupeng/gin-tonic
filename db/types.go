package db

import (
	"database/sql/driver"
	"encoding/json"
)

type JSONArray []interface{}

func (t *JSONArray) Scan(v interface{}) (err error) {
	// Should be more strictly to check this type.
	if b, ok := v.([]byte); ok {
		ret := JSONArray{}

		err = json.Unmarshal(b, &ret)

		if err == nil {
			*t = ret
		}
	}
	return
}

func (t JSONArray) Value() (val driver.Value, err error) {
	val, err = json.Marshal(t)

	return
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

type IncludableTable interface {
	IncludableTables() []string
}
