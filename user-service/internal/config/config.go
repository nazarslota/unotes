package config

import (
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	UserService struct {
		Host string `mapstructure:"USER_SERVICE_HOST"`
		Port string `mapstructure:"USER_SERVICE_PORT"`
	} `mapstructure:",squash"`
	PostgreSQL struct {
		Host     string `mapstructure:"USER_SERVICE_POSTGRESQL_HOST"`
		Port     string `mapstructure:"USER_SERVICE_POSTGRESQL_PORT"`
		Username string `mapstructure:"USER_SERVICE_POSTGRESQL_USERNAME"`
		Password string `mapstructure:"USER_SERVICE_POSTGRESQL_PASSWORD"`
		DBName   string `mapstructure:"USER_SERVICE_POSTGRESQL_DBNAME"`
		SSLMode  string `mapstructure:"USER_SERVICE_POSTGRESQL_SSLMODE"`
	} `mapstructure:",squash"`
}

var (
	once     sync.Once
	instance *Config
)

func C() *Config {
	once.Do(func() {
		v := viper.New()
		v.AddConfigPath("configs")

		environment := os.Getenv("ENVIRONMENT")
		if environment == "DEVELOPMENT" {
			v.AutomaticEnv()
			v.SetConfigName("development")

			bindEnvAuthService(v)
			if err := v.ReadInConfig(); err != nil {
				panic(err)
			}
		} else if environment == "STAGE" {
			v.AutomaticEnv()
			v.SetConfigName("stage")

			bindEnvAuthService(v)
			if err := v.ReadInConfig(); err != nil {
				panic(err)
			}
		} else {
			v.AutomaticEnv()
			v.SetConfigName("production")

			bindEnvAuthService(v)
			if err := v.ReadInConfig(); err != nil {
				panic(err)
			}
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

func bindEnvAuthService(v *viper.Viper) {
	bindEnvPostgreSQL(v)
}

func bindEnvPostgreSQL(v *viper.Viper) {
	_ = v.BindEnv("USER_SERVICE_POSTGRESQL_HOST")
	_ = v.BindEnv("USER_SERVICE_POSTGRESQL_PORT")
	_ = v.BindEnv("USER_SERVICE_POSTGRESQL_USERNAME")
	_ = v.BindEnv("USER_SERVICE_POSTGRESQL_PASSWORD")
	_ = v.BindEnv("USER_SERVICE_POSTGRESQL_DBNAME")
	_ = v.BindEnv("USER_SERVICE_POSTGRESQL_SSLMODE")
}
