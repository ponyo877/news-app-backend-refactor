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
	DBPort     string `mapstructure:"DB_PORT"`
}

type RedisConfig struct {
	KVSPassword string `mapstructure:"KVS_PASSWORD"`
	KVSHost     string `mapstructure:"KVS_HOST"`
	KVSDatabase int    `mapstructure:"KVS_DATABASE"`
	KVSPort     string `mapstructure:"KVS_PORT"`
}

type ElasticSearchConfig struct {
	SEUser     string `mapstructure:"SE_USER"`
	SEPassword string `mapstructure:"SE_PASSWORD"`
	SESHost    string `mapstructure:"SE_HOST"`
	SEPort     string `mapstructure:"SE_PORT"`
}

type WebDAVConfig struct {
	WDUser     string `mapstructure:"WD_USER"`
	WDPassword string `mapstructure:"WD_PASSWORD"`
	WDSHost    string `mapstructure:"WD_HOST"`
	WDPort     string `mapstructure:"WD_PORT"`
}

type AppConfig struct {
	APRoot string `mapstructure:"AP_ROOT"`
	APPort string `mapstructure:"AP_PORT"`
}

// LoadMysqlConfig
func LoadMysqlConfig() (MysqlConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_DATABASE")
	viper.BindEnv("DB_PORT")
	var config MysqlConfig
	if err := viper.Unmarshal(&config); err != nil {
		return MysqlConfig{}, err
	}
	log.Infof("[MySQL] user: %v, pass: %v, host: %v, db: %v, port: %v", config.DBUser, config.DBPassword, config.DBHost, config.DBDatabase, config.DBPort)
	return config, nil
}

// LoadRedisConfig
func LoadRedisConfig() (RedisConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("KVS_PASSWORD")
	viper.BindEnv("KVS_HOST")
	viper.BindEnv("KVS_DATABASE")
	viper.BindEnv("KVS_PORT")
	var config RedisConfig
	if err := viper.Unmarshal(&config); err != nil {
		return RedisConfig{}, err
	}
	log.Infof("[Redis] pass: %v, host: %v, db: %v, port: %v", config.KVSPassword, config.KVSHost, config.KVSDatabase, config.KVSPort)
	return config, nil

}

// LoadElasticSearchConfig
func LoadElasticSearchConfig() (ElasticSearchConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("SE_USER")
	viper.BindEnv("SE_PASSWORD")
	viper.BindEnv("SE_HOST")
	viper.BindEnv("SE_PORT")
	var config ElasticSearchConfig
	if err := viper.Unmarshal(&config); err != nil {
		return ElasticSearchConfig{}, err
	}
	log.Infof("[ElasticSearch] user: %v, pass: %v, host: %v, port: %v", config.SEUser, config.SEPassword, config.SESHost, config.SEPort)
	return config, nil
}

// LoadWebDAVConfig
func LoadWebDAVConfig() (WebDAVConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("WD_USER")
	viper.BindEnv("WD_PASSWORD")
	viper.BindEnv("WD_HOST")
	viper.BindEnv("WD_PORT")
	var config WebDAVConfig
	if err := viper.Unmarshal(&config); err != nil {
		return WebDAVConfig{}, err
	}
	log.Infof("[WebDAV] user: %v, pass: %v, host: %v, port: %v", config.WDUser, config.WDPassword, config.WDSHost, config.WDPort)
	return config, nil
}

// LoadAppConfig
func LoadAppConfig() (AppConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("AP_ROOT")
	viper.BindEnv("AP_PORT")
	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return AppConfig{}, err
	}
	log.Infof("[App] root: %v, port: %v", config.APRoot, config.APPort)
	return config, nil
}
