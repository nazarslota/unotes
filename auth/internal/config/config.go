package config

import (
	"os"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Auth struct {
		HostHTTP              string        `mapstructure:"AUTH_HOST_HTTP"`
		PortHTTP              string        `mapstructure:"AUTH_PORT_HTTP"`
		HostGRPC              string        `mapstructure:"AUTH_HOST_GRPC"`
		PortGRPC              string        `mapstructure:"AUTH_PORT_GRPC"`
		AccessTokenSecret     string        `mapstructure:"AUTH_ACCESS_TOKEN_SECRET"`
		AccessTokenExpiresIn  time.Duration `mapstructure:"AUTH_ACCESS_TOKEN_EXPIRES_IN"`
		RefreshTokenSecret    string        `mapstructure:"AUTH_REFRESH_TOKEN_SECRET"`
		RefreshTokenExpiresIn time.Duration `mapstructure:"AUTH_REFRESH_TOKEN_EXPIRES_IN"`
	} `mapstructure:",squash"`
	PostgreSQL struct {
		Host     string `mapstructure:"AUTH_POSTGRESQL_HOST"`
		Port     string `mapstructure:"AUTH_POSTGRESQL_PORT"`
		Username string `mapstructure:"AUTH_POSTGRESQL_USERNAME"`
		Password string `mapstructure:"AUTH_POSTGRESQL_PASSWORD"`
		DBName   string `mapstructure:"AUTH_POSTGRESQL_DBNAME"`
		SSLMode  string `mapstructure:"AUTH_POSTGRESQL_SSLMODE"`
	} `mapstructure:",squash"`
	MongoDB struct {
		Host     string `mapstructure:"AUTH_MONGODB_HOST"`
		Port     string `mapstructure:"AUTH_MONGODB_PORT"`
		Username string `mapstructure:"AUTH_MONGODB_USERNAME"`
		Password string `mapstructure:"AUTH_MONGODB_PASSWORD"`
		Database string `mapstructure:"AUTH_MONGODB_DATABASE"`
	} `mapstructure:",squash"`
	Redis struct {
		Addr     string `mapstructure:"AUTH_REDIS_ADDR"`
		Password string `mapstructure:"AUTH_REDIS_PASSWORD"`
		DB       int    `mapstructure:"AUTH_REDIS_DB"`
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

			bindEnvAuth(v)
			if err := v.ReadInConfig(); err != nil {
				panic(err)
			}
		} else if environment == "STAGE" {
			v.AutomaticEnv()
			v.SetConfigName("stage")

			bindEnvAuth(v)
			if err := v.ReadInConfig(); err != nil {
				panic(err)
			}
		} else {
			v.AutomaticEnv()
			v.SetConfigName("production")

			bindEnvAuth(v)
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

func bindEnvAuth(v *viper.Viper) {
	_ = v.BindEnv("AUTH_ACCESS_TOKEN_SECRET")
	_ = v.BindEnv("AUTH_REFRESH_TOKEN_SECRET")
	bindEnvPostgreSQL(v)
	bindEnvMongoDB(v)
	bindEnvRedis(v)
}

func bindEnvPostgreSQL(v *viper.Viper) {
	_ = v.BindEnv("AUTH_POSTGRESQL_HOST")
	_ = v.BindEnv("AUTH_POSTGRESQL_PORT")
	_ = v.BindEnv("AUTH_POSTGRESQL_USERNAME")
	_ = v.BindEnv("AUTH_POSTGRESQL_PASSWORD")
	_ = v.BindEnv("AUTH_POSTGRESQL_DBNAME")
	_ = v.BindEnv("AUTH_POSTGRESQL_SSLMODE")
}

func bindEnvMongoDB(v *viper.Viper) {
	_ = v.BindEnv("AUTH_MONGODB_HOST")
	_ = v.BindEnv("AUTH_MONGODB_PORT")
	_ = v.BindEnv("AUTH_MONGODB_USERNAME")
	_ = v.BindEnv("AUTH_MONGODB_PASSWORD")
	_ = v.BindEnv("AUTH_MONGODB_DATABASE")
}

func bindEnvRedis(v *viper.Viper) {
	_ = v.BindEnv("AUTH_REDIS_ADDR")
	_ = v.BindEnv("AUTH_REDIS_PASSWORD")
	_ = v.BindEnv("AUTH_REDIS_DB")
}
