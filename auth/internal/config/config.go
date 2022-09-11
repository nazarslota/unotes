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
	Auth struct {
		Host                  string        `mapstructure:"host"`
		Port                  string        `mapstructure:"port"`
		AccessTokenSecret     string        `mapstructure:"access_token_secret"`
		RefreshTokenSecret    string        `mapstructure:"refresh_token_secret"`
		AccessTokenExpiresIn  time.Duration `mapstructure:"access_token_expires_in"`
		RefreshTokenExpiresIn time.Duration `mapstructure:"refresh_token_expires_in"`
	} `mapstructure:"auth"`
}

const File = "configs/config.json"

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
