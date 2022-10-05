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
}

// LoadMysqlConfig
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

// LoadRedisConfig
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

// LoadElasticSearchConfig
func LoadElasticSearchConfig() (ElasticSearchConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("SE_USER")
	viper.BindEnv("SE_PASSWORD")
	viper.BindEnv("SE_HOST")
	var config ElasticSearchConfig
	if err := viper.Unmarshal(&config); err != nil {
		return ElasticSearchConfig{}, err
	}
	log.Infof("user: %v, pass: %v, host: %v", config.SEUser, config.SEPassword, config.SESHost)
	return config, nil
}

// LoadWebDAVConfig
func LoadWebDAVConfig() (WebDAVConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("WD_USER")
	viper.BindEnv("WD_PASSWORD")
	viper.BindEnv("WD_HOST")
	var config WebDAVConfig
	if err := viper.Unmarshal(&config); err != nil {
		return WebDAVConfig{}, err
	}
	log.Infof("user: %v, pass: %v, host: %v", config.WDUser, config.WDPassword, config.WDSHost)
	return config, nil
}

// LoadAppConfig
func LoadAppConfig() (AppConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("AP_ROOT")
	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return AppConfig{}, err
	}
	log.Infof("root: %v", config.APRoot)
	return config, nil
}
