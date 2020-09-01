[![Build Status](https://travis-ci.org/webliupeng/gin-tonic.svg?branch=master)](https://travis-ci.org/webliupeng/gin-tonic) [![Coverage Status](https://coveralls.io/repos/github/webliupeng/gin-tonic/badge.svg?branch=master)](https://coveralls.io/github/webliupeng/gin-tonic?branch=master)

中文 | [English](https://github.com/webliupeng/gin-tonic/blob/master/README_EN.md)

Gin-tonic 能够帮助研发人员快速构建Restful API，它可以将GORM的model暴露为CRUD handler.


**快速开始**

需要为model实现API访问相关的接口
```go
/* 定义模型 */
type MyModel struct {
  Name string `json:name`
  Age uint `json:age`
  Address string `json:address`
}

// 通过InsertableFields方法返回值定义可以插入的字段,例如返回 'age','name' 和 'address' 字段可以插入
func (f *Foo)InsertableFields() []string {
	return []string{"age", "name", "address"}
}

// 通过UpdateableFields方法返回值定义可以更新的字段,例如返回 'age' 和 'address' 字段可更新
func (f *Foo)UpdateableFields() []string {
	return []string{"age", "address"}
}

// 通过SortableFields方法返回值定义可以用于list接口参数指定排序的字段,例如返回 'age' 字段可用于排序
func (f *Foo)SortableFields() []string {
	return []string{"age"}
}

// 通过FilterableFields方法返回值定义可以用于list接口参数过滤字段白名单,例如返回 'age' 字段可用于筛选
func (f *Foo)FilterableFields() []string {
	return []string{"age", "name", "address"}
}
```

*列表接口示例*
```go
import (
	"github.com/webliupeng/gin-tonic/helpers"
)

router.Get("/models", helpers.List(&MyModel{}))


// 也可以自己定义gorm的查询条件
router.Get("/models2", helpers.List(&MyModel{}, func(db *gorm.DB, c *gin.Context) {

  return db.Where("file", "=", 321)

}))
```

 gin-tonic 生成的list handler默认支持排序，过滤，分页。
 
 querystring参数表达式：

| 表达式 | 用法                    |
| ------- | ------------------------ |
| >       | field_gt=val             |
| >=      | field_gte=val            |   
| <       | field_lt=val             |
| <=      | field_lte=val            |
| like    | field_like=val           |
| in      | field_in=val1,val2,valn  |
| not in  | field_not=val1,val2,valn |


分页参数：

| 表达式 | 用法                    |
| ------- | ------------------------ |
| .maxResults       |       最大返回记录数       |
| .offset     | 记录起始位置            |
```shell
curl http://hostname/models?.maxResults=100&.offset=10 # equals LIMIT 10, 100
curl http://hostname/models?age_lt=10  #list custerms filtered by age less than 10
curl http://hostname/models?age=10
```
*创建记录处理器示例*

创建处理器会通过gin的ShouldBindBodyWith解析请求body，如果model定义了[validator](https://github.com/go-playground/validator)的规则，将会校验请求的数据，通过才会插入数据库。

```go
router.Post("/models", helpers.Create(func(c *gin.Context) interface{} {
  mol := &Models{}

  // 可以在这里设定默认值

  return mol
}))
```

```shell
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name": "a"}' \
  http://yourdomain/models

```
*更新处理器示例*
更新处理器会通过gin的ShouldBindBodyWith解析请求body，如果model定义了[validator](https://github.com/go-playground/validator)的规则，将会校验请求的数据，通过才会更新入数据库。
```go
router.Put("/models/:id", 
  helpers.FindOne(&MyModel{}, "id", "mymodel"), // 通过 id 参数查询一条MyModel的记录，并将结果暂存在context
  helpers.Update("customer") // 更新mymodel
)
```

*删除处理器示例* 
```go
router.Delete("/models/:id", 
  helpers.FindOne(&Customer{}, "id", "customer"), // 通过 id 参数查询一条MyModel的记录，并将结果暂存在context
  helpers.Delete("customer") // delete the record by context name 
)

```

*配置*
- Gin-tonic是通过[Viper](https://github.com/spf13/viper) 来读取配置的
`gin-tonic`，推荐使用以下环境变量来配置数据库。

```shell
export GTC_DB_HOST=localhost
export GTC_DB_NAME=dbname
export GTC_DB_PORT=3306
export GTC_DB_USER=username
export GTC_DB_PASSWORD=password
```
