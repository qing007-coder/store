package config

import (
	"github.com/spf13/viper"
	"store/pkg/errors"
)

type GlobalConfig struct {
	Mysql struct {
		Addr     string `yaml:"addr"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
		Name     string `yaml:"name"`
		Conf     string `yaml:"conf"`
		Password string `yaml:"password"`
	} `yaml:"mysql"`

	Redis struct {
		Addr     string `yaml:"addr"`
		Port     string `yaml:"port"`
		DB       int    `yaml:"db"`
		Password string `yaml:"password"`
	} `yaml:"redis"`

	JWT struct {
		SecretKey     string `yaml:"secretKey"`
		AccessExpiry  int    `yaml:"accessExpiry"`
		RefreshExpiry int    `yaml:"refreshExpiry"`
	} `yaml:"jwt"`

	Elasticsearch struct {
		Addr string `yaml:"addr"`
		Port string `yaml:"port"`
	} `yaml:"elasticsearch"`

	SecretKey string `yaml:"secretKey"`
	Logger    struct {
		MaxSize    int `yaml:"maxSize"`
		MaxBackups int `yaml:"maxBackups"`
		MaxAge     int `yaml:"maxAge"`
	} `yaml:"logger"`

	Email struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
	} `yaml:"email"`

	Consul struct {
		Addr string `yaml:"addr"`
		Port string `yaml:"port"`
		Name string `yaml:"name"`
	} `yaml:"consul"`

	Minio struct {
		Endpoint  string `yaml:"endpoint"`
		Port      string `yaml:"port"`
		AccessKey string `yaml:"accessKey"`
		SecretKey string `yaml:"secretKey"`
	} `yaml:"minio"`

	Kafka struct {
		Addr string `yaml:"addr"`
		Port string `yaml:"port"`
	} `yaml:"kafka"`
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
