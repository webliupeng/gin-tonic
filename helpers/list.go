package helpers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/webliupeng/gin-tonic/db"
	"github.com/webliupeng/gin-tonic/utils"
)

type condition struct {
	Name     string
	Oparator string
}

// CriteriaCreator - 条件构造
type CriteriaCreator func(c *gin.Context) (string, string)

func CriteriaByParam(key string) CriteriaCreator {
	f := func(c *gin.Context) (string, string) {
		if val, ok := c.Params.Get(key); ok {
			return key, val
		} else {
			return key, ""
		}
	}
	return f
}

func ListHandlerWithoutServe(modelIns interface{}, c *gin.Context, paramCreators ...CriteriaCreator) (int, interface{}) {
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
		}
		for opr, handle := range oparators {
			if val := c.Query(fmt.Sprintf("%v_%v", field, opr)); val != "" {
				expression, qv := handle(field, val)
				expressions = append(expressions, expression)
				values = append(values, qv)
				break
			}
		}
	}

	for _, paramCreator := range paramCreators {
		key, value := paramCreator(c)
		expressions = append(expressions, fmt.Sprintf("%v = ?", key))
		values = append(values, value)
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
					query = query.Preload(utils.UpperInitial(table))
				} else {
					ErrorResponse(c, http.StatusBadRequest, "Can not includes "+table)
					return 0, nil
				}
			}
		} else {
			ErrorResponse(c, http.StatusBadRequest, "Can not includes a non-includable model")
			return 0, nil
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

	var total int

	query.Limit(pageSize).Offset(offset).Find(x.Interface())
	db.DB().Model(x.Interface()).Where(strings.Join(expressions, " AND "), values...).Count(&total)
	return total, x.Interface()
}

// List - Generate a handler to handle a list query , if the li
func List(modelIns interface{}, paramCreators ...CriteriaCreator) gin.HandlerFunc {
	var retFunc gin.HandlerFunc
	retFunc = func(c *gin.Context) {
		ptr := reflect.ValueOf(retFunc).Pointer()
		ptr2 := reflect.ValueOf(c.Handler()).Pointer()
		total, data := ListHandlerWithoutServe(modelIns, c, paramCreators...)

		json := gin.H{
			"total": total,
			"data":  data,
		}
		if ptr == ptr2 {
			c.JSON(200, json)
		} else {
			c.Set("list", json)
		}
	}

	return retFunc
}
