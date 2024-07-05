package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Memungkinkan viper dapat otomatis membaca config dari .env
	// contoh MYSQL_DBHOST akan dibaca viper.GetString("mysql.dbhost")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("failed to load the config, error: %v", err)
	}
}

func Port() string {
	return viper.GetString("port")
}
