package utils

import (
	"fmt"

	"strings"

	"log"

	"flag"

	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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
		Db       int    `json:"db"`
	} `json:"redis"`
	Open struct {
		URL string `json:"url"`
	} `json:"open"`
	ES struct {
		URL string `json:"url"`
	} `json:"es"`
}

func (c *Config) GetExt(keyPath string) *objx.Value {
	//ext := objx.New(c.Ext)
	//return ext.Get(keyPath)
	log.Println("GetExt is deprecated. Use viper.Get to instead")
	return nil
}

var globalConfig = &Config{}
var configInited = false

var (
	configFile = flag.String("c", "", "config file")
)

var viperRuntime = viper.New()

func GetConfig() *viper.Viper {
	if configInited {
		return viperRuntime
	}

	flag.Parse()

	viperRuntime.SetConfigName("config")
	viperRuntime.SetConfigType("json")
	viperRuntime.AddConfigPath("./")

	var err error
	if *configFile != "" {
		if file, err := os.Open(*configFile); err == nil {
			if err = viperRuntime.ReadConfig(file); err != nil {
				panic(err)
			}
		}
	} else {
		err = viperRuntime.ReadInConfig() // Find and read the config file
	}

	viperRuntime.SetEnvPrefix("GTC")
	viperRuntime.AutomaticEnv()
	viperRuntime.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viperRuntime.Unmarshal(globalConfig)
	if err != nil { // Handle errors reading the config file
		fmt.Println(fmt.Errorf("Fatal error config file: %s", err))
		fmt.Println("db host", viper.Get("db.host"), globalConfig.Db.Host)
	} else {
		if !configInited {

			configInited = true

			viperRuntime.WatchConfig()
			viperRuntime.OnConfigChange(func(in fsnotify.Event) {
				log.Println("config file change")
				//viper.Unmarshal(globalConfig)
			})
		}
	}

	return viperRuntime
}
