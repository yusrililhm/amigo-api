package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type appConfig struct {
	AppPort       string
	DBUser        string
	DBHost        string
	DBPort        string
	DBName        string
	DBPassword    string
	DBDialect     string
	JWTSecretKey  string
	AdminFullName string
	AdminEmail    string
	AdminPassword string
	RedisHost     string
	RedisPort     string
	RedisPass     string
}

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("error while load .env file :", err.Error())
	}
}

func NewAppConfig() *appConfig {
	return &appConfig{
		AppPort:       os.Getenv("APP_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBName:        os.Getenv("DB_NAME"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBDialect:     os.Getenv("DB_DIALECT"),
		JWTSecretKey:  os.Getenv("JWT_SECRET_KEY"),
		AdminFullName: os.Getenv("ADMIN_FULL_NAME"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPass:     os.Getenv("REDIS_PASSWORD"),
	}
}
