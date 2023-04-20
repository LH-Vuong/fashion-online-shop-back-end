package initializers

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST" validate:"required"`
	DBUserName     string `mapstructure:"POSTGRES_USER" validate:"required"`
	DBUserPassword string `mapstructure:"POSTGRES_PASS" validate:"required"`
	DBName         string `mapstructure:"POSTGRES_DB" validate:"required"`
	DBPort         string `mapstructure:"POSTGRES_PORT" validate:"required"`
	ServerPort     string `mapstructure:"PORT" validate:"required"`

	MongoUrl string `mapstructure:"MONGO_URL" validate:"required"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
