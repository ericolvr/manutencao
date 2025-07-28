package config

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL  string
	ServerPort   string
	AllowOrigins string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	PostgresHost string
	PostgresPort string
	JWTSecret    string
}

var (
	cfg  *Config
	once sync.Once
	db   *sql.DB
)

func LoadConfig() *Config {
	once.Do(func() {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		cfg = &Config{
			DatabaseURL:  viper.GetString("DATABASE_URL"),
			ServerPort:   viper.GetString("SERVER_PORT"),
			AllowOrigins: viper.GetString("ALLOW_ORIGINS"),
			PostgresUser: viper.GetString("POSTGRES_USER"),
			PostgresPass: viper.GetString("POSTGRES_PASSWORD"),
			PostgresDB:   viper.GetString("POSTGRES_DB"),
			PostgresHost: viper.GetString("POSTGRES_HOST"),
			PostgresPort: viper.GetString("POSTGRES_PORT"),
			JWTSecret:    viper.GetString("JWT_SECRET"),
		}
	})
	return cfg
}

func GetDB() *sql.DB {
	if db != nil {
		return db
	}

	config := LoadConfig()

	connStr := config.DatabaseURL
	if connStr == "" {
		postgresHost := config.PostgresHost
		if postgresHost == "" {
			postgresHost = "postgres"
		}

		postgresPort := config.PostgresPort
		if postgresPort == "" {
			postgresPort = "5432"
		}

		postgresUser := config.PostgresUser
		if postgresUser == "" {
			postgresUser = "postgres"
		}

		postgresPass := config.PostgresPass
		if postgresPass == "" {
			postgresPass = "postgres"
		}

		postgresDB := config.PostgresDB
		if postgresDB == "" {
			postgresDB = "app_db"
		}

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			postgresHost, postgresPort, postgresUser, postgresPass, postgresDB)
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	return db
}
