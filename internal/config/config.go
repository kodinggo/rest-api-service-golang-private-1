package config

import (
	"log"
	"strings"
	"time"

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

func GRPCPort() string {
	return viper.GetString("grpc_port")
}

func JWTSigningKey() string {
	return viper.GetString("jwt.signing_key")
}

func JWTExp() time.Duration {
	return viper.GetDuration("jwt.exp")
}

func CommentgRPCHost() string {
	return viper.GetString("comment_service.grpc_host")
}
