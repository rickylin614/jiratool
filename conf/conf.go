package conf

import (
	"github.com/spf13/viper"
)

var baseConfig *BaseConfig = &BaseConfig{}

// 讀取設定檔
func ConfigInit() {
	viper.SetConfigFile("conf.yml")
	viper.ReadInConfig() // ignore error

	baseConfig = &BaseConfig{
		AdminPath:     viper.GetString("admin_path"),
		AdminFileName: viper.GetString("admin_filename"),
		WebPath:       viper.GetString("source_path"),
		WebFileName:   viper.GetString("source_filename"),
	}

}

func GetConfig() *BaseConfig {
	return baseConfig
}
