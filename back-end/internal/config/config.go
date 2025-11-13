package config

import (
	"context"
	"entry-project/back-end/internal/model"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	AppPort   string
	DBHost    string
	DBUser    string
	DBPass    string
	DBName    string
	DBPort    string
	RedisAddr string

	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessTTL     string
	JWTRefreshTTL    string
}

func LoadConfig() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		DBHost:  getEnv("DB_HOST", "localhost"),
		DBUser:  getEnv("DB_USER", "khue"),
		DBPass:  getEnv("DB_PASSWORD", "27112000"),
		DBName:  getEnv("DB_NAME", "myapp"),
		DBPort:  getEnv("DB_PORT", "5432"),

		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),

		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "access-secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "refresh-secret"),
		JWTAccessTTL:     getEnv("JWT_ACCESS_TTL", "15m"),
		JWTRefreshTTL:    getEnv("JWT_REFRESH_TTL", "720h"),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func ConnectDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	log.Println("Connect to PostgreSQL")
	return db
}

func ConnectRedis(cfg *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, PoolSize: 1000})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect Redis:", err)
	}
	log.Println(("Connect to Redis"))
	return rdb
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&model.User{}, &model.Login{})
	log.Println("Database migrated")
}
