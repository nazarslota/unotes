package config

import (
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Auth struct {
		Host                  string        `mapstructure:"AUTH_HOST"`
		Port                  string        `mapstructure:"AUTH_PORT"`
		AccessTokenSecret     string        `mapstructure:"AUTH_ACCESS_TOKEN_SECRET"`
		AccessTokenExpiresIn  time.Duration `mapstructure:"AUTH_ACCESS_TOKEN_EXPIRES_IN"`
		RefreshTokenSecret    string        `mapstructure:"AUTH_REFRESH_TOKEN_SECRET"`
		RefreshTokenExpiresIn time.Duration `mapstructure:"AUTH_REFRESH_TOKEN_EXPIRES_IN"`
	} `mapstructure:",squash"`
	PostgreSQL struct {
		Host     string `mapstructure:"POSTGRESQL_HOST"`
		Port     string `mapstructure:"POSTGRESQL_PORT"`
		Username string `mapstructure:"POSTGRESQL_USERNAME"`
		Password string `mapstructure:"POSTGRESQL_PASSWORD"`
		DBName   string `mapstructure:"POSTGRESQL_DBNAME"`
		SSLMode  string `mapstructure:"POSTGRESQL_SSLMODE"`
	} `mapstructure:",squash"`
	MongoDB struct {
		Host     string `mapstructure:"MONGODB_HOST"`
		Port     string `mapstructure:"MONGODB_PORT"`
		Username string `mapstructure:"MONGODB_USERNAME"`
		Password string `mapstructure:"MONGODB_PASSWORD"`
		Database string `mapstructure:"MONGODB_DATABASE"`
	} `mapstructure:",squash"`
	Redis struct {
		Addr     string `mapstructure:"REDIS_ADDR"`
		Password string `mapstructure:"REDIS_PASSWORD"`
		DB       int    `mapstructure:"REDIS_DB"`
	} `mapstructure:",squash"`
}

var (
	once     sync.Once
	instance *Config
)

func C() *Config {
	once.Do(func() {
		v := viper.New()
		v.SetConfigName("app")
		v.AddConfigPath("configs")

		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := v.Unmarshal(&instance); err != nil {
			panic(err)
		}

		if err := validator.New().Struct(instance); err != nil {
			panic(err)
		}
	})
	return instance
}
