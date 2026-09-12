package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	APPport       string
	DSN           string
	REDISaddr     string
	REDISpassword string
	RabbitMQ      string
	IpGeo         string
}

func NewConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found!")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
        log.Fatal("APP_PORT environment variable is required!")
    }

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environtment variable is required!")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        log.Fatal("REDIS_ADDR environment variable is required!")
    }

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
        log.Fatal("REDIS_PASSWORD environment variable is required!")
    }

	rabbitMQ := os.Getenv("RABBITMQ_DSN")
	if rabbitMQ == "" {
		log.Fatal("RABBITMQ_DSN environment variable is required!")

	}

	ipGeo := os.Getenv("IPGEO_API_KEY")
	if ipGeo == "" {
		log.Fatal("ADMIN_SECRET_KEY environment variable is required!")
	}

	return AppConfig{
		APPport:       appPort,
		DSN:           dsn,
		REDISaddr:     redisAddr,
		REDISpassword: redisPassword,
		RabbitMQ:      rabbitMQ,
		IpGeo:         ipGeo,
	}
}