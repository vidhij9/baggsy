package db

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     os.Getenv("DATABASE_URL"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBName:     os.Getenv("POSTGRES_DB"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
