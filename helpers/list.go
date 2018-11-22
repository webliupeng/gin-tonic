package helpers

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/webliupeng/gin-tonic/db"
)

type condition struct {
	Name     string
	Oparator string
}

// CriteriaCreator - 条件构造
type CriteriaCreator func(db *gorm.DB, c *gin.Context) *gorm.DB

func CriteriaByParam(key string) CriteriaCreator {
	f := func(db *gorm.DB, c *gin.Context) *gorm.DB {
		if val, ok := c.Params.Get(key); ok {
			return db.Where(key+" = ?", val)
		}
		return db
	}
	return f
}

// func BuildPaginion(db *gorm.DB, c *gin.Context) {

//  pageSize, err := strconv.Atoi(c.DefaultQuery(".maxResults", "10"))
//  offset, err2 := strconv.Atoi(c.DefaultQuery(".offset", "0"))
//  if pageSize > 1000 {
//   panic("page size too large")
//  }
//  if err != nil {
//   panic(err)
//  }

//  if err2 != nil {
//   panic(err2)
//  }

// }

func BuildQueryDB(modelIns interface{}, c *gin.Context) (*gorm.DB, error) {
	myType := reflect.TypeOf(modelIns)

	slice := reflect.MakeSlice(reflect.SliceOf(myType), 0, 0)
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)

	expressions := []string{}
	values := []interface{}{}

	filterableFields := []string{}
	if fb, ok := modelIns.(db.Filterable); ok {
		filterableFields = fb.FilterableFields()
	}

	oparators := map[string]func(key, val string) (string, interface{}){
		"lt": func(key, val string) (string, interface{}) {
			return key + " < ?", val
		},
		"lte": func(key, val string) (string, interface{}) {
			return key + " <= ?", val
		},
		"gt": func(key, val string) (string, interface{}) {
			return key + " > ?", val
		},
		"gte": func(key, val string) (string, interface{}) {
			return key + " >= ?", val
		},
		"in": func(key, val string) (string, interface{}) {
			return key + " IN (?)", strings.Split(val, ",")
		},
		"not": func(key, val string) (string, interface{}) {
			return key + " NOT IN (?)", val
		},
		"like": func(key, val string) (string, interface{}) {
			return key + " LIKE ?", strings.Replace(val, "*", "%", -1)
		},
	}

	for _, field := range filterableFields {
		if fieldValue := c.Query(field); len(fieldValue) > 0 {
			expressions = append(expressions, field+"= ?")
			values = append(values, fieldValue)
		} else {
			for opr, handle := range oparators {
				if val := c.Query(fmt.Sprintf("%v_%v", field, opr)); val != "" {
					expression, qv := handle(field, val)
					expressions = append(expressions, expression)
					values = append(values, qv)
				}
			}
		}
	}

	query := db.DB().Where(strings.Join(expressions, " AND "), values...)
	if includes := c.Query(".includes"); includes != "" {
		if includable, ok := modelIns.(db.IncludableTable); ok {
			tbs := map[string]bool{}
			for _, val := range includable.IncludableTables() {
				tbs[val] = true
			}
			for _, table := range strings.Split(includes, ",") {
				if tbs[table] {
					query = query.Preload(strings.Title(table))
				} else {
					return nil, errors.New("Can not includes " + table)
				}
			}
		} else {
			return nil, errors.New("Can not includes a non-includable model")
		}
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery(".maxResults", "10"))
	offset, err2 := strconv.Atoi(c.DefaultQuery(".offset", "0"))
	if pageSize > 1000 {
		panic("page size too large")
	}
	if err != nil {
		panic(err)
	}

	if err2 != nil {
		panic(err2)
	}

	if sortable, ok := modelIns.(db.Sortable); ok {
		fields := sortable.SortableFields()
		sortableFields := map[string]interface{}{}

		for _, val := range fields {
			sortableFields[val] = true
		}

		orderby := c.DefaultQuery(".sort", "id")
		orderField := orderby
		isDesc := false
		if orderby[0:1] == "-" {
			orderField = orderby[1:]
			isDesc = true
		}

		//fmt.Println("sortable", fields, orderField)

		if sortableFields[orderField] != nil {
			if isDesc {
				query = query.Order(orderField + " DESC")
			} else {
				query = query.Order(orderField + " ASC")
			}
		}
	}

	return query.Limit(pageSize).Offset(offset), nil
}

// ListHandlerWithoutServe -
// Deprecated: use BuildQueryDB to instead
func ListHandlerWithoutServe(modelIns interface{}, c *gin.Context, paramCreators ...CriteriaCreator) (int, interface{}) {
	myType := reflect.TypeOf(modelIns)
	slice := reflect.MakeSlice(reflect.SliceOf(myType), 0, 0)
	x := reflect.New(slice.Type())
	if query, err := BuildQueryDB(modelIns, c); err == nil {
		for _, paramCreator := range paramCreators {
			query = paramCreator(query, c)
		}

		query.Find(x.Interface())

		var total int
		query.Model(modelIns).Count(&total)

		//total, data := ListHandlerWithoutServe(modelIns, c, paramCreators...)
		return total, x.Interface()
	} else {
		ErrorResponse(c, 400, err.Error())
	}

	return 0, nil
}

// List - Generate a handler to handle a list query.
func List(modelIns interface{}, paramCreators ...CriteriaCreator) gin.HandlerFunc {
	var retFunc gin.HandlerFunc
	retFunc = func(c *gin.Context) {
		ptr := reflect.ValueOf(retFunc).Pointer()
		ptr2 := reflect.ValueOf(c.Handler()).Pointer()

		myType := reflect.TypeOf(modelIns)

		slice := reflect.MakeSlice(reflect.SliceOf(myType), 0, 0)
		x := reflect.New(slice.Type())
		x.Elem().Set(slice)

		if query, err := BuildQueryDB(modelIns, c); err == nil {
			for _, paramCreator := range paramCreators {
				query = paramCreator(query, c)
			}

			query.Find(x.Interface())

			var total int
			query.Model(modelIns).Count(&total)

			//total, data := ListHandlerWithoutServe(modelIns, c, paramCreators...)

			json := gin.H{
				"total": total,
				"data":  x.Interface(),
			}

			if ptr == ptr2 {
				c.JSON(200, json)
			} else {
				c.Set("list", json)
			}
		} else {
			ErrorResponse(c, 400, err.Error())
		}

	}

	return retFunc
}
