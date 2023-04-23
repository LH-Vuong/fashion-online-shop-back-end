package initializers

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST" validate:"required"`
	DBUserName     string `mapstructure:"POSTGRES_USER" validate:"required"`
	DBUserPassword string `mapstructure:"POSTGRES_PASS" validate:"required"`
	DBName         string `mapstructure:"POSTGRES_DB" validate:"required"`
	DBPort         string `mapstructure:"POSTGRES_PORT" validate:"required"`
	ServerPort     string `mapstructure:"PORT" validate:"required"`
	ClientUrl      string `mapstructure:"CLIENT_URL" validate:"required"`
	MongoUrl       string `mapstructure:"MONGO_URL" validate:"required"`

	EmailFrom string `mapstructure:"EMAIL_FROM" validate:"required"`
	SMTPHost  string `mapstructure:"SMTP_HOST" validate:"required"`
	SMTPPort  int    `mapstructure:"SMTP_PORT" validate:"required"`
	SMTPUser  string `mapstructure:"SMTP_USER" validate:"required"`
	SMTPPass  string `mapstructure:"SMTP_PASS" validate:"required"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
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
