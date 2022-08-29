package config

import (
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Host                  string        `mapstructure:"host" validate:"required"`
	Port                  string        `mapstructure:"port" validate:"required"`
	LoggerOutputFile      string        `mapstructure:"logger_output" validate:"required"`
	AccessTokenSecret     string        `mapstructure:"access_token_secret" validate:"required"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"access_token_expires_in" validate:"required"`
	RefreshTokenSecret    string        `mapstructure:"refresh_token_secret" validate:"required"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"refresh_token_expires_in" validate:"required"`
	Debug                 bool          `mapstructure:"debug"`
	MongoDB               struct {
		Host            string `mapstructure:"host" validate:"required"`
		Port            string `mapstructure:"port"`
		Username        string `mapstructure:"username"`
		Password        string `mapstructure:"password"`
		Database        string `mapstructure:"database" validate:"required"`
		UsersCollection string `mapstructure:"users_collection" validate:"required"`
	} `mapstructure:"mongodb" validate:"required"`
	Redis struct {
		Addr     string `mapstructure:"addr" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
		Database int    `mapstructure:"database"`
	} `mapstructure:"redis" validate:"required"`
}

const File = "./configs/config.json"

var (
	once     sync.Once
	instance Config
)

func C() *Config {
	once.Do(func() {
		filename := path.Base(File)
		filenameWithoutExt := filename[:len(filename)-len(filepath.Ext(filename))]

		v := viper.New()
		v.AddConfigPath(path.Dir(File))
		v.SetConfigName(filenameWithoutExt)

		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := v.Unmarshal(&instance); err != nil {
			panic(err)
		}

		if err := validator.New().Struct(&instance); err != nil {
			panic(err)
		}
	})
	return &instance
}
