package configs

import (
	"github.com/spf13/viper"
)

func init() {

	viper.SetConfigName("env")
	viper.SetConfigType("json")
	viper.AddConfigPath("./test_configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}
