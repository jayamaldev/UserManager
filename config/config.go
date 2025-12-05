package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	APPPort    int
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Println("error on parsing DB port")
		return nil, err
	}
	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		fmt.Println("error on parsing App port")
		return nil, err
	}

	if dbHost == "" || dbPort <= 0 || dbUser == "" || dbPwd == "" || dbName == "" || appPort <= 0 {
		fmt.Println("error on configurations")
		return nil, fmt.Errorf("error on configurations")
	}

	return &Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPwd,
		DBName:     dbName,
		APPPort:    appPort,
	}, nil
}
