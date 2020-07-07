package utils

import (
	"fmt"
	"path"
	"path/filepath"
	"sync"

	"strings"

	"log"

	"flag"

	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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

var globalConfig = &Config{}
var configInited = false

var (
	ConfigFile = flag.String("c", "", "config file")
)

var viperRuntime = viper.New()

var locker sync.Mutex

func GetConfig() *viper.Viper {
	if configInited {
		return viperRuntime
	}

	locker.Lock()
	defer locker.Unlock()
	
	testing.Init()
	
	flag.Parse()

	viperRuntime.SetConfigType("json")
	var err error
	if *ConfigFile != "" { // 参数指定了配置文件
		if file, err := os.Open(*ConfigFile); err == nil {
			log.Println("read specified config file", *ConfigFile)
			if err = viperRuntime.ReadConfig(file); err != nil {
				panic(err)
			}
		}

		viperRuntime.AddConfigPath(path.Dir(*ConfigFile))
		var basename = filepath.Base(*ConfigFile)
		var extension = filepath.Ext(basename)
		var configname = basename[0 : len(basename)-len(extension)]
		viperRuntime.SetConfigName(configname)
	} else {
		viperRuntime.SetConfigName("config")
		viperRuntime.AddConfigPath("./")
		err = viperRuntime.ReadInConfig() // Find and read the config file
	}

	viperRuntime.SetEnvPrefix("GTC")
	viperRuntime.AutomaticEnv()
	viperRuntime.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viperRuntime.Unmarshal(globalConfig)
	if err != nil { // Handle errors reading the config file
		fmt.Println(fmt.Errorf("Fatal error config file: %s", err))
	} else {
		if !configInited {
			configInited = true
			viperRuntime.WatchConfig()
			viperRuntime.OnConfigChange(func(in fsnotify.Event) {
				log.Println("config file change")
			})
		}
	}

	return viperRuntime
}
