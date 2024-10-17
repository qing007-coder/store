package config

import (
	"fmt"
	"github.com/spf13/viper"
	"store/pkg/errors"
)

type GlobalConfig struct {
	Mysql struct {
		Address  string `yaml:"address"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
		Name     string `yaml:"name"`
		Conf     string `yaml:"conf"`
		Password string `yaml:"password"`
	} `yaml:"mysql"`

	Redis struct {
		Address  string `yaml:"address"`
		Port     string `yaml:"port"`
		DB       int    `yaml:"db"`
		Password string `yaml:"password"`
	} `yaml:"redis"`

	JWT struct {
		SecretKey string `yaml:"secretKey"`
		Expiry    int    `yaml:"expiry"`
	} `yaml:"jwt"`

	Elasticsearch struct {
		Address string `yaml:"address"`
	} `yaml:"elasticsearch"`

	SecretKey string `yaml:"secretKey"`
	Logger    struct {
		MaxSize    int `yaml:"max_size"`
		MaxBackups int `yaml:"max_backups"`
		MaxAge     int `yaml:"max_age"`
	} `yaml:"logger"`

	Email struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
	} `yaml:"email"`
}

func NewGlobalConfig() (*GlobalConfig, error) {
	conf := new(GlobalConfig)
	if err := conf.init(); err != nil {
		return nil, err
	}
	return conf, nil
}

func (c *GlobalConfig) init() error {
	v := viper.New()
	v.AddConfigPath("./config/")
	v.SetConfigName("common")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			fmt.Println(1111)
			return errors.ConfigFileNotFound
		} else {
			return errors.OtherError
		}
	}

	if err := v.Unmarshal(c); err != nil {
		return errors.UnmarshalError
	}

	return nil
}
