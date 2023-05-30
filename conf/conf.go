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
		UserName:  viper.GetString("account"),
		Pwd:       viper.GetString("password"),
		Jiraurl:   viper.GetString("jiraurl"),
		Project:   viper.GetString("project"),
		Label:     viper.GetString("label"),
		Component: viper.GetString("component"),
	}

}

func GetConfig() *BaseConfig {
	return baseConfig
}
