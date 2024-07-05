package bootstrap

import (
	"fmt"
	"github.com/NumberMan1/MMO-server/game_server/config"
	"github.com/NumberMan1/MMO-server/game_server/database"
	"github.com/NumberMan1/common/global"
	"github.com/NumberMan1/common/summer/proto_helper"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

func initImpl() {
	//解析yaml
	database.Init()
	config.ServerInit()
	//初始化pb字典
	proto_helper.Init()
}

func Init() {
	// Do your application logic here
	fmt.Println("Running game server with config file:", cfgFile)
	// Load the configuration file
	viper.SetConfigFile(cfgFile)
	// Watch for config file changes and re-read it when necessary
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}
		// Example of reading a config value
		port := viper.GetInt("server.port")
		fmt.Println("Server port:", port)
	})
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// TODO:也通过viper解析
	global.Init(cfgFile)
	initImpl()
}
