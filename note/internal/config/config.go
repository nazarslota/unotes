package config

import (
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Note struct {
		HostGRPC          string `mapstructure:"NOTE_HOST_GRPC"`
		PortGRPC          string `mapstructure:"NOTE_PORT_GRPC"`
		HostREST          string `mapstructure:"NOTE_HOST_REST"`
		PortREST          string `mapstructure:"NOTE_PORT_REST"`
		Debug             bool   `mapstructure:"NOTE_DEBUG"`
		Log               string `mapstructure:"NOTE_LOG"`
		AccessTokenSecret string `mapstructure:"NOTE_ACCESS_TOKEN_SECRET"`
	} `mapstructure:",squash"`
	MongoDB struct {
		Host     string `mapstructure:"NOTE_MONGODB_HOST"`
		Port     string `mapstructure:"NOTE_MONGODB_PORT"`
		Username string `mapstructure:"NOTE_MONGODB_USERNAME"`
		Password string `mapstructure:"NOTE_MONGODB_PASSWORD"`
		Database string `mapstructure:"NOTE_MONGODB_DATABASE"`
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

		environment := os.Getenv("NOTE_ENVIRONMENT")
		if environment == "DEVELOPMENT" {
			v.AutomaticEnv()
			v.SetConfigName("development")

			bindEnv(v)
			if err := v.ReadInConfig(); err != nil {
				panic(err)
			}
		} else if environment == "STAGE" {
			v.AutomaticEnv()
			v.SetConfigName("stage")

			bindEnv(v)
			if err := v.ReadInConfig(); err != nil {
				panic(err)
			}
		} else {
			v.AutomaticEnv()
			v.SetConfigName("production")

			bindEnv(v)
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

func bindEnv(v *viper.Viper) {
	_ = v.BindEnv("NOTE_ACCESS_TOKEN_SECRET")
	bindEnvMongoDB(v)
}

func bindEnvMongoDB(v *viper.Viper) {
	_ = v.BindEnv("NOTE_MONGODB_HOST")
	_ = v.BindEnv("NOTE_MONGODB_PORT")
	_ = v.BindEnv("NOTE_MONGODB_USERNAME")
	_ = v.BindEnv("NOTE_MONGODB_PASSWORD")
	_ = v.BindEnv("NOTE_MONGODB_DATABASE")
}
