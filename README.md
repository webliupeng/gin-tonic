[![Build Status](https://travis-ci.org/webliupeng/gin-tonic.svg?branch=master)](https://travis-ci.org/webliupeng/gin-tonic) [![Coverage Status](https://coveralls.io/repos/github/webliupeng/gin-tonic/badge.svg?branch=master)](https://coveralls.io/github/webliupeng/gin-tonic?branch=master)

Gin-tonic is inspired by redstone's open-rest project which helps developers building CRUD APIs with Gin and GORM fastly. 

**Features**
- Genrate list filterable handler
- Genrate delete handler
- Genrate detail handler 
- Genrate create handler with Gin's validation
- Genrate update handler with Gin's validation


**USAGE**

    import (
		"github.com/webliupeng/gin-tonic/helpers"
	)

	type Foo struct {
		Bar string
	}
	func (f *Foo)InsertableFields() []string {
		return []string{"bar"}
	}
	
	func (f *Foo)UpdateableFields() []string {
		return []string{"bar"}
	}

	router.POST("/posts", helpers.Create(func(c *gin.Context){ return &Foo{} }))
	router.GET("/posts", helpers.List(&Foo{}))
	router.GET("/posts/:id", helpers.FindOneByParam(&Foo{}, "id", "foo"), helpers.ServeJSONFromContext("foo"))
	router.DELETE("/posts/:id", helpers.FindOneByParam(&Foo{}, "id", "foo"), helpers.Delete("foo"))
	router.PUT("/posts/:id", helpers.FindOneByParam(&Foo{}, "id", "foo"), helpers.Update("foo"))
