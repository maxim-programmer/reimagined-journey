package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr             string
	DatabaseURL      string
	UploadDir        string
	ElasticsearchURL string
	RedisAddr        string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading environment variables")
	}

	return &Config{
		Addr:             getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/knowledgedb?sslmode=disable"),
		UploadDir:        getEnv("UPLOAD_DIR", "./uploads"),
		ElasticsearchURL: getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}