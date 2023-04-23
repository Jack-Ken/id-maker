package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"id-maker/config"
)

// 加载配置文件
func Load_Config() (err error) {
	// 读取配置文件
	viper.SetConfigName("config")                        // 配置文件名称（无扩展名）
	viper.SetConfigType("yaml")                          // 如果配置文件中没有扩展名，则需要配置此项
	viper.AddConfigPath("J:/go-project/id-maker/config") // 查找配置文件所在路径

	if err = viper.ReadInConfig(); err != nil { // 处理读取配置文件的错误
		fmt.Println("Fatal error config file: %v \n", err)
		return
	}

	if err = viper.Unmarshal(config.Conf); err != nil {
		fmt.Printf("viper.Unmarshal config file failed: %v \n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("The config file has cahnged...")
		if err = viper.Unmarshal(config.Conf); err != nil {
			fmt.Printf("viper.Unmarshal config file failed: %v \n", err)
		}
	})
	return
}
