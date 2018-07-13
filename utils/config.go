package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/objx"
)

type Config struct {
	App struct {
		Port string `json:"port"`
	} `json:"app"`

	Db struct {
		Name     string `json:"name"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"db"`
	Redis struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Password string `json:"password"`
	} `json:"redis"`
	Open struct {
		URL string `json:"url"`
	} `json:"open"`
	ES struct {
		URL string `json:"url"`
	} `json:"es"`

	Ext interface{} `json:"ext"`
}

func (c *Config) GetExt(keyPath string) *objx.Value {
	ext := objx.New(c.Ext)
	return ext.Get(keyPath)
}

func GetConfig() Config {
	ret := Config{}

	env := os.Getenv("ENV")

	configDir := os.Getenv("CONFIGDIR")

	if configDir == "" {
		configDir = "./configs"
	}

	if env == "" {
		env = gin.DebugMode
	}

	configFile := path.Join(configDir, "config."+env+".json")
	if d, err := ioutil.ReadFile(configFile); err == nil {
		if err := json.Unmarshal(d, &ret); err != nil {
			//panic(err)
			fmt.Println("Config format incorrect")
		}
	} else {
		fmt.Println("no config file", err)
	}

	return ret
}
