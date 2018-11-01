[![Build Status](https://travis-ci.org/webliupeng/gin-tonic.svg?branch=master)](https://travis-ci.org/webliupeng/gin-tonic) [![Coverage Status](https://coveralls.io/repos/github/webliupeng/gin-tonic/badge.svg?branch=master)](https://coveralls.io/github/webliupeng/gin-tonic?branch=master)

Gin-tonic is inspired by redstone's [Open-Rest](https://github.com/open-node/open-rest) project which helps developers building CRUD APIs with Gin and [GORM](https://github.com/jinzhu/gorm) fastly. 

**Features**
- Generate list handler
- Generate delete handler
- Generate detail handler 
- Generate create handler with Gin's validation
- Generate update handler with Gin's validation
- Integrate [Viper](https://github.com/spf13/viper) to read configuration

**Quick start**

*Define model*

*implement accessible interfaces to expose model CRUD behaviors
```go
type Customer struct {
  Name string
  Age uint
  Address string
}

// define the model only allows 'age','name' and 'name' fields can be insert
func (f *Foo)InsertableFields() []string {
	return []string{"age", "name", "address"}
}

// define the model only allows 'age' and 'address' fields can be update.
func (f *Foo)UpdateableFields() []string {
	return []string{"age", "address"}
}

// define the model only allows 'age' to sort
func (f *Foo)SortableFields() []string {
	return []string{"age"}
}

// define the model only allows 'age','name','address' to filter
func (f *Foo)FilterableFields() []string {
	return []string{"age", "name", "address"}
}
```


*List Handler Example*
```go
import (
	"github.com/webliupeng/gin-tonic/helpers"
)

router.Get("/customers", helpers.List(&Customer{}))
```
```shell
curl yourdomain/customers?.maxResults=100&.offset=10 # equals LIMIT 10, 100
curl yourdomain/customers?age_lt=10
```
*Create Handler Example*

```go
router.Post("/customer", helpers.Create(func(c *gin.Context) interface{} {
  customer := &Customers{}

  // set default values here.
  return customer
}))
```

```shell
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name": "a"}' \
  http://yourdomain/customers

```
*Update Handler Example*

```go
router.Put("/customers/:id", 
  helpers.FindOne(&Customer{}, "id", "customer"), // find a record by params 'id' and store the result to gin's Context
  helpers.Update("customer") // update the record by context name 
)
```

*Delete Handler Example* 
```go
router.Put("/customers/:id", 
  helpers.FindOne(&Customer{}, "id", "customer"), // find a record by params 'id' and store the result to gin's Context
  helpers.Delete("customer") // delete the record by context name 
)

```