package helpers

import (
	"fmt"
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

// List - 列表输出模型
func List(modelIns interface{}, paramCreators ...CriteriaCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
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
				return key + " IN (?)", val
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
			for _, table := range strings.Split(includes, ",") {
				query = query.Preload(utils.UpperInitial(table))
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

		var total int

		//fmt.Println("recevied max result", pageSize)
		query.Limit(pageSize).Offset(offset).Find(x.Interface())

		db.DB().Model(x.Interface()).Where(strings.Join(expressions, " AND "), values...).Count(&total)

		c.JSON(200, gin.H{
			"total": total,
			"data":  x.Interface(),
		})
	}
}
