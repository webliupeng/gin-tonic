package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtConfig(t *testing.T) {

	//defer os.RemoveAll("/tmp/configs")

	config := GetConfig()
	assert.Equal(t, 0, config.GetInt("redis.db"))
}

func TestCommonFuncs(t *testing.T) {

	assert.Equal(t, "Foo", UpperInitial("foo"))
}

func TestConfigThreadSafe(t *testing.T) {

	for i := 0; i < 10000; i++ {
		go func() {
			c := GetConfig()

			fmt.Println("...", c.Get("db.host"))
		}()
	}
}

func init() {
	_ = os.Mkdir("/tmp/configs", 0777)

	text := `{
		"app": {
			"port": "8001"
		},
		"db": {
			"name": "test",
			"host":"localhost",
			"user": "xxx",
			"password": "C",
		},
		"redis": {
			"host": "localhost",
			"port": "6379",
			"password": "xxx"
		},
		"es": {
			"url": "http://localhost/"
		},
		"ext": {
			"foo": "bar",
			"num": 1
		}
	}`

	if err := ioutil.WriteFile("/tmp/configs/config.json", []byte(text), 0644); err != nil {
		panic(err)
	} else {
		fmt.Println("add config")
	}

}
