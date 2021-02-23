package core

import (
	"fmt"
	"main/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Viper 加载配置文件
func Viper() {
	viper.SetConfigFile(global.ConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&global.CONFIG); err != nil {
		fmt.Println(err)
	}
	fmt.Println(global.CONFIG)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&global.CONFIG); err != nil {
			fmt.Println(err)
		}
		fmt.Println(global.CONFIG)
	})
}
