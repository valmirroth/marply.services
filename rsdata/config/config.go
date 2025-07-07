package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: não encontrou .env, usando as variáveis de ambiente do sistema")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetDBConfig() DBConfig {
	return DBConfig{
		Host:     GetEnv("DB_HOST"),
		Port:     GetEnv("DB_PORT"),
		User:     GetEnv("DB_USER"),
		Password: GetEnv("DB_PASS"),
		Name:     GetEnv("DB_NAME"),
	}
}
