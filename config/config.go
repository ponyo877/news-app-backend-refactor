package config

import (
	"github.com/labstack/gommon/log"

	"github.com/spf13/viper"
)

type MysqlConfig struct {
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBDatabase string `mapstructure:"DB_DATABASE"`
}

type RedisConfig struct {
	KVSPassword string `mapstructure:"KVS_PASSWORD"`
	KVSHost     string `mapstructure:"KVS_HOST"`
	KVSDatabase int    `mapstructure:"KVS_DATABASE"`
}

func LoadMysqlConfig() (MysqlConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_DATABASE")
	var config MysqlConfig
	if err := viper.Unmarshal(&config); err != nil {
		return MysqlConfig{}, err
	}
	log.Infof("user: %v, pass: %v, host: %v, db: %v", config.DBUser, config.DBPassword, config.DBHost, config.DBDatabase)
	return config, nil
}

func LoadRedisConfig() (RedisConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("KVS_PASSWORD")
	viper.BindEnv("KVS_HOST")
	viper.BindEnv("KVS_DATABASE")
	var config RedisConfig
	if err := viper.Unmarshal(&config); err != nil {
		return RedisConfig{}, err
	}
	log.Infof("pass: %v, host: %v, db: %v", config.KVSPassword, config.KVSHost, config.KVSDatabase)
	return config, nil

}
